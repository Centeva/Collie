package external

import (
	"flag"
	"log"
)

type IFlagProvider interface {
	String(name string, value string, usage string) *string
	StringVar(p *string, name string, value string, usage string)
	NewFlagSet(name string, usage string) IFlagSet
	Parse()
	PrintDefaults()
}

type FlagProvider struct{}

func (f *FlagProvider) String(name string, value string, usage string) *string {
	return flag.String(name, value, usage)
}

func (f *FlagProvider) StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, value, usage)
}

func (f *FlagProvider) Parse() {
	flag.Parse()
}

func (f *FlagProvider) PrintDefaults() {
	flag.PrintDefaults()
}

type IFlagSet interface {
	String(name string, value string, usage string) *string
	StringVar(p *string, name string, value string, usage string)
	Parse(arguments []string) error
	Arg(i int) string
	PrintDefaults()
}

func (f *FlagProvider) NewFlagSet(name string, usage string) IFlagSet {
	cmd := flag.NewFlagSet(name, flag.ExitOnError)
	cmd.Usage = func() {
		log.Printf("Usage of %s:\n", name)
		log.Printf("\t%s", usage)
		cmd.PrintDefaults()
		flag.PrintDefaults()
	}
	return cmd
}
