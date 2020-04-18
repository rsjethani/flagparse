package flagparse

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func Test_FlagSet_addFlagFromTag_InvalidInput(t *testing.T) {
	fs := NewFlagSet()
	testValue := newIntValue(new(int))
	data := []string{
		"hello=hi",
		"nargs=99999999999999999999999",
		"name=pos-flag,nargs=0",
		"name=--opt:no-prefix",
	}
	for _, input := range data {
		if err := fs.addFlagFromTag(testValue, input, ""); err == nil {
			t.Errorf("Testing: addFlagFromTag(%q); expected: error; got: no error", input)
		}
	}
}

func Test_FlagSet_addFlagFromTag_ValidInput(t *testing.T) {
	testValue := newIntValue(new(int))
	data := []string{
		"",
		"name=pos-name,nargs=10",
		"name=--opt-name,nargs=10,usage=hello",
	}
	for _, input := range data {
		fs := NewFlagSet()
		if err := fs.addFlagFromTag(testValue, input, "field-name"); err != nil {
			t.Errorf("Testing: addFlagFromTag(%q); expected: error; got: no error", input)
		}
	}
}

func Test_NewFlagSet(t *testing.T) {
	flagSet := NewFlagSet()
	expected := &FlagSet{
		optFlags: make(map[string]*Flag),
		usageOut: os.Stderr,
		name:     os.Args[0],
		CmdArgs:  os.Args[1:],
	}
	if !reflect.DeepEqual(flagSet, expected) {
		t.Errorf("Testing: NewFlagSet(); Expected: %#v; Got: %#v", expected, flagSet)
	}
}

func Test_FlagSet_Add_Invalid(t *testing.T) {
	testVar := 100
	posFlag := NewIntFlag(&testVar, true, "")
	optFlag := NewIntFlag(&testVar, false, "")
	fs := NewFlagSet()
	fs.optFlags[defaultOptPrefix+"opt1"] = optFlag
	fs.posFlags = append(fs.posFlags, posWithName{"pos1", posFlag})
	data := []struct {
		fl       *Flag
		name     string
		optNames []string
	}{
		// Test adding new positional flag with existing name
		{fl: posFlag, name: "pos1"},
		// Test adding new optional flag with existing name
		{fl: optFlag, name: defaultOptPrefix + "opt1"},
		{fl: posFlag, name: defaultOptPrefix + "name-with-prefix"},
		{fl: optFlag, name: "name-without-prefix"},
		{fl: optFlag, name: helpLong},
		{fl: optFlag, name: helpShort},
		{optFlag, defaultOptPrefix + "name-with-prefix", []string{"opt-name-no-prefix"}},
	}
	for _, input := range data {
		if err := fs.Add(input.fl, input.name, input.optNames...); err == nil {
			t.Errorf("Testing: FlagSet.Add(%p,%q,%q); Expected: error; Got: nil", input.fl,
				input.name, input.optNames)
		}
	}
}

func Test_FlagSet_Add_Valid(t *testing.T) {
	fs := NewFlagSet()
	data := []struct {
		fl       *Flag
		name     string
		optNames []string
		validate func() error
	}{
		// Test: passing nil as flag does not affect flagset state
		{
			fl: nil,
			validate: func() error {
				if len(fs.optFlags) == 0 && len(fs.posFlags) == 0 {
					return nil
				}
				return fmt.Errorf("flag got added to flagset")
			}},
		// Test: adding positional flag with valid name
		{
			fl:   NewIntFlag(new(int), true, ""),
			name: "name-without-prefix",
			validate: func() error {
				if len(fs.posFlags) == 1 {
					return nil
				}
				return fmt.Errorf("positional flag not added")
			}},
		// Test: adding optional flag with valid name and valid extra names
		{
			fl:       NewIntFlag(new(int), false, ""),
			name:     defaultOptPrefix + "name-with-prefix",
			optNames: []string{defaultOptPrefix + "another-name"},
			validate: func() error {
				if len(fs.optFlags) == 2 {
					return nil
				}
				return fmt.Errorf("optional flag not added")
			}},
	}
	for _, input := range data {
		if err := fs.Add(input.fl, input.name, input.optNames...); err != nil {
			t.Errorf("Testing: FlagSet.Add(%p,%q,%q); Expected: no error; Got: %s", input.fl,
				input.name, input.optNames, err)
		}
		if err := input.validate(); err != nil {
			t.Errorf("Testing: FlagSet.Add(%p,%q,%q); Expected: no error; Got: %v", input.fl,
				input.name, input.optNames, err)
		}
	}
}

