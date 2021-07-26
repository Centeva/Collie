package lib

import (
	"flag"
	"log"
)

type ICommandParser interface {
	GetBranch() *string
	ParseFlags()
}

type CommandParser struct{}

func (c CommandParser) GetBranch() *string {
	return flag.String("CleanBranch", "", "Name of a branch to format")
}

func (c CommandParser) ParseFlags() {
	flag.Parse()
}

func Entry(cmd ICommandParser) {
	log.SetFlags(0)
	cleanedBranch := cmd.GetBranch()
	cmd.ParseFlags()

	if *cleanedBranch != "" {
		log.Printf("%s", CleanBranch(*cleanedBranch))
	}
}
