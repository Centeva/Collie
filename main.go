package main

import (
	"log"

	"bitbucket.org/centeva/collie/packages/command"
	"bitbucket.org/centeva/collie/packages/external"
)

func main() {
	log.SetFlags(0)

	flagProvider := external.NewFlagProvider()
	bitbucketManager := external.NewBitbucketManager()
	gitProviderFactory := external.NewGitProviderFactory(bitbucketManager)
	kubernetesManager := &external.KubernetesManager{}
	postgresManager := external.NewPostgresManager()
	fileReader := &external.FileReader{}

	cmd := command.NewCommandParser(flagProvider, gitProviderFactory, kubernetesManager, postgresManager, fileReader)

	err := cmd.ParseCommands()

	if err != nil {
		panic(err)
	}
}
