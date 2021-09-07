package command_test

import (
	"testing"

	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/packages/external"
	"bitbucket.org/centeva/collie/testutils"
)

func Test_ExecuteBasicAuth(t *testing.T) {
	mockFileReader := testutils.NewMockFileReader("testFile")
	mockBitbucketManager := testutils.NewMockGitProvider()
	mockGitProviderFactory := &external.GitProviderFactory{
		BitbucketManager: mockBitbucketManager,
	}
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewCleanupCommand(mockFlagProvider, mockKubernetesManager, mockFileReader, mockGitProviderFactory)

	bitbucketArgs := &command.ConfigBitbucketArgs{
		ClientId:  "testClientId",
		Secret:    "testSecret",
		Workspace: "testWorkspace",
		Repo:      "testRepo",
	}
	gitProvider := &command.ConfigGitProvider{
		Bitbucket: bitbucketArgs,
	}

	jobConfig := &external.CleanupJobConfig{}

	cleanupConfig := &command.CleanupConfig{
		Kubeconfig:  "kubeconfig",
		GitProvider: gitProvider,
		JobConfig:   jobConfig,
	}

	namespaceLabel := "testLabel"

	sut.CleanupConfig = cleanupConfig
	sut.NamespaceLabel = &namespaceLabel
	sut.Execute()

	if mockBitbucketManager.Called["basicauth"] != 1 {
		t.Errorf("BasicAuth() should have been called once")
	}

	firstArgs := mockBitbucketManager.CalledWith["basicauth"][0]
	switch args := firstArgs.(type) {
	case *testutils.GPAuthArgs:
		flat := *args

		if flat.ClientId != bitbucketArgs.ClientId {
			t.Errorf("BasicAuth() should have been called with ClientId: %s but got %+v", bitbucketArgs.ClientId, firstArgs)
		}
		if flat.Secret != bitbucketArgs.Secret {
			t.Errorf("BasicAuth() should have been called with Secret: %s but got %+v", bitbucketArgs.Secret, firstArgs)
		}
	}
}

func Test_ExecuteGetOpenPRBranches(t *testing.T) {
	mockFileReader := testutils.NewMockFileReader("testFile")
	mockBitbucketManager := testutils.NewMockGitProvider()
	mockGitProviderFactory := &external.GitProviderFactory{
		BitbucketManager: mockBitbucketManager,
	}
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewCleanupCommand(mockFlagProvider, mockKubernetesManager, mockFileReader, mockGitProviderFactory)

	bitbucketArgs := &command.ConfigBitbucketArgs{
		ClientId:  "testClientId",
		Secret:    "testSecret",
		Workspace: "testWorkspace",
		Repo:      "testRepo",
	}
	gitProvider := &command.ConfigGitProvider{
		Bitbucket: bitbucketArgs,
	}

	jobConfig := &external.CleanupJobConfig{}

	cleanupConfig := &command.CleanupConfig{
		Kubeconfig:  "kubeconfig",
		GitProvider: gitProvider,
		JobConfig:   jobConfig,
	}

	namespaceLabel := "testLabel"

	sut.CleanupConfig = cleanupConfig
	sut.NamespaceLabel = &namespaceLabel
	sut.Execute()

	if mockBitbucketManager.Called["getopenprbranches"] != 1 {
		t.Errorf("GetOpenPRBranches() should have been called once")
	}

	firstArgs := mockBitbucketManager.CalledWith["getopenprbranches"][0]
	switch args := firstArgs.(type) {
	case *testutils.GPGetOpenPRBranchesArgs:
		flat := *args

		if flat.Workspace != bitbucketArgs.Workspace {
			t.Errorf("GetOpenPRBranches() should have been called with Workspace: %s but got %+v", bitbucketArgs.Workspace, firstArgs)
		}
		if flat.Repo != bitbucketArgs.Repo {
			t.Errorf("GetOpenPRBranches() should have been called with Repo: %s but got %+v", bitbucketArgs.Repo, firstArgs)
		}
	}
}

