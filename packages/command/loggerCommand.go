package command

import (
	"bitbucket.org/centeva/collie/packages/external"
	"github.com/pkg/errors"
)

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

func (c *LoggerCommand) GetFlags(flagProvider external.IFlagProvider) (err error) {
	flagProvider.StringVar(&c.Logger, "Logger", string(CLI), "Log output style to use [cli,teamcity]")
	return
}

func (c *LoggerCommand) FlagsValid() (err error) {
	if c.Logger == "" {
		return errors.New("Logger must have a value")
	}

	if c.Logger != "cli" && c.Logger != "teamcity" {
		return errors.New("Logger must be either 'cli' or 'teamcity'")
	}

	return
}

func (c *LoggerCommand) BeforeOthers(globals *GlobalCommandOptions) {
	globals.Logger = (LoggerTypes)(c.Logger)
}

func (c *LoggerCommand) Execute(globals *GlobalCommandOptions) (err error) {
	return
}
