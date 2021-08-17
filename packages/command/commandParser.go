package command

import (
	"os"

	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
)

type ICommand interface {
	GetFlags(external.IFlagProvider) (err error)
	FlagsValid() (err error)
	Execute(*GlobalCommandOptions) (err error)
}

type ICommandParser interface {
	ParseCommands() (err error)
}

type IBeforeOthers interface {
	BeforeOthers(*GlobalCommandOptions)
}

type IIsCurrentSubcommand interface {
	IsCurrentSubcommand() bool
}

type GlobalCommandOptions struct {
	Logger LoggerTypes
}

type CommandParser struct {
	flagProvider      external.IFlagProvider
	kubernetesManager external.IKubernetesManager
	gitProvider       *external.GitProviderFactory
	globals           *GlobalCommandOptions
	commands          []ICommand
}

func NewCommandParser(flagProvider external.IFlagProvider, gitProviderFactory *external.GitProviderFactory, kubernetesManager external.IKubernetesManager, postgresManager external.IPostgresManager) *CommandParser {
	parser := &CommandParser{
		flagProvider:      flagProvider,
		kubernetesManager: kubernetesManager,
		globals:           &GlobalCommandOptions{},
		commands: []ICommand{
			&LoggerCommand{},
			&CleanBranchCommand{},
			NewPRCommentCommand(gitProviderFactory),
			NewNamespaceCommand(kubernetesManager),
			NewDatabaseCommand(postgresManager),
		},
	}
	return parser
}

func (parser CommandParser) ParseCommands() (err error) {

	if len(os.Args) < 1 || os.Args[0] == "" {
		return errors.New("Must specify a command to run")
	}

	parser.flagProvider.Parse()

	// Some commands are more global (Logger), these command flags should be added to all subcommands and might have extra logic that needs to happen before we execute other commands.
	for _, c := range parser.commands {
		// If a command doesn't have a subcommand then it must be global.
		if _, ok := c.(IIsCurrentSubcommand); !ok {
			if err := c.GetFlags(parser.flagProvider); err != nil {
				return errors.Wrapf(err, "Failed to get flags for command: %T", c)
			}
		}

		if ft, ok := c.(IBeforeOthers); ok {
			ft.BeforeOthers(parser.globals)
		}
	}

	for _, c := range parser.commands {
		if ft, ok := c.(IIsCurrentSubcommand); ok {
			if ft.IsCurrentSubcommand() {
				if err := c.GetFlags(parser.flagProvider); err != nil {
					return errors.Wrapf(err, "Failed to get flags for command: %T", c)
				}

				if err = c.FlagsValid(); err != nil {
					return errors.Wrap(err, "Failed to validate arguments")
				}

				if err = c.Execute(parser.globals); err != nil {
					return errors.Wrap(err, "Failed to execute command")
				}
			}
		}
	}

	return
}
