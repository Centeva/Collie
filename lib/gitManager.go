package lib

import (
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage"
)

type IStorageProvider interface {
	storage.Storer
}

type IGitRemote interface {
	List(o *git.ListOptions) (rfs []*plumbing.Reference, err error)
}

type IGitProvider interface {
	NewRemote(s storage.Storer, c *config.RemoteConfig) IGitRemote
}

type GitProvider struct{}

func (g GitProvider) NewRemote(s storage.Storer, c *config.RemoteConfig) *git.Remote {
	return git.NewRemote(s, c)
}

type GitManager struct {
	GitProvider     IGitProvider
	StorageProvider IStorageProvider
}

type ListRemoteRefOptions struct {
	RemoteConfig *config.RemoteConfig
	ListOptions  *git.ListOptions
}

func (m GitManager) ListRemoteRef(options *ListRemoteRefOptions) []string {
	rem := m.GitProvider.NewRemote(m.StorageProvider, options.RemoteConfig)
	refs, err := rem.List(options.ListOptions)
	checkIfError(err)

	refNames := make([]string, len(refs))
	for i, ref := range refs {
		refNames[i] = ref.Name().Short()
	}

	return refNames
}

func checkIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
