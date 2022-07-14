package external

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

type GHPullRequestModel struct {
	Id     string     `json:"id"`
	Number string     `json:"number"`
	Head   GHRefModel `json:"head"`
	Base   GHRefModel `json:"base"`
}

type GHPullRequestFlatModel struct {
	Id          string `json:"id"`
	Number      int    `json:"number"`
	HeadRefName string `json:"headRefName"`
	BaseRefName string `json:"baseRefName"`
}

type GHRefModel struct {
	Ref string `json:"ref"`
}

type GHAuth struct {
	Username  string
	patSecret string
}

type GithubManager struct {
	client    *http.Client
	gqlClient *graphql.Client
	ctx       context.Context
	auth      *GHAuth
}

func NewGithubManager() *GithubManager {
	client := &http.Client{}
	return &GithubManager{
		client:    client,
		ctx:       context.Background(),
		gqlClient: graphql.NewClient("https://api.github.com/graphql", graphql.WithHTTPClient(client)),
		auth:      &GHAuth{},
	}
}

func (m *GithubManager) BasicAuth(clientId string, secret string) (auth *AuthModel, err error) {
	m.auth = &GHAuth{
		Username:  clientId,
		patSecret: secret,
	}

	return
}

type Setable interface {
	Set(key string, value string)
}

func (m *GithubManager) setAuth(req Setable) {
	req.Set("Authorization", fmt.Sprintf("Bearer %s", m.auth.patSecret))
}

func (m *GithubManager) Comment(workspace string, repo string, branch string, comment string, username *string, password *string) (err error) {
	pr, err := m.getPrForBranch(workspace, repo, branch)

	jsonStr := []byte(fmt.Sprintf(`{"body":"%s"}`, comment))

	if err != nil {
		return errors.Wrapf(err, "Failed to find pr for branch: %s", branch)
	}

	commentPath := fmt.Sprintf(`https://api.github.com/repos/%s/%s/issues/%d/comments`, workspace, repo, pr.Number)
	commentUrl, err := buildUrl(commentPath, make(map[string]string))

	if err != nil {
		return errors.Wrap(err, "Failed to build Url")
	}
	var req *http.Request
	if req, err = http.NewRequest("POST", commentUrl, bytes.NewBuffer(jsonStr)); err != nil {
		return errors.Wrap(err, "Failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	m.setAuth(req.Header)

	commentRes, err := m.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to make request")
	}

	if commentRes.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(commentRes.Body)
		return errors.Errorf("Request Error: %s %s", commentRes.Status, string(body))
	}

	return
}

func (m *GithubManager) getPrForBranch(workspace string, repo string, branch string) (pr *GHPullRequestFlatModel, err error) {
	req := graphql.NewRequest(`query($owner: String!, $repo: String!, $branch: String!) {
		repository(owner: $owner, name: $repo) {
			pullRequests(first: 1, states: [OPEN], headRefName: $branch) {
				nodes {
					headRefName
					baseRefName
					number
					id
				}
			}
		}
	}`)

	req.Var("owner", workspace)
	req.Var("repo", repo)
	req.Var("branch", branch)

	m.setAuth(req.Header)

	var resData struct {
		Repository struct {
			PullRequest struct {
				Nodes []struct {
					HeadRefName string `json:"headRefName"`
					BaseRefName string `json:"baseRefName"`
					Number      int    `json:"number"`
					Id          string `json:"id"`
				} `json:"nodes"`
			} `json:"pullRequests"`
		} `json:"repository"`
	}

	if err = m.gqlClient.Run(m.ctx, req, &resData); err != nil {
		return nil, errors.Wrap(err, "Failed to make request")
	}

	if len(resData.Repository.PullRequest.Nodes) == 0 {
		return nil, errors.Wrap(err, "No pull requests found")
	}

	node := resData.Repository.PullRequest.Nodes[0]

	pr = &GHPullRequestFlatModel{
		Id:          node.Id,
		Number:      node.Number,
		HeadRefName: node.HeadRefName,
		BaseRefName: node.BaseRefName,
	}

	return
}

func (m *GithubManager) GetOpenPRBranches(workspace string, repo string) (branches []string, err error) {
	prPath := fmt.Sprintf(`https://api.github.com/repos/%s/%s/pulls`, workspace, repo)
	prUrl, err := buildUrl(prPath, make(map[string]string))

	if err != nil {
		return nil, errors.Wrap(err, "Failed to build Url")
	}

	var req *http.Request
	if req, err = http.NewRequest("GET", prUrl, nil); err != nil {
		return nil, errors.Wrap(err, "Failed to get pullRequests")
	}

	m.setAuth(req.Header)

	prRes, err := m.client.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get open pullRequests")
	}

	if prRes.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(prRes.Body)
		return nil, errors.Errorf("Request Error: %s %s", prRes.Status, string(body))
	}

	var resModel *[]GHPullRequestModel
	if err = jsonUnmarshal(&resModel, prRes); err != nil {
		return nil, errors.Wrap(err, "Failed to Unmarshal request")
	}

	for _, b := range *resModel {
		branches = append(branches, b.Head.Ref)
	}

	return
}
