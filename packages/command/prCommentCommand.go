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
	GitProvider        string
	GitSource          IGitSource
}

func NewPRCommentCommand(gitProviderFactory *external.GitProviderFactory) *PRCommentCommand {
	return &PRCommentCommand{
		gitProviderFactory: gitProviderFactory,
	}
}

func (c *PRCommentCommand) GetFlags(FlagProvider external.IFlagProvider) (err error) {
	cmd := FlagProvider.NewFlagSet("Comment", "Create a comment on a pull request Usage: Comment <GitProvider:[bitbucket]> <Args>")
	if len(os.Args) <= 2 || os.Args[2] == "" {
		return errors.New("Comment must have a GitProvider, must be [bitbucket]")
	}
	c.GitProvider = os.Args[2]

	switch c.GitProvider {
	case "bitbucket":
		otherSource := &BitBucketSource{
			Branch:    cmd.String("Branch", "", "(required) Source branch of the Pull Request"),
			ClientId:  cmd.String("ClientId", "", "(required) BitBucket OAuth ClientId/key"),
			Comment:   cmd.String("Comment", "", "(required) Comment message to add to the Pull Request"),
			Repo:      cmd.String("Repo", "", "(required) Repository name"),
			Secret:    cmd.String("Secret", "", "(required) BitBucket OAuth Secret"),
			Workspace: cmd.String("Workspace", "", "(required) BitBucket workspace"),
			Username:  cmd.String("Username", "", "Optional Username of comment author"),
			Password:  cmd.String("Password", "", "Optional Password of comment author"),
		}
		c.GitSource = otherSource
		cmd.Parse(os.Args[3:])
	}
	return
}

func (c *PRCommentCommand) IsCurrentSubcommand() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "Comment")
}

func (c *PRCommentCommand) FlagsValid() (err error) {
	switch s := c.GitSource.(type) {
	case *BitBucketSource:
		log.Printf("source: %+v", s)
		if s.Branch == nil || *s.Branch == "" {
			return errors.New("Branch is required")
		}

		if s.ClientId == nil || *s.ClientId == "" {
			return errors.New("ClientId is required")
		}

		if s.Comment == nil || *s.Comment == "" {
			return errors.New("Comment is required")
		}

		if s.Repo == nil || *s.Repo == "" {
			return errors.New("Repo is required")
		}

		if s.Secret == nil || *s.Secret == "" {
			return errors.New("Secret is required")
		}

		if s.Workspace == nil || *s.Workspace == "" {
			return errors.New("Workspace is required")
		}

		// if (s.Username != nil && *s.Username != "") == (s.Password != nil && *s.Password != "") {
		// 	return errors.New("Both Username and Password are required to set comment author")
		// }

	default:
		return errors.New("Could not recognize GitProvider")
	}

	return
}

func (c *PRCommentCommand) Execute(globals *GlobalCommandOptions) (err error) {

	switch s := c.GitSource.(type) {
	case *BitBucketSource:
		{
			if s.Username != nil && *s.Username != "" {
				if _, err := c.gitProviderFactory.BitbucketManager.UserAuth(*s.ClientId, *s.Secret, *s.Username, *s.Password); err != nil {
					return errors.Wrap(err, "Failed to authenticate with bitbucket api while executing UserAuth")
				}
			} else {
				if _, err := c.gitProviderFactory.BitbucketManager.BasicAuth(*s.ClientId, *s.Secret); err != nil {
					return errors.Wrap(err, "Failed to authenticate with bitbucket api while executing BasicAuth")
				}
			}

			if err := c.gitProviderFactory.BitbucketManager.Comment(*s.Workspace, *s.Repo, *s.Branch, *s.Comment); err != nil {
				return errors.Wrap(err, "Failed to add comment through bitbucket api")
			}
		}
	}

	log.Printf("Added comment to pull request")
	return
}
