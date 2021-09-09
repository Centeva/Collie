package testutils

import "bitbucket.org/centeva/collie/packages/external"

type MockGitProvider struct {
	Called         map[string]int
	CalledWith     map[string][]interface{}
	AuthRes        *external.AuthModel
	GetBranchesRes []string
}

func NewMockGitProvider() *MockGitProvider {
	return &MockGitProvider{
		Called:     make(map[string]int),
		CalledWith: make(map[string][]interface{}),
	}
}

type GPCommentArgs struct {
	Workspace string
	Repo      string
	Branch    string
	Comment   string
}

func (m *MockGitProvider) Comment(workspace string, repo string, branch string, comment string) (err error) {
	m.Called["comment"]++
	m.CalledWith["comment"] = append(m.CalledWith["comment"], &GPCommentArgs{
		workspace,
		repo,
		branch,
		comment,
	})
	return
}

type GPGetOpenPRBranchesArgs struct {
	Workspace string
	Repo      string
}

func (m *MockGitProvider) GetOpenPRBranches(workspace string, repo string) (branches []string, err error) {
	m.Called["getopenprbranches"]++
	m.CalledWith["getopenprbranches"] = append(m.CalledWith["getopenprbranches"], &GPGetOpenPRBranchesArgs{
		workspace,
		repo,
	})

	return m.GetBranchesRes, nil
}

type GPAuthArgs struct {
	ClientId string
	Secret   string
	Username string
	Password string
}

func (m *MockGitProvider) BasicAuth(clientId string, secret string) (auth *external.AuthModel, err error) {
	m.Called["basicauth"]++

	m.CalledWith["basicauth"] = append(m.CalledWith["basicauth"], &GPAuthArgs{
		ClientId: clientId,
		Secret:   secret,
		Username: "",
		Password: "",
	})
	return m.AuthRes, nil
}

func (m *MockGitProvider) UserAuth(clientId string, secret string, username string, password string) (auth *external.AuthModel, err error) {
	m.Called["userauth"]++

	m.CalledWith["userauth"] = append(m.CalledWith["userauth"], &GPAuthArgs{
		ClientId: clientId,
		Secret:   secret,
		Username: username,
		Password: password,
	})
	return m.AuthRes, nil
}
