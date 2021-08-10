package main

import (
	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/packages/core"
	"bitbucket.org/centeva/collie/packages/external"
)

func main() {
	cmd := command.NewCommandParser(&external.FlagProvider{})
	err := core.Entry(cmd)
	if err != nil {
		panic(err)
	}
}
