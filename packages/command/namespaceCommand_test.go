package command_test

import (
	"log"
	"testing"

	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/testutils"
)

func Test_executeOutCluster(t *testing.T) {
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewNamespaceCommand(mockFlagProvider, mockKubernetesManager)

	log.Printf("setup")
	kubeconfig := "/test/path.conf"
	sut.Kubeconfig = &kubeconfig
	timeout := "10m"
	sut.Timeout = &timeout
	sut.Execute()

	log.Printf("executed")
	if mockKubernetesManager.Called["outclusterconfig"] != 1 {
		t.Errorf("OutClusterConfig() should have been called once")
	}

	firstArg := mockKubernetesManager.CalledWith["outclusterconfig"][0]
	switch arg := firstArg.(type) {
	case *testutils.KMOutClusterConfigArgs:
		flat := *arg
		if flat.Kubeconfig == "/test/path.conf" {
			return
		}
	}

	t.Errorf("OutClusterConfig() should have been called with Kueconfig: %s but got %+v", kubeconfig, firstArg)
}

func Test_executeInCluster(t *testing.T) {
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewNamespaceCommand(mockFlagProvider, mockKubernetesManager)

	timeout := "10m"
	sut.Timeout = &timeout
	sut.Execute()

	if mockKubernetesManager.Called["inclusterconfig"] != 1 {
		t.Errorf("InClusterConfig() should have been called once")
	}
}

func Test_executeGetNamespaces(t *testing.T) {
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewNamespaceCommand(mockFlagProvider, mockKubernetesManager)

	timeout := "10m"
	sut.Timeout = &timeout
	sut.Execute()

	if mockKubernetesManager.Called["getnamespaces"] != 1 {
		t.Errorf("GetNamespaces() should have been called once")
	}

	firstArg := mockKubernetesManager.CalledWith["getnamespaces"][0]
	switch arg := firstArg.(type) {
	case *testutils.KMGetNamespacesArgs:
		flat := *arg
		if flat.Label == "" {
			return
		}
	}

	t.Errorf("GetNamespaces() should have been called with empty Label but got %+v", firstArg)
}

func Test_executeDeleteNamespace(t *testing.T) {
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockKubernetesManager.GetNamespacesRes = []string{"test-1", "other"}
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewNamespaceCommand(mockFlagProvider, mockKubernetesManager)

	timeout := "10m"
	sut.Timeout = &timeout
	sut.Namespace = "test-1"
	sut.Execute()

	if mockKubernetesManager.Called["deletenamespace"] != 1 {
		t.Errorf("GetNamespaces() should have been called once")
	}

	firstArg := mockKubernetesManager.CalledWith["deletenamespace"][0]
	switch arg := firstArg.(type) {
	case *testutils.KMDeleteNamespaceArgs:
		flat := *arg
		if flat.Namespace == "test-1" {
			return
		}
	}

	t.Errorf("GetNamespaces() should have been called with Namespace: %s but got %+v", sut.Namespace, firstArg)
}
