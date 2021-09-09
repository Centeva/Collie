package command

import (
	"log"
	"os"
	"strings"

	"bitbucket.org/centeva/collie/packages/external"
)

type HelpCommand struct {
	flagProvider external.IFlagProvider
	cmd          external.IFlagSet
}

func NewHelpCommand(flagProvider external.IFlagProvider) *HelpCommand {
	return &HelpCommand{
		flagProvider: flagProvider,
		cmd:          flagProvider.NewFlagSet("Help", "Collie is a cli tool full of useful devops commands! 🐶"),
	}
}

func (h *HelpCommand) IsCurrent() bool {
	return len(os.Args) > 1 && strings.EqualFold(os.Args[1], "Help")
}

func (h *HelpCommand) GetFlags() (err error) {
	return
}

func (h *HelpCommand) Execute() (err error) {

	for name, usage := range h.flagProvider.GetUsage() {
		log.Printf("%s:", name)
		log.Printf("\t%s", usage)
	}

	return
}