func Test_NewFlagSetFrom_InvalidInputs(t *testing.T) {
	data := []interface{}{
		// Test nil as input
		nil,
		// Test non-pointer as input
		*new(bool),
		// Test pointer to a non-struct as input
		new(bool),
		// Test unsupported field type as input
		&struct {
			Field1 int8 `flagparse:""`
		}{},
		// Test error returned from newFlagFromKVs()
		&struct {
			Field1 int `flagparse:"name=A_B"`
		}{},
	}
	for _, input := range data {
		if flagSet, err := NewFlagSetFrom(input); flagSet != nil || err == nil {
			t.Errorf("testing: NewFlagSet(%#v); expected: (nil, error); got: (%v, %#v)", input, flagSet, err)
		}
	}
}

func Test_NewFlagSetFrom_ValidInputs(t *testing.T) {
	args := struct {
		Field0 int // should get ignored
		field1 int `flagparse:""`              // should get ignored
		Field2 int `flagparse:""`              // expected a positional flag
		Field3 int `flagparse:"name=--field3"` // expected an optional flag
	}{}

	flagSet, err := NewFlagSetFrom(&args)
	if err != nil {
		t.Errorf("Testing: NewFlagSetFrom(%#v); Expected: no error; Got: %v", args, err)
	}
	if len(flagSet.optFlags) != 1 {
		t.Errorf("Testing: NewFlagSetFrom(%#v); Expected: one optional Flag in FlagSet; Got: %d", args, len(flagSet.optFlags))
	}

	if len(flagSet.posFlags) != 1 {
		t.Errorf("Testing: NewFlagSetFrom(%#v); Expected: one positional Flag in FlagSet; Got: %d", args, len(flagSet.posFlags))
	}
}

func Test_usage_UserDefined(t *testing.T) {
	fs := NewFlagSet()
	var called bool
	fs.Usage = func() {
		called = true
	}
	fs.usage()
	if !called {
		t.Errorf("Testing: Flagset.usage(); Expected: user defined function should be called; Got: not called")
	}
}

type testConfig struct {
	Pos1 int       `flagparse:"usage=pos1 usage"`
	Pos2 []float64 `flagparse:"usage=pos2 usage,nargs=2"`
	Opt1 string    `flagparse:"name=-t:--opt1,usage=--opt1 usage"`
	Opt2 []int     `flagparse:"name=--opt2,usage=--opt2 usage,nargs=-1"`
	Opt3 []string  `flagparse:"name=--opt3,usage=--opt3 usage,nargs=2"`
	Sw1  bool      `flagparse:"name=-s,usage=-s usage,nargs=0"`
}

func Test_Parse_InvalidInputs(t *testing.T) {
	cfg := &testConfig{}
	data := [][]string{
		// no arguments, not even for positional flags
		{},
		// argument given for pos1 but wrong type
		{"not a number"},
		// pos1 ok but not enough args for pos2: required 2; given 1
		{"10", "1.1"},
		// pos1 ok, no. of args for pos2 is ok but invalid format
		{"10", "1.1", "abc"},
		// pos1 ok, pos2 ok, but unrecognized optional flag given
		{"10", "1.1", "2.2", "--dummy", "dummy's value"},
		// pos1 ok, pos2 ok, but opt1 requires 1 argument
		{"10", "1.1", "2.2", "--opt1"},
		// pos1 ok, pos2 ok, but opt2 requires at least one argument
		{"10", "1.1", "2.2", "--opt2"},
		// pos1 ok, pos2 ok, but opt3 requires two arguments
		{"10", "1.1", "2.2", "--opt3", "one"},
		// pos1 ok, pos2 ok, -s ok, for opt2 no. of args is ok but invalid format
		{"10", "1.1", "2.2", "--opt2", "23", "24b", "-s"},
		// arguments for given flags are ok but 'extra arg' is unwanted/unrecognized argument
		{"10", "1.1", "2.2", "--opt2", "11", "22", "--opt1", "hello", "--opt3", "one", "two", "-s", "extra arg"},
	}
	for _, input := range data {
		fs, err := NewFlagSetFrom(cfg)
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}
		f, _ := os.Create(os.DevNull)
		fs.SetOutput(f)
		fs.CmdArgs = input
		fs.ContinueOnError = true
		if err := fs.Parse(); err == nil {
			t.Errorf("Testing: FlagSet.Parse(); Expected: error with %q as args; Got: no error", input)
		}
		if _, ok := err.(*ErrHelpInvoked); ok {
			t.Errorf("Testing: FlagSet.Parse(); Expected: error should not be of %[1]T type with %[2]q args; Got: error %[1]T type", &ErrHelpInvoked{}, input)
		}
	}
}

