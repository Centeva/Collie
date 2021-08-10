package external_test

import (
	"strings"
	"testing"

	"bitbucket.org/centeva/collie/packages/external"
	"bitbucket.org/centeva/collie/testUtils"
	"github.com/go-git/go-git/v5/storage/memory"
)

func Test_ListRemoteRef(t *testing.T) {
	type args struct {
		gitProviderArgs *testUtils.GitProviderTestArgs
	}

	defaultArgs := args{
		&testUtils.GitProviderTestArgs{
			Branches: []string{
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
			gitProvider := &testUtils.MockGitProvider{
				TestArgs: *tt.args.gitProviderArgs,
			}

			storageProvider := memory.NewStorage()

			sut := &external.GitManager{
				GitProvider:     gitProvider,
				StorageProvider: storageProvider,
			}

			refs := sut.ListRemoteRef(&external.ListRemoteRefOptions{})
			got := strings.Join(refs, ",")

			if tt.want != got {
				t.Errorf("ListRemoteRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
