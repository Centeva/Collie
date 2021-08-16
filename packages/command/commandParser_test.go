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
	sut := command.NewCommandParser(flagProvider, gitProviderFactory, nil)

	sut.ParseCommands()

	gotParse := flagProvider.Called["parse"]
	gotStringVar := flagProvider.Called["stringvar"]

	if gotParse != 1 {
		t.Errorf("ParseFlags(): flagProvider.Parse() Should have been called once got: %v", flagProvider.Called)
	}
	if gotStringVar != 1 {
		t.Errorf("ParseFlags(): flagProvider.StringVar() Should have been called once got: %v", flagProvider.Called)
	}
}
