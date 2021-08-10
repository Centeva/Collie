package command_test

import (
	"testing"

	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/testUtils"
)

func Test_ParseCommands(t *testing.T) {
	flagProvider := testUtils.NewMockFlagProvider()
	sut := command.NewCommandParser(flagProvider)

	sut.ParseCommands()

	gotParse := flagProvider.Called["parse"]
	gotStringVar := flagProvider.Called["stringvar"]

	if gotParse != 1 {
		t.Errorf("ParseFlags(): flagProvider.Parse() Should have been called once got: %v", flagProvider.Called)
	}
	if gotStringVar != 2 {
		t.Errorf("ParseFlags(): flagProvider.StringVar() Should have been called twice got: %v", flagProvider.Called)
	}
}
