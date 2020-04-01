package flagparse

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

const (
	stateInit int = iota
	statePosFlag
	stateOptFlag
	stateNoArgsLeft
	defaultOptFlagPrefix string = "--"
	packageTag           string = "flagparse"
)

type posFlagWithName struct {
	name string
	flag *Flag
}

type FlagSet struct {
	continueOnError bool
	name            string
	Desc            string
	CmdArgs         []string
	Usage           func()
	usageOut        io.Writer
	OptFlagPrefix   string
	optFlags        map[string]*Flag
	posFlags        []posFlagWithName
}

func (fs *FlagSet) SetOutput(w io.Writer) {
	if w != nil {
		fs.usageOut = w
	}
}

func (fs *FlagSet) addHelp() {
	var help bool
	fs.Add("help", NewSwitchFlag(NewBool(&help), "Show this help message and exit"))
}

func NewFlagSet() *FlagSet {
	fs := &FlagSet{
		OptFlagPrefix: defaultOptFlagPrefix,
		optFlags:      make(map[string]*Flag),
		usageOut:      os.Stderr,
		name:          os.Args[0],
		CmdArgs:       os.Args[1:],
	}
	fs.addHelp()
	return fs
}

func NewFlagSetFrom(src interface{}) (*FlagSet, error) {
	if src == nil {
		return nil, fmt.Errorf("src cannot be nil")
	}
	// get Type data of src, verify that it is of pointer type
	srcTyp := reflect.TypeOf(src)
	if srcTyp.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("src must be a pointer")
	}

	// get Type data of the actual struct pointed by the pointer,
	// verify that it is a struct
	srcTyp = srcTyp.Elem()
	if srcTyp.Kind() != reflect.Struct {
		return nil, fmt.Errorf("src must be a pointer to a struct")
	}

	srcVal := reflect.ValueOf(src).Elem()

	fs := NewFlagSet()
	// iterate over all fields of the struct, parse the value of 'packageTag'
	// and create flags accordingly. Skip any field not having the tag.
	for i := 0; i < srcTyp.NumField(); i++ {
		fieldType := srcTyp.Field(i)
		fieldVal := srcVal.Field(i)
		structTags, tagged := srcTyp.Field(i).Tag.Lookup(packageTag)
		if !tagged {
			continue
		}

		if !fieldVal.Addr().CanInterface() {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, "unexported struct field")
		}
		val, err := NewValue(fieldVal.Addr().Interface())
		if err != nil {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, err)
		}

		flag, name, err := newFlagFromTags(val, fieldType.Name, structTags)
		if err != nil {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, err)
		}

		fs.Add(name, flag)
	}

	return fs, nil
}

func (fs *FlagSet) Add(name string, fl *Flag) {
	if fl == nil {
		return
	}
	if fl.positional {
		fs.posFlags = append(fs.posFlags, posFlagWithName{name: name, flag: fl})
		return
	}
	fs.optFlags[fs.OptFlagPrefix+name] = fl
}

func (fs *FlagSet) defaultUsage() {
	out := fs.usageOut
	fmt.Fprintf(out, "\nUsage of %s:\n", fs.name)
	if fs.Desc != "" {
		fmt.Fprintf(out, "\n%s\n", fs.Desc)
	}
	fmt.Fprint(out, "\nPositional Arguments:")
	for _, p := range fs.posFlags {
		val := p.flag.value.Get()
		fmt.Fprintf(out, "\n  %[1]s  %[2]T\n\t%[3]s  (Default: %[2]v)", p.name, val, p.flag.help, val)
	}

	// TODO: show list of opt args in sorted order
	fmt.Fprint(out, "\n\nOptional Arguments:")
	for name, arg := range fs.optFlags {
		if arg.isSwitch() {
			fmt.Fprintf(out, "\n  %[1]s\n\t%s", name, arg.help)
			continue
		}
		val := arg.value.Get()
		fmt.Fprintf(out, "\n  %[1]s  %[2]T\n\t%[3]s  (Default: %[2]v)", name, val, arg.help)
	}

	fmt.Fprint(out, "\n")
}

// usage calls the Usage method for the ArgSet if one is specified,
// or the appropriate default usage function otherwise.
func (fs *FlagSet) usage() {
	if fs.Usage == nil {
		fs.defaultUsage()
	} else {
		fs.Usage()
	}
}

func (fs *FlagSet) parse() error {
	argsToParse := fs.CmdArgs
	curState := stateInit
	var curArg string
	visited := make(map[string]bool)
	var posIndex, argsIndex int

	getArg := func(i int) string {
		if i < len(argsToParse) {
			return argsToParse[i]
		}
		return ""
	}

	for {
		switch curState {
		case stateInit:
			arg := getArg(argsIndex)
			if arg == "" {
				curState = stateNoArgsLeft
				break
			}
			curArg = arg

			// if curArg starts with the configured prefix then process it as an optional arg
			if strings.HasPrefix(curArg, fs.OptFlagPrefix) {
				if _, found := fs.optFlags[curArg]; found {
					if visited[curArg] { // if curArg is defined but already processed then return error
						return fmt.Errorf("option '%s' already given", curArg)
					}
					curState = stateOptFlag
					break
				} else { // if curArg is not defined as an opt arg then return error
					return fmt.Errorf("unknown optional flag: %s", curArg)
				}
			}

			// if all positional args have not been processed yet then consider
			// curArg as the value for next positional arg
			if posIndex < len(fs.posFlags) {
				curState = statePosFlag
				break
			}

			// since all defined positional and optional args have been processed, curArg
			// is an undefined positional arg
			return fmt.Errorf("Unknown positional flag: %s", curArg)
		case statePosFlag:
			if err := fs.posFlags[posIndex].flag.value.Set(curArg); err != nil {
				return fmt.Errorf("error while setting option '%s': %s", fs.posFlags[posIndex].name, err)
			}
			visited[fs.posFlags[posIndex].name] = true
			posIndex++
			argsIndex++
			curState = stateInit
		case stateOptFlag:
			if fs.optFlags[curArg].nArgs == 0 {
				if curArg == "--help" {
					fs.usage()
					return nil
				}
				fs.optFlags[curArg].value.Set()
				argsIndex++
			} else if fs.optFlags[curArg].nArgs < 0 {
				if err := fs.optFlags[curArg].value.Set(argsToParse[argsIndex+1:]...); err != nil {
					return fmt.Errorf("error while setting option '%s': %s", curArg, err)
				}
				argsIndex = len(argsToParse)

			} else {
				inp := []string{}
				for i := 1; i <= fs.optFlags[curArg].nArgs; i++ {
					v := getArg(i + argsIndex)
					if v == "" {
						return fmt.Errorf("invalid no. of arguments for option '%s'; required: %d, given: %d", curArg, fs.optFlags[curArg].nArgs, i-1)
					}
					inp = append(inp, v)
				}
				if err := fs.optFlags[curArg].value.Set(inp...); err != nil {
					return fmt.Errorf("error while setting option '%s': %s", curArg, err)
				}
				argsIndex += fs.optFlags[curArg].nArgs + 1
			}
			curState = stateInit
		case stateNoArgsLeft:
			for _, pos := range fs.posFlags {
				if !visited[pos.name] {
					return fmt.Errorf("Error: value for positional flag '%s' not given", pos.name)
				}
			}
			return nil
		}
	}
}

func (fs *FlagSet) Parse() error {
	err := fs.parse()
	if err != nil && !fs.continueOnError {
		fmt.Fprintln(fs.usageOut, err)
		fs.usage()
		os.Exit(1)
	}
	return err
}
