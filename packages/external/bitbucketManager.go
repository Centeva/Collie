package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type GitProviderFactory struct {
	BitbucketManager IGitProvider
}

func NewGitProviderFactory(bitbucketManager IGitProvider) *GitProviderFactory {
	return &GitProviderFactory{
		BitbucketManager: bitbucketManager,
	}
}

type IGitProvider interface {
	GetOpenPRBranches(workspace string, repo string) (branches []string, err error)
	Comment(workspace string, repo string, branch string, comment string, username *string, password *string) (err error)
	BasicAuth(clientId string, secret string) (auth *AuthModel, err error)
}

type ErrorModel struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type AuthModel struct {
	ErrorModel
	Scopes       string `json:"scopes"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	State        string `json:"state"`
	RefreshToken string `json:"refresh_token"`
}

type BranchModel struct {
	Name string `json:"name"`
}

type CommitModel struct {
	Hash string `json:"hash"`
}

type RefModel struct {
	Commit CommitModel `json:"commit"`
	Branch BranchModel `json:"branch"`
}

type PullRequestModel struct {
	Description string   `json:"description"`
	Title       string   `json:"title"`
	Id          int      `json:"id"`
	Destination RefModel `json:"destination"`
	Source      RefModel `json:"source"`
	State       string   `json:"state"`
}

type PaginatedResponse struct {
	ErrorModel
	PageLen int `json:"pagelen"`
	Page    int `json:"page"`
	Size    int `json:"size"`
}

type PaginatedPullRequestModel struct {
	PaginatedResponse
	Values []PullRequestModel `json:"values"`
}

type BitbucketManager struct {
	client *http.Client
	auth   *AuthModel
}

func NewBitbucketManager() *BitbucketManager {
	return &BitbucketManager{
		client: &http.Client{},
	}
}

func (m *BitbucketManager) authenticate(clientId string, secret string, data *url.Values) (auth *AuthModel, err error) {
	authUrl := "https://bitbucket.org/site/oauth2/access_token"

	dataEncoded := data.Encode()
	var req *http.Request
	if req, err = http.NewRequest("POST", authUrl, strings.NewReader(dataEncoded)); err != nil {
		return nil, errors.Wrap(err, "Failed to create request")
	}

	req.SetBasicAuth(clientId, secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(dataEncoded)))

	var res *http.Response
	if res, err = m.client.Do(req); err != nil {
		return nil, errors.Wrap(err, "Failed to make request")
	}

	if err = jsonUnmarshal(&auth, res); err != nil {
		return nil, errors.Wrap(err, "Failed to Unmarsal request")
	}

	if auth.ErrorCode != "" {
		return nil, errors.Errorf("API Error: %s %s", auth.ErrorCode, auth.ErrorDescription)
	}

	m.auth = auth
	return
}

func (m *BitbucketManager) BasicAuth(clientId string, secret string) (auth *AuthModel, err error) {
	data := &url.Values{
		"grant_type": []string{"client_credentials"},
	}

	return m.authenticate(clientId, secret, data)
}

func (m *BitbucketManager) getPrForBranch(workspace string, repo string, branch string) (pr *PullRequestModel, err error) {
	prPath := fmt.Sprintf(`https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests`, workspace, repo)
	prUrl, err := buildUrl(prPath, map[string]string{
		"q": fmt.Sprintf(`source.branch.name="%s"`, branch),
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to build Url")
	}

	var req *http.Request
	if req, err = http.NewRequest("GET", prUrl, nil); err != nil {
		return nil, errors.Wrap(err, "Failed to create request")
	}

	if err = m.addAuthHeader(req); err != nil {
		return nil, errors.Wrap(err, "Failed to add auth headers")
	}

	var res *http.Response
	if res, err = m.client.Do(req); err != nil {
		return nil, errors.Wrap(err, "Failed to make request")
	}

	var prData *PaginatedPullRequestModel
	if err = jsonUnmarshal(&prData, res); err != nil {
		return nil, errors.Wrap(err, "Failed to Unmarshal request")
	}

	if prData.ErrorCode != "" {
		return nil, errors.Errorf("API Error: %s %s", prData.ErrorCode, prData.ErrorDescription)
	}

	first := prData.Values[0]
	return &first, nil
}

func (m *BitbucketManager) addAuthHeader(req *http.Request) (err error) {
	if m.auth == nil {
		return fmt.Errorf("BitbucketManager: Auth must be called before Comment")
	}

	bearer := "Bearer " + m.auth.AccessToken
	req.Header.Add("Authorization", bearer)

	return
}

func (m *BitbucketManager) GetOpenPRBranches(workspace string, repo string) (branches []string, err error) {
	prPath := fmt.Sprintf(`https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests`, workspace, repo)
	prUrl, err := buildUrl(prPath, map[string]string{
		"state":  "OPEN",
		"fields": "values.source.branch.name,values.id",
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to build Url")
	}

	var req *http.Request
	if req, err = http.NewRequest("GET", prUrl, nil); err != nil {
		return nil, errors.Wrap(err, "Failed to get open pullRequests")
	}

	if err = m.addAuthHeader(req); err != nil {
		return nil, errors.Wrap(err, "Failed to add auth headers")
	}

	prRes, err := m.client.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get open pull requests")
	}

	if prRes.StatusCode == 400 {
		body, _ := ioutil.ReadAll(prRes.Body)
		return nil, errors.Errorf("Request Error: %s %s", prRes.Status, string(body))
	}

	var resModel *PaginatedPullRequestModel
	if err = jsonUnmarshal(&resModel, prRes); err != nil {
		return nil, errors.Wrap(err, "Failed to Unmarshal request")
	}

	if resModel.ErrorCode != "" {
		return nil, errors.Errorf("API Error: %s %s", resModel.ErrorCode, resModel.ErrorDescription)
	}

	for _, b := range resModel.Values {
		branches = append(branches, b.Source.Branch.Name)
	}

	return
}

func (m *BitbucketManager) Comment(workspace string, repo string, branch string, comment string, username *string, password *string) (err error) {
	pr, err := m.getPrForBranch(workspace, repo, branch)

	jsonStr := []byte(fmt.Sprintf(`{"content":{"raw":"%s"}}`, comment))

	if err != nil {
		return errors.Wrapf(err, "Failed to find pr for branch: %s", branch)
	}

	commentPath := fmt.Sprintf(`https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests/%d/comments`, workspace, repo, pr.Id)
	commentUrl, err := buildUrl(commentPath, make(map[string]string))

	if err != nil {
		return errors.Wrap(err, "Failed to build Url")
	}

	var req *http.Request
	if req, err = http.NewRequest("POST", commentUrl, bytes.NewBuffer(jsonStr)); err != nil {
		return errors.Wrap(err, "Failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	if username != nil && password != nil {
		req.SetBasicAuth(*username, *password)
	} else {
		if err = m.addAuthHeader(req); err != nil {
			return errors.Wrap(err, "Failed to add auth headers")
		}
	}

	commentRes, err := m.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to make request")
	}

	if commentRes.StatusCode == 400 {
		body, _ := ioutil.ReadAll(commentRes.Body)
		return errors.Errorf("Request Error: %s %s", commentRes.Status, string(body))
	}

	var resModel *struct {
		ErrorModel
	}
	if err = jsonUnmarshal(&resModel, commentRes); err != nil {
		return errors.Wrap(err, "Failed to Unmarshal request")
	}

	if resModel.ErrorCode != "" {
		return errors.Errorf("API Error: %s %s", resModel.ErrorCode, resModel.ErrorDescription)
	}

	return
}

func buildUrl(path string, queryParams map[string]string) (resUrl string, err error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to parse path: %s", path)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to parse query params: %s", u.RawQuery)
	}

	for key, param := range queryParams {
		q.Add(key, param)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func jsonUnmarshal(t interface{}, r *http.Response) (err error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&t); err != nil {
		return errors.Wrapf(err, "Failed to Unmarshal data to type: %T", &t)
	}

	return
}
