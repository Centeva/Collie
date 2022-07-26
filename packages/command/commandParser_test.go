package command_test

import (
	"testing"

	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/packages/external"
	"bitbucket.org/centeva/collie/testutils"
)

func Test_ParseCommands(t *testing.T) {
	flagProvider := testutils.NewMockFlagProvider()
	gitProviderFactory := &external.GitProviderFactory{BitbucketManager: &testutils.MockGitProvider{}}
	sut := command.NewCommandParser(flagProvider, gitProviderFactory, nil, nil, nil)

	sut.ParseCommands()

	gotParse := flagProvider.Called["parse"]

	if gotParse != 0 {
		t.Errorf("ParseFlags(): flagProvider.Parse() Should not have been called; got: %v", flagProvider.Called)
	}
}
