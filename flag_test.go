package flagparse

import (
	"testing"
)

func TestNewPosFlag(t *testing.T) {
	fl := NewPosFlag(nil, "help string")
	if !fl.positional || fl.value != nil || fl.nArgs != 1 || fl.help != "help string" {
		t.Errorf("Expected: positional=true, value=nil, nargs=1, help=help string; Got: %+v", fl)
	}
}

func TestNewOptFlag(t *testing.T) {
	fl := NewOptFlag(nil, "help string")
	if fl.positional || fl.value != nil || fl.nArgs != 1 || fl.help != "help string" {
		t.Errorf("Expected: positional=false, value=nil, nargs=1, help=help string; Got: %+v", fl)
	}
}

func TestNewSwitchFlag(t *testing.T) {
	fl := NewSwitchFlag(nil, "help string")
	if fl.positional || fl.value != nil || fl.nArgs != 0 || fl.help != "help string" {
		t.Errorf("Expected: positional=false, value=nil, nargs=0, help=help string; Got: %+v", fl)
	}
}

func TestSetNArgs(t *testing.T) {
	fl := NewPosFlag(nil, "")
	if err := fl.SetNArgs(10); err != nil || fl.nArgs != 10 {
		t.Errorf("Expected: for positional flag %[1]T.SetNArgs(10) suceeds with nil error setting %[1]T.nArgs==10; Got: error", fl)
	}
	if err := fl.SetNArgs(0); err == nil {
		t.Errorf("Expected: for positional flag %T.SetNArgs(0) results in error; Got: nil error", fl)
	}

	fl = NewOptFlag(nil, "")
	if err := fl.SetNArgs(10); err != nil || fl.nArgs != 10 {
		t.Errorf("Expected: for optional flag %[1]T.SetNArgs(10) suceeds with nil error setting %[1]T.nArgs==10; Got: error", fl)
	}
	if err := fl.SetNArgs(0); err != nil || fl.nArgs != 0 {
		t.Errorf("Expected: for optional flag %[1]T.SetNArgs(0) suceeds with no error setting %[1]T.nArgs==0; Got: error", fl)
	}
}
