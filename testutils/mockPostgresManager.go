package testutils

type MockPostgresManager struct {
	Called     map[string]int
	CalledWith map[string][]interface{}
}

func NewMockPostgresManager() *MockPostgresManager {
	return &MockPostgresManager{
		Called:     make(map[string]int),
		CalledWith: make(map[string][]interface{}),
	}
}

type PMConnectArgs struct {
	ConnectionString string
}

func (m *MockPostgresManager) Connect(connectionString string) (err error) {
	m.Called["connect"]++
	m.CalledWith["connect"] = append(m.CalledWith["connect"], &PMConnectArgs{connectionString})
	return
}

type PMDeleteDatabaseArgs struct {
	Database string
}

func (m *MockPostgresManager) DeleteDatabase(database string) (err error) {
	m.Called["deletedatabase"]++
	m.CalledWith["deletedatabase"] = append(m.CalledWith["deletedatabase"], &PMDeleteDatabaseArgs{database})
	return
}

func (m *MockPostgresManager) Close() {
	m.Called["close"]++
}
