package testutils

import "bitbucket.org/centeva/collie/packages/external"

type MockGitProvider struct {
	Called         map[string]int
	CalledWith     map[string][]interface{}
	AuthRes        *external.AuthModel
	getBranchesRes []string
}

func NewMockGitProvider() *MockGitProvider {
	return &MockGitProvider{
		Called:     make(map[string]int),
		CalledWith: make(map[string][]interface{}),
	}
}

type commentArgs struct {
	workspace string
	repo      string
	branch    string
	comment   string
}

func (m *MockGitProvider) Comment(workspace string, repo string, branch string, comment string) (err error) {
	m.Called["comment"]++
	m.CalledWith["comment"] = append(m.CalledWith["comment"], &commentArgs{
		workspace,
		repo,
		branch,
		comment,
	})
	return
}

type getOpenPRBranchesArgs struct {
	workspace string
	repo      string
}

func (m *MockGitProvider) GetOpenPRBranches(workspace string, repo string) (branches []string, err error) {
	m.Called["getopenprbranches"]++
	m.CalledWith["getopenprbranches"] = append(m.CalledWith["getopenprbranches"], &getOpenPRBranchesArgs{
		workspace,
		repo,
	})

	return m.getBranchesRes, nil
}

type authArgs struct {
	clientId string
	secret   string
	username string
	password string
}

func (m *MockGitProvider) BasicAuth(clientId string, secret string) (auth *external.AuthModel, err error) {
	m.Called["basicauth"]++

	m.CalledWith["basicauth"] = append(m.CalledWith["basicauth"], &authArgs{
		clientId: clientId,
		secret:   secret,
		username: "",
		password: "",
	})
	return m.AuthRes, nil
}

func (m *MockGitProvider) UserAuth(clientId string, secret string, username string, password string) (auth *external.AuthModel, err error) {
	m.Called["userauth"]++

	m.CalledWith["userauth"] = append(m.CalledWith["userauth"], &authArgs{
		clientId: clientId,
		secret:   secret,
		username: username,
		password: password,
	})
	return m.AuthRes, nil
}