func Test_ExecuteOutClusterConfig(t *testing.T) {
	mockFileReader := testutils.NewMockFileReader("testFile")
	mockBitbucketManager := testutils.NewMockGitProvider()
	mockGitProviderFactory := &external.GitProviderFactory{
		BitbucketManager: mockBitbucketManager,
	}
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewCleanupCommand(mockFlagProvider, mockKubernetesManager, mockFileReader, mockGitProviderFactory)

	bitbucketArgs := &command.ConfigBitbucketArgs{
		ClientId:  "testClientId",
		Secret:    "testSecret",
		Workspace: "testWorkspace",
		Repo:      "testRepo",
	}
	gitProvider := &command.ConfigGitProvider{
		Bitbucket: bitbucketArgs,
	}

	jobConfig := &external.CleanupJobConfig{}

	cleanupConfig := &command.CleanupConfig{
		Kubeconfig:  "kubeconfig",
		GitProvider: gitProvider,
		JobConfig:   jobConfig,
	}

	namespaceLabel := "testLabel"

	sut.CleanupConfig = cleanupConfig
	sut.NamespaceLabel = &namespaceLabel
	sut.Execute()

	if mockKubernetesManager.Called["outclusterconfig"] != 1 {
		t.Errorf("OutClusterConfig() should have been called once")
	}

	firstArgs := mockKubernetesManager.CalledWith["outclusterconfig"][0]
	switch args := firstArgs.(type) {
	case *testutils.KMOutClusterConfigArgs:
		flat := *args

		if flat.Kubeconfig != cleanupConfig.Kubeconfig {
			t.Errorf("OutClusterConfig() should have been called with Kubeconfig: %s but got %+v", cleanupConfig.Kubeconfig, firstArgs)
		}
	}
}

func Test_ExecuteInClusterConfig(t *testing.T) {
	mockFileReader := testutils.NewMockFileReader("testFile")
	mockBitbucketManager := testutils.NewMockGitProvider()
	mockGitProviderFactory := &external.GitProviderFactory{
		BitbucketManager: mockBitbucketManager,
	}
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewCleanupCommand(mockFlagProvider, mockKubernetesManager, mockFileReader, mockGitProviderFactory)

	bitbucketArgs := &command.ConfigBitbucketArgs{
		ClientId:  "testClientId",
		Secret:    "testSecret",
		Workspace: "testWorkspace",
		Repo:      "testRepo",
	}
	gitProvider := &command.ConfigGitProvider{
		Bitbucket: bitbucketArgs,
	}

	jobConfig := &external.CleanupJobConfig{}

	cleanupConfig := &command.CleanupConfig{
		Kubeconfig:  "",
		GitProvider: gitProvider,
		JobConfig:   jobConfig,
	}

	namespaceLabel := "testLabel"

	sut.CleanupConfig = cleanupConfig
	sut.NamespaceLabel = &namespaceLabel
	sut.Execute()

	if mockKubernetesManager.Called["inclusterconfig"] != 1 {
		t.Errorf("InClusterConfig() should have been called once")
	}
}

