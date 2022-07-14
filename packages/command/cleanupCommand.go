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
	gitProviderFactory *external.GitProviderFactory
	kubernetesManager  external.IKubernetesManager
	cmd                external.IFlagSet

	cleanupConfigPath string
	NamespaceLabel    *string
	CleanupConfig     *CleanupConfig
	fileReader        external.IFileReader
}

func NewCleanupCommand(flagProvider external.IFlagProvider, kubernetesManager external.IKubernetesManager, FileReader external.IFileReader, gitProviderFactory *external.GitProviderFactory) *CleanupCommand {
	return &CleanupCommand{
		fileReader:         FileReader,
		gitProviderFactory: gitProviderFactory,
		kubernetesManager:  kubernetesManager,
		cmd:                flagProvider.NewFlagSet("Cleanup", "Compares open Pull Requests with namespaces in the cluster and cleanup extras, Usage: Cleanup <CleanupConfigPath>"),
	}
}

func (c *CleanupCommand) IsCurrent() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "Cleanup")
}

type CleanupConfig struct {
	Kubeconfig  string                     `yaml:"kubeconfig"`
	GitProvider *ConfigGitProvider         `yaml:"gitProvider,omitempty"`
	JobConfig   *external.CleanupJobConfig `yaml:"job,omitempty"`
}

type ConfigGitProvider struct {
	Bitbucket *ConfigBitbucketArgs `yaml:"bitbucket,omitempty"`
	Github    *ConfigGithubArgs    `yaml:"github,omitempty"`
}

type ConfigBitbucketArgs struct {
	ClientId  string `yaml:"clientId"`
	Secret    string `yaml:"secret"`
	Workspace string `yaml:"workspace"`
	Repo      string `yaml:"repo"`
}

type ConfigGithubArgs struct {
	Organization string `yaml:"organization"`
	Repo         string `yaml:"repo"`
	Token        string `yaml:"token"`
	Username     string `yaml:"username"`
}

func (c *CleanupCommand) GetFlags() (err error) {
	c.NamespaceLabel = c.cmd.String("NamespaceLabel", "dev.centeva.meta=PullRequest", "Set the label used to check if a namespace can be cleaned up")

	if len(os.Args) <= 2 || os.Args[2] == "" {
		c.cmd.PrintDefaults()
		return errors.New("Cleanup requires a cleanupConfig file, check usage.")
	}

	c.cleanupConfigPath = os.Args[2]
	c.cmd.Parse(os.Args[3:])

	c.CleanupConfig, err = readConfigFile(c.fileReader, c.cleanupConfigPath)

	if err != nil {
		return errors.Wrap(err, "Failed to read config")
	}

	return
}

func readConfigFile(fileReader external.IFileReader, path string) (config *CleanupConfig, err error) {
	file, err := fileReader.ReadFile(path)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to read config file")
	}

	if err = yaml.Unmarshal(file, &config); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal config file")
	}

	return
}

func (c *CleanupCommand) Execute() (err error) {

	var branchesRaw []string

	switch {
	case c.CleanupConfig.GitProvider.Bitbucket != nil:
		{
			config := c.CleanupConfig.GitProvider.Bitbucket

			if _, err := c.gitProviderFactory.BitbucketManager.BasicAuth(config.ClientId, config.Secret); err != nil {
				return errors.Wrap(err, "Failed to auth")
			}

			if branchesRaw, err = c.gitProviderFactory.BitbucketManager.GetOpenPRBranches(config.Workspace, config.Repo); err != nil {
				return errors.Wrap(err, "Failed to get branches")
			}
		}
	case c.CleanupConfig.GitProvider.Github != nil:
		{
			config := c.CleanupConfig.GitProvider.Github

			c.gitProviderFactory.GithubManager.BasicAuth(config.Username, config.Token)

			if branchesRaw, err = c.gitProviderFactory.GithubManager.GetOpenPRBranches(config.Organization, config.Repo); err != nil {
				return errors.Wrap(err, "Failed to get branches")
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

	if c.CleanupConfig != nil && c.CleanupConfig.Kubeconfig != "" {
		if _, err := c.kubernetesManager.OutClusterConfig(context.Background(), c.CleanupConfig.Kubeconfig); err != nil {
			return errors.Wrap(err, "Failed to create outCluster config")
		}
	} else {
		if _, err := c.kubernetesManager.InClusterConfig(context.Background()); err != nil {
			return errors.Wrap(err, "Failed to create inCluster config")
		}
	}
	if namespaces, err = c.kubernetesManager.GetNamespaces(*c.NamespaceLabel); err != nil {
		return errors.Wrap(err, "Failed to get namespaces")
	}

	var cleanupList []string

	for _, name := range namespaces {
		if !Contains(branches, name) {
			cleanupList = append(cleanupList, name)
		}
	}

	log.Printf("Cleaning up %s", cleanupList)

	jobConfig := c.CleanupConfig.JobConfig
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
