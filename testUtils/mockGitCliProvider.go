package testutils

import (
	"bitbucket.org/centeva/collie/packages/external"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage"
)

type GitProviderTestArgs struct {
	Branches []string
}

type MockGitCliProvider struct {
	TestArgs GitProviderTestArgs
}

func (g MockGitCliProvider) Clone(s storage.Storer, worktree billy.Filesystem, o *git.CloneOptions) (*git.Repository, error) {
	repo := &git.Repository{}
	return repo, nil
}

type MockRemote struct {
	Branches []string
}

func (m MockRemote) List(o *git.ListOptions) (rfs []*plumbing.Reference, err error) {
	rfs = make([]*plumbing.Reference, len(m.Branches))

	for i, b := range m.Branches {
		rfs[i] = plumbing.NewReferenceFromStrings(b, "origin")
	}

	return
}

func (g MockGitCliProvider) NewRemote(s storage.Storer, c *config.RemoteConfig) external.IGitRemote {
	val := &MockRemote{
		Branches: g.TestArgs.Branches,
	}
	return val
}
