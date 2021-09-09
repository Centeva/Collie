package command

import (
	"os"

	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
)

type LoggerTypes string

const (
	CLI      LoggerTypes = "cli"
	TEAMCITY LoggerTypes = "teamcity"
)

type ICommand interface {
	GetFlags() (err error)
	Execute() (err error)
	IsCurrent() bool
}

type CommandParser struct {
	flagProvider      external.IFlagProvider
	kubernetesManager external.IKubernetesManager
	commands          []ICommand
}

func NewCommandParser(flagProvider external.IFlagProvider, gitProviderFactory *external.GitProviderFactory, kubernetesManager external.IKubernetesManager, postgresManager external.IPostgresManager, fileReader external.IFileReader) *CommandParser {
	parser := &CommandParser{
		flagProvider:      flagProvider,
		kubernetesManager: kubernetesManager,
		commands: []ICommand{
			NewCleanBranchCommand(flagProvider),
			NewPRCommentCommand(flagProvider, gitProviderFactory),
			NewNamespaceCommand(flagProvider, kubernetesManager),
			NewDatabaseCommand(flagProvider, postgresManager),
			NewCleanupCommand(flagProvider, kubernetesManager, fileReader, gitProviderFactory),
			NewHelpCommand(flagProvider),
		},
	}
	return parser
}

func (parser CommandParser) ParseCommands() (err error) {

	if len(os.Args) < 1 || os.Args[0] == "" {
		return errors.New("Must specify a command to run")
	}

	for _, c := range parser.commands {
		if c.IsCurrent() {
			if err := c.GetFlags(); err != nil {
				return errors.Wrapf(err, "Failed to get flags for command: %T", c)
			}

			if err = c.Execute(); err != nil {
				return errors.Wrap(err, "Failed to execute command")
			}

			return
		}
	}

	help := NewHelpCommand(parser.flagProvider)

	if err = help.Execute(); err != nil {
		return errors.Wrap(err, "Failed to output help")
	}

	return
}
