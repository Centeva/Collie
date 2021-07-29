package lib

import (
	"flag"
)

type IFlagProvider interface {
	StringVar(p *string, name string, value string, usage string)
	Parse()
}

type FlagProvider struct{}

func (f *FlagProvider) StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, value, usage)
}

func (f *FlagProvider) Parse() {
	flag.Parse()
}
