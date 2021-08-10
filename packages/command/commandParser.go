package command

import (
	"bitbucket.org/centeva/collie/packages/external"
)

type ICommand interface {
	GetFlags(external.IFlagProvider)
	FlagsValid() bool
	Execute(*GlobalCommandOptions) (err error)
}

type ICommandParser interface {
	ParseCommands() (err error)
}

type IBeforeOthers interface {
	BeforeOthers(*GlobalCommandOptions)
}

type CommandParser struct {
	flagProvider external.IFlagProvider
	globals      *GlobalCommandOptions
	commands     []ICommand
}

func NewCommandParser(flagProvider external.IFlagProvider) *CommandParser {
	parser := &CommandParser{
		flagProvider: flagProvider,
		globals:      &GlobalCommandOptions{},
		commands: []ICommand{
			&LoggerCommand{},
			&CleanBranchCommand{},
			&PRCommentCommand{},
		},
	}
	return parser
}

func (parser CommandParser) ParseCommands() (err error) {
	for _, command := range parser.commands {
		command.GetFlags(parser.flagProvider)
	}
	parser.flagProvider.Parse()

	for _, command := range parser.commands {
		if ft, ok := command.(IBeforeOthers); ok {
			ft.BeforeOthers(parser.globals)
		}
	}

	for _, c := range parser.commands {
		if c.FlagsValid() {
			err = c.Execute(parser.globals)
			if err != nil {
				return
			}
		}
	}

	return
}
