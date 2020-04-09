package flagparse

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func Test_NewFlagSet(t *testing.T) {
	flagSet := NewFlagSet()
	expected := &FlagSet{
		OptPrefix: defaultOptPrefix,
		optFlags:  make(map[string]*Flag),
		usageOut:  os.Stderr,
		name:      os.Args[0],
		CmdArgs:   os.Args[1:],
	}
	if !reflect.DeepEqual(flagSet, expected) {
		t.Errorf("Testing: NewFlagSet(); Expected: %#v; Got: %#v", expected, flagSet)
	}
}

func Test_Add_FlagToFlagSet(t *testing.T) {
	testValue := NewInt(new(int))
	posFlag := NewPosFlag("pos1", testValue, "help")
	optFlag := NewOptFlag("opt1", testValue, "help")
	swFlag := NewSwitchFlag("sw1", testValue, "help")

	expected := NewFlagSet()

	fs := NewFlagSet()
	fs.Add(nil)
	if !reflect.DeepEqual(fs, expected) {
		t.Errorf("Testing: FlagSet.Add(nil); Expected: %+v; Got: %+v", expected, fs)
	}

	expected.posFlags = append(expected.posFlags, posFlag)
	fs.Add(posFlag)
	if !reflect.DeepEqual(fs, expected) {
		t.Errorf("Testing: FlagSet.Add(%+v); Expected: %+v; Got: %+v", posFlag, expected, fs)
	}

	expected.optFlags[optFlag.name] = optFlag
	fs.Add(optFlag)
	if !reflect.DeepEqual(fs, expected) {
		t.Errorf("Testing: FlagSet.Add(%+v); Expected: %+v; Got: %+v", optFlag, expected, fs)
	}

	expected.optFlags[swFlag.name] = swFlag
	fs.Add(swFlag)
	if !reflect.DeepEqual(fs, expected) {
		t.Errorf("Testing: FlagSet.Add(%+v); Expected: %+v; Got: %+v", swFlag, expected, fs)
	}
}

