package external

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesManager struct {
	clientset *kubernetes.Clientset
	context   context.Context
}

type IKubernetesManager interface {
	InClusterConfig(context context.Context) (client *kubernetes.Clientset, err error)
	OutClusterConfig(context context.Context, kubeconfig string) (client *kubernetes.Clientset, err error)
	DeleteNamespace(namespace string) (err error)
	GetNamespaces() (namespaces []string, err error)
}

func (k *KubernetesManager) InClusterConfig(context context.Context) (client *kubernetes.Clientset, err error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get config from k8s api")
	}

	client, err = kubernetes.NewForConfig(config)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create client")
	}

	k.clientset = client
	k.context = context
	return
}

func (k *KubernetesManager) OutClusterConfig(context context.Context, kubeconfig string) (client *kubernetes.Clientset, err error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to build config from file: %s", kubeconfig)
	}

	client, err = kubernetes.NewForConfig(config)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create client")
	}

	k.clientset = client
	k.context = context
	return
}

func (k *KubernetesManager) DeleteNamespace(namespace string) (err error) {

	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	if err := k.clientset.CoreV1().Namespaces().Delete(k.context, namespace, *opts); err != nil {
		return errors.Wrap(err, "Failed to delete namespace")
	}
	return
}

func (k *KubernetesManager) GetNamespaces() (namespaces []string, err error) {

	res, err := k.clientset.CoreV1().Namespaces().List(k.context, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get namespaces")
	}

	namespaces = []string{}

	for _, item := range (*res).Items {
		namespaces = append(namespaces, item.Name)
	}

	return
}
