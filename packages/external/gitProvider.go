package external

type GitProviderFactory struct {
	BitbucketManager IGitProvider
	GithubManager    IGitProvider
}

func NewGitProviderFactory(bitbucketManager IGitProvider, githubManager IGitProvider) *GitProviderFactory {
	return &GitProviderFactory{
		BitbucketManager: bitbucketManager,
		GithubManager:    githubManager,
	}
}

type IGitProvider interface {
	GetOpenPRBranches(workspace string, repo string) (branches []string, err error)
	Comment(workspace string, repo string, branch string, comment string, username *string, password *string) (err error)
	BasicAuth(clientId string, secret string) (auth *AuthModel, err error)
}