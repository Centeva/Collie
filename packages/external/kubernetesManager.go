package external

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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
	GetNamespaces(label string) (namespaces []string, err error)
	CreateCleanupJob(config *CleanupJobConfig) (err error)
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

func (k *KubernetesManager) GetNamespaces(label string) (namespaces []string, err error) {

	res, err := k.clientset.CoreV1().Namespaces().List(k.context, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get namespaces")
	}

	namespaces = []string{}

	for _, item := range (*res).Items {
		namespaces = append(namespaces, item.Name)
	}

	return
}

type CleanupJobConfig struct {
	Image            string `yaml:"image"`
	ImagePullSecret  string `yaml:"imagePullSecret"`
	JobNamespace     string `yaml:"namespace"`
	ConnectionString string `yaml:"connectionString"`
	ServiceAccount   string `yaml:"serviceAccountName"`
	Timeout          string `yaml:"timeout"`
	Name             string
}

func (k *KubernetesManager) CreateCleanupJob(config *CleanupJobConfig) (err error) {
	if config.Timeout == "" {
		config.Timeout = "2m"
	}

	ttlSecondsAfterFinished := int32(120)
	name := fmt.Sprintf("cleanup-%s", config.Name)

	cleanupJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: config.JobNamespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					ImagePullSecrets: []v1.LocalObjectReference{
						{Name: config.ImagePullSecret},
					},
					Containers: []v1.Container{
						{
							Name:  "deletedb",
							Image: config.Image,
							Args:  []string{"DeleteDatabase", config.Name, fmt.Sprintf(`--ConnectionString=%s`, config.ConnectionString)},
						},
						{
							Name:  "deletenamespace",
							Image: config.Image,
							Args:  []string{"DeleteNamespace", config.Name},
						},
					},
				},
			},
		},
	}
	// log.Printf("Create job: %+v", cleanupJob)
	_, err = k.clientset.BatchV1().Jobs(config.JobNamespace).Create(k.context, cleanupJob, metav1.CreateOptions{})

	if err != nil {
		return errors.Wrap(err, "Failed to create cleanup job")
	}

	dur, err := time.ParseDuration(config.Timeout)

	if err != nil {
		return errors.Wrap(err, "Failed to parse Timeout")
	}

	durSeconds := int64(dur.Seconds())

	watcher, err := k.clientset.BatchV1().Jobs(config.JobNamespace).Watch(k.context, metav1.ListOptions{
		LabelSelector:  fmt.Sprintf("name=%s", name),
		TimeoutSeconds: &durSeconds,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to create watcher")
	}

	for event := range watcher.ResultChan() {
		job := event.Object.(*batchv1.Job)

		switch event.Type {
		case watch.Deleted:
			log.Printf("Job '%s' finished sucessfully", job.Name)
			return
		}
	}

	return errors.Errorf("Timeout error watching job %s", name)
}
