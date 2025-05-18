package config

import (
	"flag"
	"fmt"
	"os"
)

type FlagSet struct {
	flags  *flag.FlagSet
	values map[string]any
}

func NewFlagSet() *FlagSet {
	return &FlagSet{
		flags:  flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		values: make(map[string]any),
	}
}

func (f *FlagSet) DefineString(name, value, usage string) {
	var val string
	f.flags.StringVar(&val, name, value, usage)
	f.values[name] = &val
}

func (f *FlagSet) DefineInt(name string, value int, usage string) {
	var val int
	f.flags.IntVar(&val, name, value, usage)
	f.values[name] = &val
}

func (f *FlagSet) DefineBool(name string, value bool, usage string) {
	var val bool
	f.flags.BoolVar(&val, name, value, usage)
	f.values[name] = &val
}

func (f *FlagSet) Collect() map[string]any {
	f.flags.Parse(os.Args[1:])

	result := make(map[string]any)

	for name, ptr := range f.values {
		switch v := ptr.(type) {
		case *string:
			if *v != "" {
				result[name] = *v
			}
		case *int:
			if *v != 0 {
				result[name] = *v
			}
		case *bool:
			if *v {
				result[name] = *v
			}
		}
	}

	return result
}

func (f *FlagSet) String() string {
	f.flags.Parse(os.Args[1:]) // Ensure flags are parsed
	output := ""

	for name, ptr := range f.values {
		switch v := ptr.(type) {
		case *string:
			output += name + "=" + *v + " "
		case *int:
			output += name + "=" + fmt.Sprintf("%d", *v) + " "
		case *bool:
			output += name + "=" + fmt.Sprintf("%t", *v) + " "
		}
	}

	return output
}
