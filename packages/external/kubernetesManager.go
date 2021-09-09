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

func buildCleanupJob(name string, config *CleanupJobConfig) *batchv1.Job {
	ttlSecondsAfterFinished := int32(1)

	return &batchv1.Job{
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
							Name:            "deletedb",
							Image:           config.Image,
							ImagePullPolicy: v1.PullAlways,
							Args:            []string{"DeleteDatabase", config.Name, fmt.Sprintf(`--ConnectionString=%s`, config.ConnectionString)},
						},
						{
							Name:            "deletenamespace",
							Image:           config.Image,
							ImagePullPolicy: v1.PullAlways,
							Args:            []string{"DeleteNamespace", config.Name},
						},
					},
				},
			},
		},
	}
}

func (k *KubernetesManager) CreateCleanupJob(config *CleanupJobConfig) (err error) {
	if config.Timeout == "" {
		config.Timeout = "2m"
	}

	name := fmt.Sprintf("cleanup-%s", config.Name)
	cleanupJob := buildCleanupJob(name, config)

	dur, err := time.ParseDuration(config.Timeout)
	if err != nil {
		return errors.Wrap(err, "Failed to parse Timeout")
	}

	_, err = k.clientset.BatchV1().Jobs(config.JobNamespace).Create(k.context, cleanupJob, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "Failed to create cleanup job")
	}

	durSeconds := int64(dur.Seconds())

	watcher, err := k.clientset.BatchV1().Jobs(config.JobNamespace).Watch(k.context, metav1.ListOptions{
		LabelSelector:  fmt.Sprintf("job-name=%s", name),
		TimeoutSeconds: &durSeconds,
	})

	if err != nil {
		return errors.Wrapf(err, "Failed to create watcher: %+v", watcher)
	}

	for event := range watcher.ResultChan() {
		job := event.Object.(*batchv1.Job)

		switch event.Type {
		case watch.Modified:
			ok, err := checkJobCompleted(job)

			if err != nil {
				return errors.Wrapf(err, "Error watching job %s", job.Name)
			}

			if ok {
				log.Printf("Job '%s' finished sucessfully", job.Name)

				deletePolicy := metav1.DeletePropagationForeground
				err = k.clientset.BatchV1().Jobs(config.JobNamespace).Delete(k.context, name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})

				if err != nil {
					return errors.Wrap(err, "Failed to delete job")
				}

				return nil
			}
		}
	}

	return errors.Errorf("Timeout error watching job %s", name)
}

func checkJobCompleted(job *batchv1.Job) (res bool, err error) {

	if job == nil {
		return false, errors.New("Job is nil")
	}

	if len(job.Status.Conditions) > 0 {
		latest := job.Status.Conditions[len(job.Status.Conditions)-1]
		switch latest.Type {
		case batchv1.JobComplete:
			return true, nil
		default:
			return false, errors.Errorf("Job %T: %s Details: %s", latest.Type, latest.Reason, latest.Message)
		}
	}

	return false, nil
}
