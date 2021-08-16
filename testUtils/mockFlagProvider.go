package testutils

import "bitbucket.org/centeva/collie/packages/external"

type MockFlagProvider struct {
	Called     map[string]int
	CalledWith map[string][]interface{}
	stringRes  string
}

func NewMockFlagProvider() *MockFlagProvider {
	res := &MockFlagProvider{
		Called:     make(map[string]int),
		CalledWith: make(map[string][]interface{}),
		stringRes:  "",
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

type StringArgs struct {
	name  string
	value string
	usage string
}

func (m *MockFlagProvider) PrintDefaults() {
	m.Called["printdefaults"]++
}

func (m *MockFlagProvider) String(name string, value string, usage string) *string {
	m.Called["string"]++
	m.CalledWith["string"] = append(m.CalledWith["stringvar"], &StringArgs{
		name,
		value,
		usage,
	})

	return &m.stringRes
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

func (m *MockFlagProvider) NewFlagSet(name string, usage string) external.IFlagSet {
	m.Called["newflagset"]++

	m.CalledWith["newflagset"] = append(m.CalledWith["newflagset"], struct {
		name  string
		usage string
	}{name, usage})
	return NewMockFlagSet("args1")
}

type mockFlagSet struct {
	called     map[string]int
	calledWith map[string][]interface{}
	argRes     string
	stringRes  string
}

func NewMockFlagSet(argRes string) *mockFlagSet {
	return &mockFlagSet{
		called:     make(map[string]int),
		calledWith: make(map[string][]interface{}),
		argRes:     argRes,
	}
}

func (m *mockFlagSet) PrintDefaults() {
	m.called["printdefauts"]++
}

func (m *mockFlagSet) String(name string, value string, usage string) *string {
	m.called["string"]++
	m.calledWith["stringvar"] = append(m.calledWith["stringvar"], &StringArgs{
		name,
		value,
		usage,
	})

	return &m.stringRes
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
