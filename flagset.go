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
	defaultOptPrefix string = "--"
	helpFlag         string = "--help"
	packageTag       string = "flagparse"
)

type ErrHelpInvoked struct{}

func (e *ErrHelpInvoked) Error() string { return "" }

type posWithName struct {
	name string
	flag *Flag
}

type FlagSet struct {
	ContinueOnError bool
	name            string
	Desc            string
	OptPrefix       string
	Usage           func()
	usageOut        io.Writer
	CmdArgs         []string
	posFlags        []posWithName
	optFlags        map[string]*Flag
}

func (fs *FlagSet) SetOutput(w io.Writer) {
	if w != nil {
		fs.usageOut = w
	}
}

func NewFlagSet() *FlagSet {
	fs := &FlagSet{
		OptPrefix: defaultOptPrefix,
		optFlags:  make(map[string]*Flag),
		usageOut:  os.Stderr,
		name:      os.Args[0],
		CmdArgs:   os.Args[1:],
	}
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
		// ignore fields which are untagged and/or unexported
		if !tagged || !fieldVal.Addr().CanInterface() {
			continue
		}

		keyValues, err := parseKVs(structTags)
		if err != nil {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, err)
		}

		// if 'name' key not specified then simply use field's name in lower case
		if keyValues[nameKey] == "" {
			keyValues[nameKey] = strings.ToLower(fieldType.Name)
		}

		val, err := newValue(fieldVal.Addr().Interface())
		if err != nil {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, err)
		}

		flag, err := newFlagFromKVs(val, keyValues)
		if err != nil {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, err)
		}

		fs.Add(keyValues["name"], flag)
	}

	return fs, nil
}

func (fs *FlagSet) Add(name string, fl *Flag) {
	if fl == nil {
		return
	}
	switch fl.positional {
	case true:
		fs.posFlags = append(fs.posFlags, posWithName{name, fl})
	case false:
		fs.optFlags[fs.OptPrefix+name] = fl
	}
}

func (fs *FlagSet) defaultUsage() {
	out := fs.usageOut
	fmt.Fprintf(out, "\nUsage of %s:\n", fs.name)
	if fs.Desc != "" {
		fmt.Fprintf(out, "\n%s\n", fs.Desc)
	}
	fmt.Fprint(out, "\nPositional Arguments:")
	for _, fl := range fs.posFlags {
		fmt.Fprintf(out, "\n  %s  %T\n\t%s", fl.name, fl.flag.value.Get(), fl.flag.usage)
	}

	fmt.Fprint(out, "\n\nOptional Arguments:")
	fmt.Fprintf(out, "\n  %s\n\t%s", helpFlag, "Show this usage message and exit")
	for name, fl := range fs.optFlags {
		if fl.isSwitch() {
			fmt.Fprintf(out, "\n  %s\n\t%s", name, fl.usage)
			continue
		}
		fmt.Fprintf(out, "\n  %s  %T\n\t%s  (Default: %s)", name, fl.value.Get(), fl.usage, fl.defVal)
	}

	fmt.Fprint(out, "\n")
}

// usage calls the Usage method if one is specified,
// or the appropriate default usage function otherwise.
func (fs *FlagSet) usage() {
	if fs.Usage == nil {
		fs.defaultUsage()
	} else {
		fs.Usage()
	}
}

func (fs *FlagSet) parse() error {
	cmdArgs := fs.CmdArgs
	visited := make(map[string]bool)

	for iArgs, iPos, curState := 0, 0, stateInit; iArgs < len(cmdArgs); {
		curArg := cmdArgs[iArgs]
		switch curState {
		case stateInit:
			if curArg == helpFlag {
				return &ErrHelpInvoked{}
			}

			// if curArg starts with the configured prefix then process it as an optional arg
			if strings.HasPrefix(curArg, fs.OptPrefix) {
				if _, found := fs.optFlags[curArg]; found {
					if visited[curArg] { // if defined but already processed then return error
						return fmt.Errorf("flag '%s' already given", curArg)
					}
					curState = stateOptFlag
					break
				} else { // if not defined as an opt arg then return error
					return fmt.Errorf("unknown optional flag: %s", curArg)
				}
			}

			// if all positional flags have not been processed yet then consider
			// curArg as the value for next positional arg
			if iPos < len(fs.posFlags) {
				curState = statePosFlag
				break
			}

			// since all defined positional and optional args have been processed, curArg
			// is an undefined positional flag
			return fmt.Errorf("Unknown positional flag: %s", curArg)
		case statePosFlag:
			name := fs.posFlags[iPos].name
			val := fs.posFlags[iPos].flag.value
			nargs := fs.posFlags[iPos].flag.nArgs
			given := len(cmdArgs) - iArgs
			if nargs < 0 {
				if err := val.Set(cmdArgs[iArgs:]...); err != nil {
					return fmt.Errorf("error while setting flag '%s': %s", name, err)
				}
				iArgs = len(cmdArgs)
			} else {
				if given < nargs {
					return fmt.Errorf("invalid no. of arguments for flag '%s'; required: %d, given: %d", name, nargs, given)
				}
				if err := val.Set(cmdArgs[iArgs : iArgs+nargs]...); err != nil {
					return fmt.Errorf("error while setting flag '%s': %s", name, err)
				}
				iArgs += nargs
			}
			iPos++
			visited[name] = true
			curState = stateInit
		case stateOptFlag:
			name := curArg
			val := fs.optFlags[name].value
			nargs := fs.optFlags[name].nArgs
			given := len(cmdArgs) - 1 - iArgs
			if nargs < 0 { // unlimited no. of arguments
				if given < 1 {
					return fmt.Errorf("invalid no. of arguments for flag '%s'; required: at least one, given: 0", name)
				}
				if err := val.Set(cmdArgs[iArgs+1:]...); err != nil {
					return fmt.Errorf("error while setting flag '%s': %s", name, err)
				}
				iArgs = len(cmdArgs)
			} else if nargs > 0 { // limited no. of arguments
				if given < nargs {
					return fmt.Errorf("invalid no. of arguments for flag '%s'; required: %d, given: %d", name, nargs, given)
				}
				if err := val.Set(cmdArgs[iArgs+1 : iArgs+1+nargs]...); err != nil {
					return fmt.Errorf("error while setting flag '%s': %s", name, err)
				}
				iArgs += nargs + 1
			} else { // zero arguments i.e. a switch
				val.Set()
				iArgs++
			}
			visited[name] = true
			curState = stateInit
		}
	}
	for _, pos := range fs.posFlags {
		if !visited[pos.name] {
			return fmt.Errorf("Error: value for positional flag '%s' not given", pos.name)
		}
	}
	return nil
}

func (fs *FlagSet) Parse() error {
	err := fs.parse()
	if err == nil {
		return nil
	}

	var exitCode int
	switch err.(type) {
	case *ErrHelpInvoked:
		exitCode = 1
	default:
		exitCode = 2
		fmt.Fprintln(fs.usageOut, err)
	}
	fs.usage()
	if !fs.ContinueOnError {
		os.Exit(exitCode)
	}
	return err
}
