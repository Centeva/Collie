package lib

type LoggerTypes string

const (
	CLI      LoggerTypes = "cli"
	TEAMCITY LoggerTypes = "teamcity"
)

type LoggerCommandOptions struct {
}

type LoggerCommand struct {
	Logger string
}

func (c *LoggerCommand) GetFlags(flagProvider IFlagProvider) {
	flagProvider.StringVar(&c.Logger, "Logger", string(CLI), "Log output style to use [cli,teamcity]")
}

func (c *LoggerCommand) FlagsValid() bool {
	return c.Logger != ""
}

func (c *LoggerCommand) BeforeOthers(globals *GlobalCommandOptions) {
	globals.Logger = (LoggerTypes)(c.Logger)
}

func (c *LoggerCommand) Execute(globals *GlobalCommandOptions) (err error) {
	return
}
