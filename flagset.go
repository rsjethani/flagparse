package flagparse

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	stateInit int = iota
	statePosFlag
	stateOptFlag
	kvSep            rune   = '='
	kvPairSep        rune   = ','
	optNameSep       string = ":"
	defaultOptPrefix string = "-"
	nameKey          string = "name"
	usageKey         string = "usage"
	nargsKey         string = "nargs"
	helpShort        string = "-h"
	helpLong         string = "--help"
	packageTag       string = "flagparse"
)

var validKVs = map[string]*regexp.Regexp{
	usageKey: regexp.MustCompile(fmt.Sprintf(`^%s%c(.+)$`, usageKey, kvSep)),
	nargsKey: regexp.MustCompile(fmt.Sprintf(`^%s%c(-?[[:digit:]]+)$`, nargsKey, kvSep)),
	nameKey: regexp.MustCompile(fmt.Sprintf(`^%s%c([-[:alnum:]]+(%s[-[:alnum:]]+)*)$`, nameKey,
		kvSep, optNameSep)),
}

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
	Usage           func()
	usageOut        io.Writer
	CmdArgs         []string
	posFlags        []posWithName
	optFlags        map[string]*Flag
	// following fields are used for state changes during parsing
	iPos        int
	curFlag     *Flag
	curFlagName string
	curFlagArgs []string
}

func NewFlagSet() *FlagSet {
	fs := &FlagSet{
		optFlags: make(map[string]*Flag),
		usageOut: os.Stderr,
		name:     os.Args[0],
		CmdArgs:  os.Args[1:],
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
		return nil, fmt.Errorf("src must be a pointer to struct")
	}

	// get Type data of the actual struct pointed by the pointer,
	// verify that it is a struct
	srcTyp = srcTyp.Elem()
	if srcTyp.Kind() != reflect.Struct {
		return nil, fmt.Errorf("src must be a pointer to struct")
	}

	srcVal := reflect.ValueOf(src).Elem()

	fs := NewFlagSet()
	// iterate over all fields of the struct, parse the value of 'packageTag'
	// and create flags accordingly. Skip any field not having the tag.
	for i := 0; i < srcTyp.NumField(); i++ {
		fieldType := srcTyp.Field(i)
		fieldVal := srcVal.Field(i)
		tagValue, tagged := srcTyp.Field(i).Tag.Lookup(packageTag)
		// ignore fields which are untagged and/or unexported
		if !tagged || !fieldVal.Addr().CanInterface() {
			continue
		}

		val, err := newValue(fieldVal.Addr().Interface())
		if err != nil {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, err)
		}

		err = fs.addFlagFromTag(val, tagValue, fieldType.Name)
		if err != nil {
			return nil, fmt.Errorf("Error while creating flag from field '%s': %s", fieldType.Name, err)
		}
	}

	return fs, nil
}

func validPosName(name string) bool {
	return regexp.MustCompile(`^[[:alnum:]][-[:alnum:]]+$`).MatchString(name)
}

func validOptName(name string) bool {
	return regexp.MustCompile(`^-[-[:alnum:]]+$`).MatchString(name)
}

func (fs *FlagSet) Add(fl *Flag, name string, optNames ...string) error {
	if fl == nil {
		return nil
	}
	if fl.positional {
		if !validPosName(name) {
			return fmt.Errorf("%q is not a valid positional flag name", name)
		}
		// check for duplicate name
		for _, v := range fs.posFlags {
			if name == v.name {
				return fmt.Errorf("positional flag with name %q already exists", name)
			}
		}
		fs.posFlags = append(fs.posFlags, posWithName{name, fl})
	} else {
		names := []string{name}
		names = append(names, optNames...)
		for _, nm := range names {
			if nm == helpShort || nm == helpLong {
				return fmt.Errorf("flag names %s,%s are reserved", helpShort, helpLong)
			}
			if !validOptName(nm) {
				return fmt.Errorf("%q is not a valid optional flag name", nm)
			}
			// check for duplicate name
			for v := range fs.optFlags {
				if nm == v {
					return fmt.Errorf("optional flag with name %q already exists", nm)
				}
			}
			fs.optFlags[nm] = fl
		}
	}
	return nil
}

