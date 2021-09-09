package command

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
)

type NamespaceCommand struct {
	kubernetesManager external.IKubernetesManager
	ctx               context.Context
	ctxCancel         context.CancelFunc
	cmd               external.IFlagSet
	Namespace         string
	Kubeconfig        *string
	Timeout           *string
}

func NewNamespaceCommand(flagProvider external.IFlagProvider, kubernetesManager external.IKubernetesManager) *NamespaceCommand {
	return &NamespaceCommand{
		kubernetesManager: kubernetesManager,
		cmd:               flagProvider.NewFlagSet("DeleteNamespace", "Delete kubernetes namespace and everything in it Usage: DeleteNamespace <namespace>"),
	}
}

func (k *NamespaceCommand) IsCurrent() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "DeleteNamespace")
}

func (k *NamespaceCommand) GetFlags() (err error) {
	k.Timeout = k.cmd.String("Timeout", "10m", "Context Timout")
	k.Kubeconfig = k.cmd.String("Kubeconfig", "", "Path to kubeconfig context file, used for running outside of the cluster")

	if len(os.Args) <= 2 || os.Args[2] == "" {
		k.cmd.PrintDefaults()
		return errors.New("Missing namespace, see usage.")
	}
	k.Namespace = os.Args[2]

	k.cmd.Parse(os.Args[3:])
	return
}

func (k *NamespaceCommand) CreateContext() (ctx context.Context, err error) {
	timeout, err := time.ParseDuration(*k.Timeout)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse timeout: %s", *k.Timeout)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	k.ctx = ctx
	k.ctxCancel = cancel
	return
}

func (k *NamespaceCommand) Execute() (err error) {
	if _, err := k.CreateContext(); err != nil {
		return errors.Wrap(err, "Failed to create context")
	}

	if k.Kubeconfig != nil && *k.Kubeconfig != "" {
		if _, err := k.kubernetesManager.OutClusterConfig(k.ctx, *k.Kubeconfig); err != nil {
			return errors.Wrap(err, "Failed to create outCluster config")
		}
	} else {
		if _, err := k.kubernetesManager.InClusterConfig(k.ctx); err != nil {
			return errors.Wrap(err, "Failed to create inCluster config, are you running in a cluster?")
		}

	}

	list, err := k.kubernetesManager.GetNamespaces("")
	match := false
	for _, v := range list {
		match = match || v == k.Namespace
	}

	if !match {
		log.Printf("Namespace %s does not exist; Nothing to delete", k.Namespace)
		return
	}

	if err != nil {
		return errors.Wrapf(err, "Failed to get namespaces")
	}

	if err := k.kubernetesManager.DeleteNamespace(k.Namespace); err != nil {
		return errors.Wrapf(err, "Failed to delete namespace %s", k.Namespace)
	}

	log.Printf("Deleted namespace %s", k.Namespace)

	return
}
