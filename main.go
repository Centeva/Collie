package main

import (
	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/packages/core"
	"bitbucket.org/centeva/collie/packages/external"
)

func main() {
	gitProviderFactory := &external.GitProviderFactory{
		BitbucketManager: external.NewBitbucketManager(),
	}
	cmd := command.NewCommandParser(&external.FlagProvider{}, gitProviderFactory)
	err := core.Entry(cmd)
	if err != nil {
		panic(err)
	}
}
