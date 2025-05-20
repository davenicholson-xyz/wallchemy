package config

import (
	"flag"
	"fmt"
	"os"
)

type FlagSet struct {
	flags      *flag.FlagSet
	values     map[string]any
	hiddenFlag map[string]bool
}

func NewFlagSet() *FlagSet {
	fs := &FlagSet{
		flags:      flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		values:     make(map[string]any),
		hiddenFlag: make(map[string]bool),
	}

	fs.flags.Usage = func() {
		maxLength := 0
		fs.flags.VisitAll(func(f *flag.Flag) {
			if !fs.hiddenFlag[f.Name] && len(f.Name) > maxLength {
				maxLength = len(f.Name)
			}
		})

		maxLength += 2

		fmt.Fprintf(fs.flags.Output(), "Usage of %s:\n", "wallchemy")
		fs.flags.VisitAll(func(f *flag.Flag) {
			if !fs.hiddenFlag[f.Name] {
				// Format with padding
				fmt.Fprintf(fs.flags.Output(), "  -%-*s%s\n", maxLength, f.Name, f.Usage)
			}
		})
	}

	return fs
}

func (f *FlagSet) DefineStringHidden(name, value string) {
	f.hiddenFlag[name] = true
	f.DefineString(name, value, "")
}

func (f *FlagSet) DefineIntHidden(name string, value int) {
	f.hiddenFlag[name] = true
	f.DefineInt(name, value, "")
}

func (f *FlagSet) DefineBoolHidden(name string, value bool) {
	f.hiddenFlag[name] = true
	f.DefineBool(name, value, "")
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
	f.flags.Parse(os.Args[1:])
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