func splitKVs(src string, sep rune) []string {
	backSlash := '\\'
	parts := make([]string, 0)
	b := &strings.Builder{}
	var prevRune rune
	for _, curRune := range src {
		switch {
		// rune is a backSlash simply skip it
		case curRune == backSlash:
		// rune is a sep but it is not escaped by backSlash
		case curRune == sep && prevRune != backSlash:
			if b.Len() != 0 {
				parts = append(parts, b.String())
				b.Reset()
			}
		// rune is either not a sep/backslash or if it is a sep then it is escaped by backskash
		default:
			b.WriteRune(curRune)
		}
		prevRune = curRune
	}
	// append any remaining runes between last sep and end of src
	if b.Len() != 0 {
		parts = append(parts, b.String())
	}
	return parts
}

func parseKVs(tagValue string) (map[string]string, error) {
	kvs := make(map[string]string)
	for _, key := range splitKVs(tagValue, kvPairSep) {
		invalid := true
		for name, regex := range validKVs {
			res := regex.FindStringSubmatch(key)
			if len(res) >= 2 {
				kvs[name] = res[1]
				invalid = false
			}
		}
		if invalid {
			return nil, fmt.Errorf("unknown key and/or invalid value: %s", key)
		}
	}
	return kvs, nil
}

func (fs *FlagSet) addFlagFromTag(value Value, tagValue string, fieldName string) error {
	keyValues, err := parseKVs(tagValue)
	if err != nil {
		return err
	}

	var fl *Flag

	// create flag
	if strings.HasPrefix(keyValues[nameKey], defaultOptPrefix) {
		fl = NewFlag(value, false, keyValues[usageKey])
	} else {
		fl = NewFlag(value, true, keyValues[usageKey])
	}

	// set nargs for the flag
	if keyValues[nargsKey] != "" {
		nargs, err := strconv.ParseInt(keyValues[nargsKey], 0, strconv.IntSize)
		if err != nil {
			return formatParseError(keyValues[nargsKey], fmt.Sprintf("%T", int(1)), err)
		}
		err = fl.SetNArgs(int(nargs))
		if err != nil {
			return err
		}
	}

	names := strings.Split(keyValues[nameKey], optNameSep)
	// if no name is given then use field's name in lower case
	if names[0] == "" {
		names[0] = strings.ToLower(fieldName)
	}
	return fs.Add(fl, names[0], names[1:]...)
}

func (fs *FlagSet) writeAndCloseFlag() error {
	if err := fs.curFlag.value.Set(fs.curFlagArgs...); err != nil {
		return err
	}
	fs.curFlagName = ""
	fs.curFlag = nil
	fs.curFlagArgs = fs.curFlagArgs[:0]
	return nil
}

func (fs *FlagSet) processArg(arg string) error {
	if arg != fs.curFlagName {
		fs.curFlagArgs = append(fs.curFlagArgs, arg)
	}
	if len(fs.curFlagArgs) == fs.curFlag.nArgs {
		return fs.writeAndCloseFlag()
	}
	return nil
}

