package testutils

import (
	"context"

	"bitbucket.org/centeva/collie/packages/external"
	"k8s.io/client-go/kubernetes"
)

type MockKubernetesManager struct {
	Called           map[string]int
	CalledWith       map[string][]interface{}
	GetNamespacesRes []string
}

func NewMockKubernetesManager() *MockKubernetesManager {
	return &MockKubernetesManager{
		Called:     make(map[string]int),
		CalledWith: make(map[string][]interface{}),
	}
}

type KMInClusterConfigArgs struct {
	Context context.Context
}

func (m *MockKubernetesManager) InClusterConfig(context context.Context) (client *kubernetes.Clientset, err error) {
	m.Called["inclusterconfig"]++
	m.CalledWith["inclusterconfig"] = append(m.CalledWith["inclusterconfig"], &KMInClusterConfigArgs{context})
	return
}

type KMOutClusterConfigArgs struct {
	Context    context.Context
	Kubeconfig string
}

func (m *MockKubernetesManager) OutClusterConfig(context context.Context, kubeconfig string) (client *kubernetes.Clientset, err error) {
	m.Called["outclusterconfig"]++
	m.CalledWith["outclusterconfig"] = append(m.CalledWith["outclusterconfig"], &KMOutClusterConfigArgs{context, kubeconfig})
	return
}

type KMDeleteNamespaceArgs struct {
	Namespace string
}

func (m *MockKubernetesManager) DeleteNamespace(namespace string) (err error) {
	m.Called["deletenamespace"]++
	m.CalledWith["deletenamespace"] = append(m.CalledWith["deletenamespace"], &KMDeleteNamespaceArgs{namespace})
	return
}

type KMGetNamespacesArgs struct {
	Label string
}

func (m *MockKubernetesManager) GetNamespaces(label string) (namespaces []string, err error) {
	m.Called["getnamespaces"]++
	m.CalledWith["getnamespaces"] = append(m.CalledWith["getnamespaces"], &KMGetNamespacesArgs{label})
	return m.GetNamespacesRes, nil
}

type KMCreateCleanupJobArgs struct {
	Config *external.CleanupJobConfig
}

func (m *MockKubernetesManager) CreateCleanupJob(config *external.CleanupJobConfig) (err error) {
	m.Called["createcleanupjob"]++
	m.CalledWith["createcleanupjob"] = append(m.CalledWith["createcleanupjob"], &KMCreateCleanupJobArgs{config})
	return
}
