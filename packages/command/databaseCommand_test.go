package command_test

import (
	"testing"

	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/testutils"
)

func Test_executeConnect(t *testing.T) {
	mockPostgresManager := testutils.NewMockPostgresManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewDatabaseCommand(mockFlagProvider, mockPostgresManager)
	connectionString := "testConnString"
	sut.ConnectionString = &connectionString
	sut.Execute()

	if mockPostgresManager.Called["connect"] != 1 {
		t.Errorf("Connect() should have been called once")
	}

	firstArg := mockPostgresManager.CalledWith["connect"][0]
	switch arg := firstArg.(type) {
	case *testutils.PMConnectArgs:
		flat := *arg
		if flat.ConnectionString == connectionString {
			return
		}
	}

	t.Errorf("Connect() should have been called with ConnectionString: %s but got %+v", connectionString, firstArg)
}

func Test_executeDeleteDatabase(t *testing.T) {
	mockPostgresManager := testutils.NewMockPostgresManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewDatabaseCommand(mockFlagProvider, mockPostgresManager)
	connectionString := "testConnString"
	sut.ConnectionString = &connectionString
	sut.Database = "testDB"
	sut.Execute()

	if mockPostgresManager.Called["deletedatabase"] != 1 {
		t.Errorf("DeleteDatabase() should have been called once")
	}

	firstArg := mockPostgresManager.CalledWith["deletedatabase"][0]
	switch arg := firstArg.(type) {
	case *testutils.PMDeleteDatabaseArgs:
		flat := *arg
		if flat.Database == sut.Database {
			return
		}
	}

	t.Errorf("Connect() should have been called with ConnectionString: %s but got %+v", connectionString, firstArg)
}

func Test_executeClose(t *testing.T) {
	mockPostgresManager := testutils.NewMockPostgresManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewDatabaseCommand(mockFlagProvider, mockPostgresManager)
	connectionString := "testConnString"
	sut.ConnectionString = &connectionString
	sut.Execute()

	if mockPostgresManager.Called["close"] != 1 {
		t.Errorf("Connect() should have been called once")
	}
}