func Test_Parse_ValidInputs(t *testing.T) {
	cfg := &testConfig{
		Opt1: "hello",
		Opt2: []int{11},
		Opt3: []string{"one", "two"},
		Sw1:  false,
	}
	data := []struct {
		args     []string
		expected *testConfig
	}{
		{ // all positional flags are satisfied
			args: []string{"10", "1.2", "3.4"},
			expected: &testConfig{
				Pos1: 10,
				Pos2: []float64{1.2, 3.4},
				Opt1: "hello",
				Opt2: []int{11},
				Opt3: []string{"one", "two"},
				Sw1:  false,
			},
		},
		{ // all possible flags are satisfied
			args: []string{"10", "1.1", "2.2", "--opt2", "11", "22", "--opt1", "hello", "--opt3", "one", "two", "-s"},
			expected: &testConfig{
				Pos1: 10,
				Pos2: []float64{1.1, 2.2},
				Opt1: "hello",
				Opt2: []int{11, 22},
				Opt3: []string{"one", "two"},
				Sw1:  true,
			},
		},
	}
	for _, input := range data {
		fs, err := NewFlagSetFrom(cfg)
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}
		f, _ := os.Create(os.DevNull)
		fs.SetOutput(f)

		fs.CmdArgs = input.args
		fs.ContinueOnError = true
		if err := fs.Parse(); err != nil {
			t.Errorf("Testing: FlagSet.Parse(); Expected: no error with %q as args; Got: error %q", input.args, err)
		}
		if !reflect.DeepEqual(cfg, input.expected) {
			t.Errorf("Testing: FlagSet.Parse(); Expected: %+v; Got:%+v", input.expected, cfg)
		}
	}
}

// func Test_Parse_PosFlagWithUnlimitedArgs(t *testing.T) {
// 	type testCfg struct {
// 		Pos1 []int `flagparse:"nargs=-1"`
// 	}
// 	cfg := &testCfg{}

// 	good := []struct {
// 		args     []string
// 		expected *testCfg
// 	}{
// 		{
// 			args:     []string{"11"},
// 			expected: &testCfg{Pos1: []int{11}},
// 		},
// 		{
// 			args:     []string{"11", "22", "33", "44", "55"},
// 			expected: &testCfg{Pos1: []int{11, 22, 33, 44, 55}},
// 		},
// 	}
// 	for _, input := range good {
// 		fs, err := NewFlagSetFrom(cfg)
// 		if err != nil {
// 			t.Fatalf("Unexpected error: %q", err)
// 		}
// 		f, _ := os.Create(os.DevNull)
// 		fs.SetOutput(f)
// 		fs.ContinueOnError = true

// 		fs.CmdArgs = input.args
// 		if err := fs.Parse(); err != nil {
// 			t.Errorf("Testing: FlagSet.Parse(); Expected: no error with %q as args; Got: error %q", input.args, err)
// 		}
// 		if !reflect.DeepEqual(cfg, input.expected) {
// 			t.Errorf("Testing: FlagSet.Parse(); Expected: %+v; Got:%+v", input.expected, cfg)
// 		}
// 	}

// 	bad := [][]string{{}, {"11", "22", "33", "44abc", "55"}}
// 	for _, input := range bad {
// 		fs, err := NewFlagSetFrom(cfg)
// 		if err != nil {
// 			t.Fatalf("Unexpected error: %q", err)
// 		}
// 		f, _ := os.Create(os.DevNull)
// 		fs.SetOutput(f)
// 		fs.ContinueOnError = true

// 		fs.CmdArgs = input
// 		if err := fs.Parse(); err == nil {
// 			t.Errorf("Testing: FlagSet.Parse(); Expected: error with %q as args; Got: no error", input)
// 		}
// 	}
// }

