package command_test

import (
	"strings"
	"testing"

	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/packages/external"
	"bitbucket.org/centeva/collie/testutils"
)

func testSetup(args *TestArgs) *command.PRCommentCommand {
	mockGitFactory := &external.GitProviderFactory{
		BitbucketManager: testutils.NewMockGitProvider(),
	}
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewPRCommentCommand(mockFlagProvider, mockGitFactory)
	sut.GitSource = newBitBucketSource(&args.GitSourceArgs)
	return sut
}

type TestArgs struct {
	GitProvider   string
	GitSourceArgs bitBucketSourceArgs
	Want          string
}

type bitBucketSourceArgs struct {
	Branch    string
	ClientId  string
	Comment   string
	Repo      string
	Secret    string
	Workspace string
	Username  string
	Password  string
}

func newBitBucketSource(args *bitBucketSourceArgs) *command.BitBucketSource {
	return &command.BitBucketSource{
		Branch:    &args.Branch,
		ClientId:  &args.ClientId,
		Comment:   &args.Comment,
		Repo:      &args.Repo,
		Secret:    &args.Secret,
		Workspace: &args.Workspace,
		Username:  &args.Username,
		Password:  &args.Password,
	}
}

func Test_prCommentCommand_HappyPath(t *testing.T) {

	args := &TestArgs{
		GitProvider: "bitbucket",
		GitSourceArgs: bitBucketSourceArgs{
			Branch:    "test1",
			ClientId:  "test2",
			Comment:   "test3",
			Repo:      "test4",
			Secret:    "test5",
			Workspace: "test6",
			Username:  "test7",
			Password:  "test8",
		},
	}

	sut := testSetup(args)
	err := sut.Execute()

	if err != nil {
		t.Errorf("prCommentCommand_HappyPath: should not error, %s", err)
	}

}
func Test_prCommentCommand_Errors(t *testing.T) {

	tests := []TestArgs{
		{
			GitProvider: "bitbucket",
			GitSourceArgs: bitBucketSourceArgs{
				Branch:    "",
				ClientId:  "test2",
				Comment:   "test3",
				Repo:      "test4",
				Secret:    "test5",
				Workspace: "test6",
				Username:  "test7",
				Password:  "test8",
			},
			Want: "Branch is required",
		},
		{
			GitProvider: "bitbucket",
			GitSourceArgs: bitBucketSourceArgs{
				Branch:    "test1",
				ClientId:  "",
				Comment:   "test3",
				Repo:      "test4",
				Secret:    "test5",
				Workspace: "test6",
				Username:  "test7",
				Password:  "test8",
			},
			Want: "ClientId is required",
		},
		{
			GitProvider: "bitbucket",
			GitSourceArgs: bitBucketSourceArgs{
				Branch:    "test1",
				ClientId:  "test2",
				Comment:   "",
				Repo:      "test4",
				Secret:    "test5",
				Workspace: "test6",
				Username:  "test7",
				Password:  "test8",
			},
			Want: "Comment is required",
		},
		{
			GitProvider: "bitbucket",
			GitSourceArgs: bitBucketSourceArgs{
				Branch:    "test1",
				ClientId:  "test2",
				Comment:   "test3",
				Repo:      "",
				Secret:    "test5",
				Workspace: "test6",
				Username:  "test7",
				Password:  "test8",
			},
			Want: "Repo is required",
		},
		{
			GitProvider: "bitbucket",
			GitSourceArgs: bitBucketSourceArgs{
				Branch:    "test1",
				ClientId:  "test2",
				Comment:   "test3",
				Repo:      "test4",
				Secret:    "",
				Workspace: "test6",
				Username:  "test7",
				Password:  "test8",
			},
			Want: "Secret is required",
		},
		{
			GitProvider: "bitbucket",
			GitSourceArgs: bitBucketSourceArgs{
				Branch:    "test1",
				ClientId:  "test2",
				Comment:   "test3",
				Repo:      "test4",
				Secret:    "test5",
				Workspace: "",
				Username:  "test7",
				Password:  "test8",
			},
			Want: "Workspace is required",
		},
	}

	for _, args := range tests {
		sut := testSetup(&args)

		var err error
		switch s := sut.GitSource.(type) {
		case *command.BitBucketSource:
			err = sut.ValidateBitbucketFlags(s)
		}

		if err == nil {
			t.Errorf("prCommentCommand_errors: FlagsValid should error but didn't. want: %s", args.Want)
			continue
		}

		if !strings.Contains(err.Error(), args.Want) {
			t.Errorf("prCommentCommand_errors: Error should contain '%s' but got '%s'", args.Want, err.Error())
		}
	}

}
