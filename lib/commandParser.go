package lib

type ICommand interface {
	GetFlags(IFlagProvider)
	FlagsValid() bool
	Execute(*GlobalCommandOptions) (err error)
}

type ICommandParser interface {
	ParseCommands() (err error)
}

type IBeforeOthers interface {
	BeforeOthers(*GlobalCommandOptions)
}

type GlobalCommandOptions struct {
	Logger LoggerTypes
}

type CommandParser struct {
	flagProvider IFlagProvider
	globals      *GlobalCommandOptions
	commands     []ICommand
}

func NewCommandParser(flagProvider IFlagProvider) *CommandParser {
	parser := &CommandParser{
		flagProvider: flagProvider,
		globals:      &GlobalCommandOptions{},
		commands: []ICommand{
			&LoggerCommand{},
			&CleanBranchCommand{},
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
