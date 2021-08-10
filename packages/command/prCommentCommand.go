package command

import (
	"bitbucket.org/centeva/collie/packages/external"
)

type PRCommentCommand struct {
	Comment string `tc:"comment"`
}

func (c *PRCommentCommand) GetFlags(FlagProvider external.IFlagProvider) {
	cmd := FlagProvider.NewFlagSet("Comment")
	c.Comment = cmd.Arg(0)
}

func (c *PRCommentCommand) FlagsValid() bool {
	return c.Comment != ""
}

func (c *PRCommentCommand) Execute(globals *GlobalCommandOptions) (err error) {

	return
}
