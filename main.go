package main

import "bitbucket.org/centeva/collie/lib"

func main() {
	cmd := lib.NewCommandParser(&lib.FlagProvider{})
	err := lib.Entry(cmd)
	if err != nil {
		panic(err)
	}
}
