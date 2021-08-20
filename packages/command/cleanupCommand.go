package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type CleanupCommand struct {
	cleanupConfigPath  string
	namespaceLabel     *string
	cleanupConfig      *CleanupConfig
	fileReader         external.IFileReader
	gitProviderFactory *external.GitProviderFactory
	kubernetesManager  external.IKubernetesManager
}

func NewCleanupCommand(FileReader external.IFileReader, gitProviderFactory *external.GitProviderFactory, kubernetesManager external.IKubernetesManager) *CleanupCommand {
	return &CleanupCommand{
		fileReader:         FileReader,
		gitProviderFactory: gitProviderFactory,
		kubernetesManager:  kubernetesManager,
	}
}

func (c *CleanupCommand) IsCurrentSubcommand() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "Cleanup")
}

type CleanupConfig struct {
	Kubeconfig  string                     `yaml:"kubeconfig"`
	GitProvider *configGitProvider         `yaml:"gitProvider,omitempty"`
	JobConfig   *external.CleanupJobConfig `yaml:"job,omitempty"`
}

type configGitProvider struct {
	Bitbucket *configBitbucketArgs `yaml:"bitbucket,omitempty"`
}

type configBitbucketArgs struct {
	ClientId  string `yaml:"clientId"`
	Secret    string `yaml:"secret"`
	Workspace string `yaml:"workspace"`
	Repo      string `yaml:"repo"`
}

func (c *CleanupCommand) GetFlags(flagProvider external.IFlagProvider) (err error) {
	cmd := flagProvider.NewFlagSet("Cleanup", "Compares open Pull Requests with namespaces in the cluster and cleanup extras.")
	if len(os.Args) <= 2 || os.Args[2] == "" {
		return errors.New("Cleanup requires a cleanupConfig file")
	}
	c.cleanupConfigPath = os.Args[2]

	c.namespaceLabel = cmd.String("NamespaceLabel", "dev.centeva.meta=PullRequest", "Set the label used to check if a namespace can be cleaned up")

	cmd.Parse(os.Args[3:])

	c.cleanupConfig, err = c.ReadConfigFile(c.fileReader, c.cleanupConfigPath)

	if err != nil {
		return errors.Wrap(err, "Failed to read config")
	}

	return
}

func (k *CleanupCommand) FlagsValid() (err error) {
	if k.cleanupConfigPath == "" {
		return errors.New("CleanupConfig is required")
	}

	return
}

func (c *CleanupCommand) ReadConfigFile(fileReader external.IFileReader, path string) (config *CleanupConfig, err error) {
	file, err := fileReader.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to read config file")
	}

	if err = yaml.Unmarshal(file, &config); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal config file")
	}

	return
}

func (c *CleanupCommand) Execute(globals *GlobalCommandOptions) (err error) {

	var branchesRaw []string

	switch {
	case c.cleanupConfig.GitProvider.Bitbucket != nil:
		{
			config := c.cleanupConfig.GitProvider.Bitbucket

			if _, err := c.gitProviderFactory.BitbucketManager.BasicAuth(config.ClientId, config.Secret); err != nil {
				return errors.Wrap(err, "Failed to auth")
			}

			if branchesRaw, err = c.gitProviderFactory.BitbucketManager.GetOpenPRBranches(config.Workspace, config.Repo); err != nil {
				return errors.Wrap(err, "Faield to get branches")
			}
		}
	default:
		return errors.New("No gitprovider found in configfile")
	}

	var branches []string

	for _, name := range branchesRaw {
		branches = append(branches, CleanBranch(name))
	}

	var namespaces []string

	if c.cleanupConfig != nil && c.cleanupConfig.Kubeconfig != "" {
		if _, err := c.kubernetesManager.OutClusterConfig(context.Background(), c.cleanupConfig.Kubeconfig); err != nil {
			return errors.Wrap(err, "Failed to create outCluster config")
		}
	} else {
		if _, err := c.kubernetesManager.InClusterConfig(context.Background()); err != nil {
			return errors.Wrap(err, "Failed to create inCluster config")
		}
	}
	if namespaces, err = c.kubernetesManager.GetNamespaces(*c.namespaceLabel); err != nil {
		return errors.Wrap(err, "Failed to get namespaces")
	}

	var cleanupList []string

	for _, name := range namespaces {
		if !Contains(branches, name) {
			cleanupList = append(cleanupList, name)
		}
	}

	log.Printf("Cleaning up %s", cleanupList)

	jobConfig := c.cleanupConfig.JobConfig
	var wg sync.WaitGroup
	wg.Add(len(cleanupList))

	createJob := func(name string, errs chan error) {
		defer wg.Done()

		if err := c.kubernetesManager.CreateCleanupJob(&external.CleanupJobConfig{
			Name:             name,
			Image:            jobConfig.Image,
			ImagePullSecret:  jobConfig.ImagePullSecret,
			JobNamespace:     jobConfig.JobNamespace,
			ConnectionString: jobConfig.ConnectionString,
		}); err != nil {
			errs <- errors.Wrap(err, "Failed to create cleanupJob")
		}
	}

	errs := make(chan error, len(cleanupList))

	for _, name := range cleanupList {
		go createJob(name, errs)
	}

	wg.Wait()
	close(errs)

	var allErrs []string
	for goErrs := range errs {
		if goErrs != nil {
			allErrs = append(allErrs, fmt.Sprintf(" %d) %s", len(allErrs), goErrs))
		}
	}

	if len(allErrs) > 0 {
		allErr := errors.New(strings.Join(allErrs, ""))
		return errors.Wrap(allErr, "jobs failed")
	}

	return
}

func Contains(arr []string, str string) bool {
	for _, val := range arr {
		if val == str {
			return true
		}
	}

	return false
}
