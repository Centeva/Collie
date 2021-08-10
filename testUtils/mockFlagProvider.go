package testUtils

import "bitbucket.org/centeva/collie/packages/external"

type MockFlagProvider struct {
	Called     map[string]int
	CalledWith map[string][]interface{}
}

func NewMockFlagProvider() *MockFlagProvider {
	res := &MockFlagProvider{
		Called:     make(map[string]int),
		CalledWith: make(map[string][]interface{}),
	}

	return res
}

func (m *MockFlagProvider) Parse() {
	m.Called["parse"]++
}

type StringVarArgs struct {
	p     *string
	name  string
	value string
	usage string
}

func (m *MockFlagProvider) StringVar(p *string, name string, value string, usage string) {
	m.Called["stringvar"]++

	m.CalledWith["stringvar"] = append(m.CalledWith["stringvar"], StringVarArgs{
		p,
		name,
		value,
		usage,
	})
}

func (m *MockFlagProvider) NewFlagSet(name string) external.IFlagSet {
	m.Called["newflagset"]++

	m.CalledWith["newflagset"] = append(m.CalledWith["newflagset"], struct{ name string }{name})
	return NewMockFlagSet("args1")
}

type mockFlagSet struct {
	called     map[string]int
	calledWith map[string][]interface{}
	argRes     string
}

func NewMockFlagSet(argRes string) *mockFlagSet {
	return &mockFlagSet{
		called:     make(map[string]int),
		calledWith: make(map[string][]interface{}),
		argRes:     argRes,
	}
}

func (m *mockFlagSet) StringVar(p *string, name string, value string, usage string) {
	m.called["stringvar"]++

	m.calledWith["stringvar"] = append(m.calledWith["stringvar"], StringVarArgs{
		p,
		name,
		value,
		usage,
	})
}

func (m *mockFlagSet) Parse(arguments []string) error {
	m.called["parse"]++
	return nil
}

func (m *mockFlagSet) Arg(i int) string {
	m.called["arg"]++

	m.calledWith["arg"] = append(m.calledWith["arg"], struct{ i int }{i})
	return m.argRes
}