func (fs *FlagSet) closeFlag() error {
	if fs.curFlag.nArgs < 0 {
		if len(fs.curFlagArgs) < 1 {
			return fmt.Errorf("flag %s expects at least one argument", fs.curFlagName)
		}
		if err := fs.writeAndCloseFlag(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("flag %s expects %d arguments, given %d", fs.curFlagName,
			fs.curFlag.nArgs, len(fs.curFlagArgs))
	}
	return nil
}

func (fs *FlagSet) openFlag(curArg string, pos bool) error {
	if pos {
		fs.curFlagName = fs.posFlags[fs.iPos].name
		fs.curFlag = fs.posFlags[fs.iPos].flag
	} else {
		fs.curFlagName = curArg
		fs.curFlag = fs.optFlags[fs.curFlagName]
	}
	return fs.processArg(curArg)
}

func (fs *FlagSet) parse() error {
	for _, curArg := range fs.CmdArgs {
		// does it looks like an optional flag?
		if strings.HasPrefix(curArg, defaultOptPrefix) {
			// is it a help flag?
			if curArg == helpShort || curArg == helpLong {
				return &ErrHelpInvoked{}
			}
			// if this is not a known flag then return error
			if _, ok := fs.optFlags[curArg]; !ok {
				return fmt.Errorf("unrecognized flag %s", curArg)
			}
			// since this is a know flag so before starting  to process it, try closing the current
			// pos/opt flag if any
			if fs.curFlagName != "" {
				if err := fs.closeFlag(); err != nil {
					return err
				}
			}
			// since this is a known opt flag, no flag is opened currently, open this flag for
			// processing
			if err := fs.openFlag(curArg, false); err != nil {
				return err
			}
			continue
		}
		// if there is a pos/opt flag open then process current argument for it
		if fs.curFlagName != "" {
			if err := fs.processArg(curArg); err != nil {
				return err
			}
			continue
		}
		// since there is no pos/opt flag open see if there are any pos flags that can be opened
		// for processing this argument
		if fs.iPos < len(fs.posFlags) {
			if err := fs.openFlag(curArg, true); err != nil {
				return err
			}
			fs.iPos++
			continue
		}
		// since all pos/opt flags have been processed, current argument is unwanted/unrecognized
		return fmt.Errorf("unrecognized argument: %s", curArg)
	}
	if fs.curFlagName != "" {
		if err := fs.closeFlag(); err != nil {
			return err
		}
	}
	if fs.iPos < len(fs.posFlags) {
		var names []string
		for _, fl := range fs.posFlags[fs.iPos:] {
			names = append(names, fl.name)
		}
		return fmt.Errorf("arguments are required for these flags: %v", strings.Join(names, ", "))
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

func (fs *FlagSet) SetOutput(w io.Writer) {
	if w != nil {
		fs.usageOut = w
	}
}

type optWithName struct {
	name string
	fl   *Flag
}

func (fs *FlagSet) optMapToList() []optWithName {
	var optList []optWithName
	for nm, f := range fs.optFlags {
		for i := range optList {
			if optList[i].fl == f {
				optList[i].name = optList[i].name + ", " + nm
				goto end
			}
		}
		optList = append(optList, optWithName{nm, f})
	end:
	}
	sort.SliceStable(optList, func(i, j int) bool { return optList[i].name < optList[j].name })
	return optList
}

func (fs *FlagSet) defaultUsage() {
	out := fs.usageOut
	fmt.Fprintf(out, "\nUsage of %s:\n", fs.name)
	if fs.Desc != "" {
		fmt.Fprintf(out, "\n%s\n", fs.Desc)
	}
	fmt.Fprint(out, "\nPositional Flags:")
	for _, fl := range fs.posFlags {
		fmt.Fprintf(out, "\n\n  %s  %T\n\t%s", fl.name, fl.flag.value.Get(), fl.flag.usage)
	}

	fmt.Fprint(out, "\n\nOptional Flags:")
	fmt.Fprintf(out, "\n  %s, %s\n\t%s", helpShort, helpLong, "Show this usage message and exit")
	for _, v := range fs.optMapToList() {
		if v.fl.isSwitch() {
			fmt.Fprintf(out, "\n\n  %s\n\t%s", v.name, v.fl.usage)
			continue
		}
		nargs := "Requires: 1 or more arguments"
		if v.fl.nArgs > 0 {
			nargs = fmt.Sprintf("Requires: %v argument(s)", v.fl.nArgs)
		}
		def := `Default: ""`
		if v.fl.defVal != "" {
			def = fmt.Sprintf("Default: %s", v.fl.defVal)
		}
		fmt.Fprintf(out, "\n\n  %s  %T\n\t%s. %s. %s", v.name, v.fl.value.Get(), v.fl.usage, nargs, def)
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
