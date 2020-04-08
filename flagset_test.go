package flagparse

import (
	"os"
	"reflect"
	"testing"
)

func Test_NewFlagSet(t *testing.T) {
	flagSet := NewFlagSet()
	expected := &FlagSet{
		OptFlagPrefix: defaultOptFlagPrefix,
		optFlags:      make(map[string]*Flag),
		usageOut:      os.Stderr,
		name:          os.Args[0],
		CmdArgs:       os.Args[1:],
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

type testArgs struct {
	Pos1 int `flagparse:"positional,help=pos1 help"`
	// Pos2 []float64 `flagparse:"positional,help=pos2 help,nargs=2"`
	Opt1 int       `flagparse:"help=opt1 help"`
	Opt2 []string  `flagparse:"help=opt2 help,nargs=3"`
	Opt3 []float64 `flagparse:"help=opt3 help,nargs=-1"`
	Sw1  bool      `flagparse:"help=sw1 help,nargs=0"`
}

func Test_parse_InvalidInputs(t *testing.T) {
	args := &testArgs{}
	fs, err := NewFlagSetFrom(args)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}
	data := [][]string{
		{},
		{"not a number"},
		{"10", "11.1"},
		{"10", "--dummy"},
		{"10", "--opt1", "hello"},
		{"10", "--opt1", "55", "--opt1", "60"},
		{"10", "--opt2", "one"},
		{"10", "--opt3", "45.6", "99.99", "not a float64"},
	}
	for _, input := range data {
		fs.CmdArgs = input
		if err := fs.parse(); err == nil {
			t.Errorf("Testing: FlagSet.parse(); Expected error with %+v as args;Got: no error", input)
		}
	}
}