func Test_Parse_HelpOption(t *testing.T) {
	fs, _ := NewFlagSetFrom(&testConfig{})
	fs.Desc = "flagset description"
	fs.ContinueOnError = true
	f, _ := os.Create(os.DevNull)
	fs.SetOutput(f)
	fs.CmdArgs = []string{helpLong}
	err := fs.Parse()
	if err == nil {
		t.Errorf("Testing: FlagSet.Parse(); Expected: error with %q args; Got: no error", fs.CmdArgs)
	}
	if _, ok := err.(*ErrHelpInvoked); !ok {
		t.Errorf("Testing: FlagSet.Parse(); Expected: error of type %T with %q args; Got: error of other type", &ErrHelpInvoked{}, fs.CmdArgs)
	}
}

func Test_Parse_ExitOnError(t *testing.T) {
	testParseExit := func() {
		fs, _ := NewFlagSetFrom(&testConfig{})
		f, _ := os.Create(os.DevNull)
		fs.SetOutput(f)
		fs.CmdArgs = []string{os.Getenv("CMD_ARG")}
		fs.Parse()
	}
	if _, ok := os.LookupEnv("CMD_ARG"); ok {
		testParseExit()
		return
	}
	data := []struct {
		arg      string
		expected int
	}{
		{helpLong, 1},
		{"dummy-flag", 2},
	}
	for _, input := range data {
		cmd := exec.Command(os.Args[0], "-test.run=Test_Parse_ExitOnError")
		cmd.Env = append(os.Environ(), "CMD_ARG="+input.arg)
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok {
			if e.ExitCode() != input.expected {
				t.Errorf("Testing: FlagSet.Parse(); Expected: exit code %d with args %q; Got: exit code %d", input.expected, input.arg, e.ExitCode())
			}
		}
	}
}

func Test_splitKV(t *testing.T) {
	data := make(map[string][]string)
	data[""] = []string{}
	data[fmt.Sprintf("%[1]c%[1]c%[1]c", kvPairSep)] = []string{}
	data[fmt.Sprintf("a%[1]cb%[1]cc%[1]cd", kvPairSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca%[1]cb%[1]cc%[1]cd%[1]c", kvPairSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca\\%[1]cb%[1]cc%[1]cd\\%[1]c", kvPairSep)] = []string{fmt.Sprintf("a%cb", kvPairSep), "c", fmt.Sprintf("d%c", kvPairSep)}
	for input, expected := range data {
		if got := splitKVs(input, kvPairSep); !reflect.DeepEqual(expected, got) {
			t.Errorf("Testing: splitKV(%q, %c); Expected: %q; Got: %q", input, kvPairSep, expected, got)
		}
	}
}

func Test_parseKVs_InvalidKeyValues(t *testing.T) {
	invalidKVs := []string{
		"hello",
		"usage=",
		"hello=hi",
		"name=flag_name",
		"name=flag-name:",
		"name=:flag-name",
		"nargs=1x",
	}

	for _, kv := range invalidKVs {
		if _, err := parseKVs(kv); err == nil {
			t.Errorf("Testing: parseTags(%q); Expected: error; Got: no error", kv)
		}
	}
}

func Test_parseKVs_ValidKeyValues(t *testing.T) {
	data := []struct {
		kvs      string
		expected map[string]string
	}{
		{
			"",
			map[string]string{},
		},
		{
			"name=flag-name123,usage=hello\\,world,nargs=10",
			map[string]string{
				nargsKey: "10",
				nameKey:  "flag-name123",
				usageKey: "hello,world",
			},
		},
		{
			"name=-f123:--Flag-Name123,usage=abc,nargs=-10",
			map[string]string{
				nargsKey: "-10",
				nameKey:  "-f123:--Flag-Name123",
				usageKey: "abc",
			},
		},
	}

	for _, input := range data {
		got, err := parseKVs(input.kvs)
		if err != nil {
			t.Errorf("Testing: parseTags(%q); Expected: no error; Got: %s", input.kvs, err)
		}

		if !reflect.DeepEqual(input.expected, got) {
			t.Errorf("Testing: parseTags(%q); Expected: %+v; Got: %+v", input.kvs, input.expected, got)
		}
	}
}
