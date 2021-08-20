package command

import (
	"log"
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
	GetFlags(external.IFlagProvider) (err error)
	FlagsValid() (err error)
	Execute(*GlobalCommandOptions) (err error)
}

type ICommandParser interface {
	ParseCommands() (err error)
}

// Should look at removing this.
type IBeforeExecute interface {
	BeforeExecute(*GlobalCommandOptions) (err error)
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
	globals           *GlobalCommandOptions
	commands          []ICommand
}

func NewCommandParser(flagProvider external.IFlagProvider, gitProviderFactory *external.GitProviderFactory, kubernetesManager external.IKubernetesManager, postgresManager external.IPostgresManager, fileReader external.IFileReader) *CommandParser {
	parser := &CommandParser{
		flagProvider:      flagProvider,
		kubernetesManager: kubernetesManager,
		globals:           &GlobalCommandOptions{},
		commands: []ICommand{
			&CleanBranchCommand{},
			NewPRCommentCommand(gitProviderFactory),
			NewNamespaceCommand(kubernetesManager),
			NewDatabaseCommand(postgresManager),
			NewCleanupCommand(fileReader, gitProviderFactory, kubernetesManager),
		},
	}
	return parser
}

func (parser CommandParser) ParseCommands() (err error) {

	if len(os.Args) < 1 || os.Args[0] == "" {
		return errors.New("Must specify a command to run")
	}

	// Some commands are more global (Logger), these command flags should be added to all subcommands and might have extra logic that needs to happen before we execute other commands.
	for _, c := range parser.commands {
		// If a command doesn't have a subcommand then it must be global.
		if _, ok := c.(IIsCurrentSubcommand); !ok {
			log.Printf("globalCommand: %T", c)
			if err := c.GetFlags(parser.flagProvider); err != nil {
				return errors.Wrapf(err, "Failed to get flags for command: %T", c)
			}
		}

		parser.flagProvider.Parse()

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

				if ft, ok := c.(IBeforeExecute); ok {
					if err := ft.BeforeExecute(parser.globals); err != nil {
						return errors.Wrap(err, "Failed BeforeOthers")
					}
				}

				if err = c.Execute(parser.globals); err != nil {
					return errors.Wrap(err, "Failed to execute command")
				}
			}
		}
	}

	return
}