func Test_NewArgSetFrom_InvalidInputs(t *testing.T) {
	data := []interface{}{
		// Test nil as input
		nil,
		// Test non-pointer as input
		*new(bool),
		// Test pointer to a non-struct as input
		new(bool),
		// Test unexported tagged field as input
		&struct {
			field1 int `flagparse:""`
		}{},
		// Test unsupported field type as input
		&struct {
			Field1 int8 `flagparse:""`
		}{},
		// Test invalid key/value as input
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

func Test_NewArgSetFrom_ValidInputs(t *testing.T) {
	args := struct {
		Field0 int // should get ignored
		Field1 int `flagparse:""`           // expected an optional flag
		Field2 int `flagparse:"positional"` // expected a positional flag
	}{Field0: 11, Field1: 22, Field2: 33}

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

// func Test_usage_defaultUsage(t *testing.T) {
// 	args1 := struct {
// 		Pos1 int     `flagparse:"positional,help=pos1 help"`
// 		Pos2 bool    `flagparse:"positional,help=pos2 help"`
// 		Pos3 string  `flagparse:"positional,help=pos3 help"`
// 		Pos4 float64 `flagparse:"positional,help=pos4 help"`
// 		Pos5 []int   `flagparse:"positional,help=pos5 help,nargs=2"`
// 		Opt1 int     `flagparse:"help=opt1 help"`
// 		Opt2 bool    `flagparse:"help=opt2 help"`
// 		Opt3 string  `flagparse:"help=opt3 help"`
// 		Opt4 float64 `flagparse:"help=opt4 help"`
// 		Opt5 []int   `flagparse:"help=opt5 help,nargs=3"`
// 		Sw1  bool    `flagparse:"help=sw1 help,nargs=0"`
// 	}{}

// 	fs, err := NewFlagSetFrom(&args1)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	buf := &bytes.Buffer{}
// 	fs.SetOutput(buf)
// 	fs.Desc = "Description about cmdline usage"
// 	fs.usage()
// }

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
	Pos1 int       `flagparse:"positional,help=pos1 help"`
	Pos2 []float64 `flagparse:"positional,help=pos2 help,nargs=2"`
	Opt1 int       `flagparse:"help=opt1 help"`
	Opt2 []string  `flagparse:"help=opt2 help,nargs=2"`
	Opt3 []float64 `flagparse:"help=opt3 help,nargs=-1"`
	Sw1  bool      `flagparse:"help=sw1 help,nargs=0"`
}

func Test_Parse_InvalidInputs(t *testing.T) {
	cfg := &testConfig{}
	fs, err := NewFlagSetFrom(cfg)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}
	f, _ := os.Create(os.DevNull)
	fs.SetOutput(f)
	data := [][]string{
		{},
		{"not a number"},
		{"10", "1.1"},
		{"10", "1.1", "2.2", "3.3"},
		{"10", "--dummy"},
		{"10", "--opt1", "hello"},
		{"10", "--opt1", "55", "--opt1", "60"},
		{"10", "--opt2", "one"},
		{"10", "--opt3"},
		{"10", "--opt3", "--opt1", "55"},
		{"10", "--opt3", "45.6", "99.99", "not a float64"},
	}
	for _, input := range data {
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
		Opt1: -999,
		Opt2: []string{"hello", "world"},
		Opt3: []float64{5.5},
		Sw1:  false,
	}
	fs, err := NewFlagSetFrom(cfg)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}
	f, _ := os.Create(os.DevNull)
	fs.SetOutput(f)
	data := []struct {
		args     []string
		expected *testConfig
	}{
		{
			args: []string{"10", "1.2", "3.4"},
			expected: &testConfig{
				Pos1: 10,
				Pos2: []float64{1.2, 3.4},
				Opt1: -999,
				Opt2: []string{"hello", "world"},
				Opt3: []float64{5.5},
				Sw1:  false,
			},
		},
		{
			args: []string{"20", "1.2", "3.4", "--sw1", "--opt1", "100", "--opt2", "one", "two", "--opt3", "1.1", "2.2", "3.3"},
			expected: &testConfig{
				Pos1: 20,
				Pos2: []float64{1.2, 3.4},
				Opt1: 100,
				Opt2: []string{"one", "two"},
				Opt3: []float64{1.1, 2.2, 3.3},
				Sw1:  true},
		},
	}
	for _, input := range data {
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

func Test_Parse_PosFlagWithUnlimitedArgs(t *testing.T) {
	type testCfg struct {
		Pos1 []int `flagparse:"positional,nargs=-1"`
	}
	cfg := &testCfg{}
	fs, err := NewFlagSetFrom(cfg)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}
	f, _ := os.Create(os.DevNull)
	fs.SetOutput(f)
	fs.ContinueOnError = true

	good := []struct {
		args     []string
		expected *testCfg
	}{
		{
			args:     []string{"11"},
			expected: &testCfg{Pos1: []int{11}},
		},
		{
			args:     []string{"11", "22", "33", "44", "55"},
			expected: &testCfg{Pos1: []int{11, 22, 33, 44, 55}},
		},
	}
	for _, input := range good {
		fs.CmdArgs = input.args
		if err := fs.Parse(); err != nil {
			t.Errorf("Testing: FlagSet.Parse(); Expected: no error with %q as args; Got: error %q", input.args, err)
		}
		if !reflect.DeepEqual(cfg, input.expected) {
			t.Errorf("Testing: FlagSet.Parse(); Expected: %+v; Got:%+v", input.expected, cfg)
		}
	}

	bad := [][]string{[]string{}, []string{"11", "22", "33", "44abc", "55"}}
	for _, input := range bad {
		fs.CmdArgs = input
		if err := fs.Parse(); err == nil {
			t.Errorf("Testing: FlagSet.Parse(); Expected: error with %q as args; Got: no error", input)
		}
	}
}

func Test_Parse_HelpOption(t *testing.T) {
	cfg := &testConfig{}
	fs, err := NewFlagSetFrom(cfg)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}
	fs.Desc = "flagset description"
	fs.ContinueOnError = true
	f, _ := os.Create(os.DevNull)
	fs.SetOutput(f)
	fs.CmdArgs = []string{helpOptFlag}
	err = fs.Parse()
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
		{helpOptFlag, 1},
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
