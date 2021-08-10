package core

import (
	"log"

	"bitbucket.org/centeva/collie/packages/command"
)

func Entry(cmd command.ICommandParser) (err error) {
	log.SetFlags(0)
	err = cmd.ParseCommands()
	return err
}
