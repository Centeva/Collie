package lib_test

import (
	"strings"
	"testing"

	"bitbucket.org/centeva/collie/lib"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
)

func Test_ListRemoteRef(t *testing.T) {
	type args struct {
		gitProviderArgs *gitProviderTestArgs
	}

	defaultArgs := args{
		&gitProviderTestArgs{
			branches: []string{
				"branch-1",
				"branch-2",
				"branch-3",
				"branch-4",
			},
		},
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"ListRemoteRef_happyPath", defaultArgs, "branch-1,branch-2,branch-3,branch-4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitProvider := &mockGitProvider{
				testArgs: *tt.args.gitProviderArgs,
			}

			storageProvider := memory.NewStorage()

			sut := &lib.GitManager{
				GitProvider:     gitProvider,
				StorageProvider: storageProvider,
			}

			refs := sut.ListRemoteRef(&lib.ListRemoteRefOptions{})
			got := strings.Join(refs, ",")

			if tt.want != got {
				t.Errorf("ListRemoteRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

type gitProviderTestArgs struct {
	branches []string
}
type mockGitProvider struct {
	testArgs gitProviderTestArgs
}

func (g mockGitProvider) Clone(s storage.Storer, worktree billy.Filesystem, o *git.CloneOptions) (*git.Repository, error) {
	repo := &git.Repository{}
	return repo, nil
}

type MockRemote struct {
	branches []string
}

func (m MockRemote) List(o *git.ListOptions) (rfs []*plumbing.Reference, err error) {
	rfs = make([]*plumbing.Reference, len(m.branches))

	for i, b := range m.branches {
		rfs[i] = plumbing.NewReferenceFromStrings(b, "origin")
	}

	return
}

func (g mockGitProvider) NewRemote(s storage.Storer, c *config.RemoteConfig) lib.IGitRemote {
	val := &MockRemote{
		branches: g.testArgs.branches,
	}
	return val
}
