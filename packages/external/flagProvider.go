package external

import (
	"flag"
)

type IFlagProvider interface {
	StringVar(p *string, name string, value string, usage string)
	NewFlagSet(name string) IFlagSet
	Parse()
}

type FlagProvider struct{}

func (f *FlagProvider) StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, value, usage)
}

func (f *FlagProvider) Parse() {
	flag.Parse()
}

type IFlagSet interface {
	StringVar(p *string, name string, value string, usage string)
	Parse(arguments []string) error
	Arg(i int) string
}

func (f *FlagProvider) NewFlagSet(name string) IFlagSet {
	return flag.NewFlagSet(name, flag.ExitOnError)
}
