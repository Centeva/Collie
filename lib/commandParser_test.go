package lib_test

import (
	"testing"

	"bitbucket.org/centeva/collie/lib"
)

type mockFlagProvider struct {
	called map[string]int
}

func NewMockFlagProvider() *mockFlagProvider {
	res := &mockFlagProvider{
		called: make(map[string]int),
	}

	return res
}

func (m *mockFlagProvider) Parse() {
	m.called["parse"]++
}

func (m *mockFlagProvider) StringVar(p *string, name string, value string, usage string) {
	m.called["stringvar"]++
}

func Test_ParseCommands(t *testing.T) {
	flagProvider := NewMockFlagProvider()
	sut := lib.NewCommandParser(flagProvider)

	sut.ParseCommands()

	gotParse := flagProvider.called["parse"]
	gotStringVar := flagProvider.called["stringvar"]

	if gotParse != 1 {
		t.Errorf("ParseFlags(): flagProvider.Parse() Should have been called once got: %v", flagProvider.called)
	}
	if gotStringVar != 2 {
		t.Errorf("ParseFlags(): flagProvider.StringVar() Should have been called twice got: %v", flagProvider.called)
	}
}
