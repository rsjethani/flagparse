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
	helpOptFlag      string = "--help"
	packageTag       string = "flagparse"
)

type ErrHelpInvoked struct{}

func (e *ErrHelpInvoked) Error() string { return "" }

type FlagSet struct {
	ContinueOnError bool
	name            string
	Desc            string
	OptPrefix       string
	Usage           func()
	usageOut        io.Writer
	CmdArgs         []string
	posFlags        []*Flag
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

		flag, err := newFlagFromTags(val, fieldType.Name, structTags)
		if err != nil {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, err)
		}

		fs.Add(flag)
	}

	return fs, nil
}

func (fs *FlagSet) Add(fl *Flag) {
	if fl == nil {
		return
	}
	if fl.positional {
		fs.posFlags = append(fs.posFlags, fl)
		return
	}
	fs.optFlags[fl.name] = fl
}

func (fs *FlagSet) defaultUsage() {
	out := fs.usageOut
	fmt.Fprintf(out, "\nUsage of %s:\n", fs.name)
	if fs.Desc != "" {
		fmt.Fprintf(out, "\n%s\n", fs.Desc)
	}
	fmt.Fprint(out, "\nPositional Arguments:")
	for _, p := range fs.posFlags {
		val := p.value.Get()
		fmt.Fprintf(out, "\n  %s  %T\n\t%s", p.name, val, p.help)
	}

	// TODO: show list of opt args in sorted order
	fmt.Fprint(out, "\n\nOptional Arguments:")
	fmt.Fprintf(out, "\n  %s\n\t%s", helpOptFlag, "Show this help message and exit")
	for _, arg := range fs.optFlags {
		if arg.isSwitch() {
			fmt.Fprintf(out, "\n  %s%s\n\t%s", fs.OptPrefix, arg.name, arg.help)
			continue
		}
		fmt.Fprintf(out, "\n  %s%s  %T\n\t%s  (Default: %s)", fs.OptPrefix, arg.name, arg.value.Get(), arg.help, arg.defVal)
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
	cmdArgs := fs.CmdArgs
	visited := make(map[string]bool)

	for iArgs, iPos, curState := 0, 0, stateInit; iArgs < len(cmdArgs); {
		curArg := cmdArgs[iArgs]
		switch curState {
		case stateInit:
			if curArg == helpOptFlag {
				return &ErrHelpInvoked{}
			}

			// if curArg starts with the configured prefix then process it as an optional arg
			if strings.HasPrefix(curArg, fs.OptPrefix) {
				if _, found := fs.optFlags[curArg[len(fs.OptPrefix):]]; found {
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
			if iPos < len(fs.posFlags) {
				curState = statePosFlag
				break
			}

			// since all defined positional and optional args have been processed, curArg
			// is an undefined positional arg
			return fmt.Errorf("Unknown positional flag: %s", curArg)
		case statePosFlag:
			if err := fs.posFlags[iPos].value.Set(curArg); err != nil {
				return fmt.Errorf("error while setting option '%s': %s", fs.posFlags[iPos].name, err)
			}
			visited[fs.posFlags[iPos].name] = true
			iPos++
			iArgs++
			curState = stateInit
		case stateOptFlag:
			flagName := curArg[len(fs.OptPrefix):]
			nargs := fs.optFlags[flagName].nArgs
			if nargs < 0 { // unlimited no. of arguments
				given := len(cmdArgs) - 1 - iArgs
				if given < 1 {
					return fmt.Errorf("invalid no. of arguments for option '%s'; required: at least one, given: 0", curArg)
				}
				if err := fs.optFlags[flagName].value.Set(cmdArgs[iArgs+1:]...); err != nil {
					return fmt.Errorf("error while setting option '%s': %s", curArg, err)
				}
				iArgs = len(cmdArgs)
			} else if nargs > 0 { // limited no. of arguments
				given := len(cmdArgs) - 1 - iArgs
				if given < nargs {
					return fmt.Errorf("invalid no. of arguments for option '%s'; required: %d, given: %d", curArg, nargs, given)
				}
				if err := fs.optFlags[flagName].value.Set(cmdArgs[iArgs+1 : iArgs+1+nargs]...); err != nil {
					return fmt.Errorf("error while setting option '%s': %s", curArg, err)
				}
				iArgs += nargs + 1
			} else { // zero arguments i.e. a switch
				fs.optFlags[flagName].value.Set()
				iArgs++
			}
			visited[curArg] = true
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
