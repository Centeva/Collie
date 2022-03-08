package command

import (
	"log"
	"os"
	"strings"

	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
)

type IGitSource interface{}

type BitBucketSource struct {
	Branch    *string
	ClientId  *string
	Comment   *string
	Repo      *string
	Secret    *string
	Workspace *string
	Username  *string
	Password  *string
}

type PRCommentCommand struct {
	gitProviderFactory *external.GitProviderFactory
	cmd                external.IFlagSet

	GitProvider string
	GitSource   IGitSource
}

func NewPRCommentCommand(flagProvider external.IFlagProvider, gitProviderFactory *external.GitProviderFactory) *PRCommentCommand {
	return &PRCommentCommand{
		gitProviderFactory: gitProviderFactory,
		cmd:                flagProvider.NewFlagSet("Comment", "Create a comment on a pull request Usage: Comment <GitProvider:<bitbucket>> <Args>"),
	}
}

func (c *PRCommentCommand) IsCurrent() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "Comment")
}

func (c *PRCommentCommand) GetFlags() (err error) {
	if len(os.Args) <= 2 || os.Args[2] == "" {
		c.cmd.PrintDefaults()
		return errors.New("Comment must have a GitProvider, must be <bitbucket>, check usage.")
	}
	c.GitProvider = os.Args[2]

	switch c.GitProvider {
	case "bitbucket":
		source := &BitBucketSource{
			Branch:    c.cmd.String("Branch", "", "(required) Source branch of the Pull Request"),
			ClientId:  c.cmd.String("ClientId", "", "(required) BitBucket OAuth ClientId/key"),
			Comment:   c.cmd.String("Comment", "", "(required) Comment message to add to the Pull Request"),
			Repo:      c.cmd.String("Repo", "", "(required) Repository name"),
			Secret:    c.cmd.String("Secret", "", "(required) BitBucket OAuth Secret"),
			Workspace: c.cmd.String("Workspace", "", "(required) BitBucket workspace"),
			Username:  c.cmd.String("Username", "", "Optional Username of comment author"),
			Password:  c.cmd.String("Password", "", "Optional Password of comment author"),
		}
		c.GitSource = source
		c.cmd.Parse(os.Args[3:])
		if err := c.ValidateBitbucketFlags(source); err != nil {
			return errors.Wrapf(err, "Failed to validate flags")
		}

	default:
		return errors.New("Could not recognize GitProvider")
	}
	return
}

func (c *PRCommentCommand) ValidateBitbucketFlags(source *BitBucketSource) error {
	if source.Branch == nil || *source.Branch == "" {
		return errors.New("Branch is required")
	}

	if source.ClientId == nil || *source.ClientId == "" {
		return errors.New("ClientId is required")
	}

	if source.Comment == nil || *source.Comment == "" {
		return errors.New("Comment is required")
	}

	if source.Repo == nil || *source.Repo == "" {
		return errors.New("Repo is required")
	}

	if source.Secret == nil || *source.Secret == "" {
		return errors.New("Secret is required")
	}

	if source.Workspace == nil || *source.Workspace == "" {
		return errors.New("Workspace is required")
	}

	return nil
}

func (c *PRCommentCommand) Execute() (err error) {

	switch s := c.GitSource.(type) {
	case *BitBucketSource:
		{
			if _, err := c.gitProviderFactory.BitbucketManager.BasicAuth(*s.ClientId, *s.Secret); err != nil {
				return errors.Wrap(err, "Failed to authenticate with bitbucket api while executing BasicAuth")
			}

			if err := c.gitProviderFactory.BitbucketManager.Comment(*s.Workspace, *s.Repo, *s.Branch, *s.Comment, s.Username, s.Password); err != nil {
				return errors.Wrap(err, "Failed to add comment through bitbucket api")
			}
		}
	}

	log.Printf("Added comment to pull request")
	return
}