func Test_ExecuteGetNamespaces(t *testing.T) {
	mockFileReader := testutils.NewMockFileReader("testFile")
	mockBitbucketManager := testutils.NewMockGitProvider()
	mockGitProviderFactory := &external.GitProviderFactory{
		BitbucketManager: mockBitbucketManager,
	}
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewCleanupCommand(mockFlagProvider, mockKubernetesManager, mockFileReader, mockGitProviderFactory)

	bitbucketArgs := &command.ConfigBitbucketArgs{
		ClientId:  "testClientId",
		Secret:    "testSecret",
		Workspace: "testWorkspace",
		Repo:      "testRepo",
	}
	gitProvider := &command.ConfigGitProvider{
		Bitbucket: bitbucketArgs,
	}

	jobConfig := &external.CleanupJobConfig{}

	cleanupConfig := &command.CleanupConfig{
		Kubeconfig:  "kubeconfig",
		GitProvider: gitProvider,
		JobConfig:   jobConfig,
	}

	namespaceLabel := "testLabel"

	sut.CleanupConfig = cleanupConfig
	sut.NamespaceLabel = &namespaceLabel
	sut.Execute()

	if mockKubernetesManager.Called["getnamespaces"] != 1 {
		t.Errorf("GetNamespaces() should have been called once")
	}

	firstArgs := mockKubernetesManager.CalledWith["getnamespaces"][0]
	switch args := firstArgs.(type) {
	case *testutils.KMGetNamespacesArgs:
		flat := *args

		if flat.Label != namespaceLabel {
			t.Errorf("GetNamespaces() should have been called with Label: %s but got %+v", namespaceLabel, firstArgs)
		}
	}
}

func Test_ExecuteCreateCleanupJob(t *testing.T) {
	mockFileReader := testutils.NewMockFileReader("testFile")
	mockBitbucketManager := testutils.NewMockGitProvider()
	mockGitProviderFactory := &external.GitProviderFactory{
		BitbucketManager: mockBitbucketManager,
	}
	mockKubernetesManager := testutils.NewMockKubernetesManager()
	mockKubernetesManager.GetNamespacesRes = []string{"test-1"}
	mockBitbucketManager.GetBranchesRes = []string{"test-2"}
	mockFlagProvider := testutils.NewMockFlagProvider()
	sut := command.NewCleanupCommand(mockFlagProvider, mockKubernetesManager, mockFileReader, mockGitProviderFactory)

	bitbucketArgs := &command.ConfigBitbucketArgs{
		ClientId:  "testClientId",
		Secret:    "testSecret",
		Workspace: "testWorkspace",
		Repo:      "testRepo",
	}
	gitProvider := &command.ConfigGitProvider{
		Bitbucket: bitbucketArgs,
	}

	jobConfig := &external.CleanupJobConfig{
		Image:            "testImage",
		ImagePullSecret:  "testImagePullSecret",
		JobNamespace:     "testJobNamespace",
		ConnectionString: "testConnectionString",
	}

	cleanupConfig := &command.CleanupConfig{
		Kubeconfig:  "kubeconfig",
		GitProvider: gitProvider,
		JobConfig:   jobConfig,
	}

	namespaceLabel := "testLabel"

	sut.CleanupConfig = cleanupConfig
	sut.NamespaceLabel = &namespaceLabel
	sut.Execute()

	if mockKubernetesManager.Called["createcleanupjob"] != 1 {
		t.Errorf("CreateCleanupJob() should have been called once")
	}

	firstArgs := mockKubernetesManager.CalledWith["createcleanupjob"][0]
	switch args := firstArgs.(type) {
	case *testutils.KMCreateCleanupJobArgs:
		flat := *args.Config

		if flat.Name != "test-1" {
			t.Errorf("CreateCleanupJob() should have been called with Name: %s but got %+v", "test-1", firstArgs)
		}
		if flat.Image != jobConfig.Image {
			t.Errorf("CreateCleanupJob() should have been called with Image: %s but got %+v", jobConfig.Image, firstArgs)
		}
		if flat.ImagePullSecret != jobConfig.ImagePullSecret {
			t.Errorf("CreateCleanupJob() should have been called with ImagePullSecret: %s but got %+v", jobConfig.ImagePullSecret, firstArgs)
		}
		if flat.JobNamespace != jobConfig.JobNamespace {
			t.Errorf("CreateCleanupJob() should have been called with JobNamespace: %s but got %+v", jobConfig.JobNamespace, firstArgs)
		}
		if flat.ConnectionString != jobConfig.ConnectionString {
			t.Errorf("CreateCleanupJob() should have been called with ConnectionString: %s but got %+v", jobConfig.ConnectionString, firstArgs)
		}
	}
}
