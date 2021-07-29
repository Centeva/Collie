package lib

import (
	"log"
)

func Entry(cmd ICommandParser) (err error) {
	log.SetFlags(0)
	err = cmd.ParseCommands()
	return err
}
