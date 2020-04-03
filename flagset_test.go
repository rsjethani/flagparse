package flagparse

import (
	"testing"
)

func TestNewFlagSet(t *testing.T) {
	argset := NewFlagSet()
	if argset.OptFlagPrefix != defaultOptFlagPrefix {
		t.Errorf("testing: NewFlagSet(); expected: argset.OptFlagPrefix==%#v; got: %#v", defaultOptFlagPrefix, argset.OptFlagPrefix)
	}

	if argset.optFlags == nil {
		t.Errorf("testing: NewFlagSet(); expected: argset.optFlags!=nil; got: nil")
	}
	if len(argset.posFlags) != 0 {
		t.Errorf("testing: DefaultArgSet(); expected: len(argset.posArgs)==0; got: len(argset.posArgs)==%d", len(argset.posFlags))
	}
}

func TestAddFlagToFlagSet(t *testing.T) {
	fs := NewFlagSet()
	fs.Add(nil)
	if len(fs.posFlags) != 0 || fs.optFlags["--dummy"] != nil {
		t.Errorf(`testing: argset.Add("dummy", nil); expected: no positional/optional flag named 'dummy should get added; got: 'dummy' got added`)
	}

	fs = NewFlagSet()
	fs.Add(NewPosFlag("", nil, ""))
	if len(fs.posFlags) == 0 {
		t.Errorf(`testing: argset.AddArgument("pos1", NewPosArg(nil, "")); expected: argset.posArgs[0].name == "pos1"; got: len(argset.posArgs) == 0`)
	}

	fs = NewFlagSet()
	fs.Add(NewOptFlag("", nil, ""))
	if len(fs.optFlags) == 0 {
		t.Errorf(`testing: argset.AddArgument("opt1", NewOptArg(nil, "")); expected: argset.optArgs["opt1"] != nil; got: argset.optArgs["opt1"] == nil`)
	}
}

func TestNewArgSetFromInvalidInputs(t *testing.T) {
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
		// Test invalid tag/value as input
		&struct {
			Field1 int `flagparse:"type=xxx"`
		}{},
	}
	for _, input := range data {
		if argset, err := NewFlagSetFrom(input); argset != nil || err == nil {
			t.Errorf("testing: NewArgSet(%#v); expected: (nil, error); got: (%v, %#v)", input, argset, err)
		}
	}
}

func TestNewArgSetFromValidInputs(t *testing.T) {
	// Test skipping of untagged fields
	args1 := struct {
		Field1 int // no 'flagparse' tag hence should be skipped
	}{}
	argset, err := NewFlagSetFrom(&args1)
	if err != nil {
		t.Errorf("testing: NewArgSet(%#v); expected: non-nil *ArgSet and nil error; got: %v", args1, err)
	}
	if len(argset.posFlags) != 0 || len(argset.optFlags) != 0 {
		t.Errorf("testing: NewArgSet(%#v); expected: no arguments except --help in argset; got: %#v", &args1, argset)
	}

	// Test parsing of tagged fields and no error with valid key/values
	args2 := struct {
		Field1 int `flagparse:""`         // optional argument
		Field2 int `flagparse:"type=pos"` // positional argument
	}{}
	argset, err = NewFlagSetFrom(&args2)
	if err != nil {
		t.Errorf("testing: NewArgSet(%#v); expected: non-nil *ArgSet and nil error; got: %v", args2, err)
	}
	if len(argset.posFlags) == 0 || len(argset.optFlags) == 0 {
		t.Errorf("testing: NewArgSet(%#v); expected: 1 optional and 1 positional arguments in argset; got: %#v", &args2, argset)
	}
}

func TestUsage(t *testing.T) {
	args1 := struct {
		Pos1 int     `flagparse:"type=pos,help=pos1 help"`
		Pos2 bool    `flagparse:"type=pos,help=pos2 help"`
		Pos3 string  `flagparse:"type=pos,help=pos3 help"`
		Pos4 float64 `flagparse:"type=pos,help=pos4 help"`
		Pos5 []int   `flagparse:"type=pos,help=pos5 help,nargs=2"`
		Opt1 int     `flagparse:"help=opt1 help"`
		Opt2 bool    `flagparse:"help=opt2 help"`
		Opt3 string  `flagparse:"help=opt3 help"`
		Opt4 float64 `flagparse:"help=opt4 help"`
		Opt5 []int   `flagparse:"help=opt5 help,nargs=3"`
		Sw1  bool    `flagparse:"type=switch,help=sw1 help"`
	}{}

	fs, err := NewFlagSetFrom(&args1)
	if err != nil {
		t.Error(err)
	}
	fs.Desc = "Description about cmdline usage"
	fs.usage()
}
