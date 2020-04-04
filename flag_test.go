package flagparse

import (
	"reflect"
	"testing"
)

func TestSwitchFlag(t *testing.T) {
	var testVar int = 100
	val := NewInt(&testVar)
	expected := &Flag{
		name:  "flag-name",
		value: val,
		help:  "help message",
	}
	got := NewSwitchFlag(expected.name, expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewSwitchFlag(); Expected: %#v; Got: %#v", expected, got)
	}
	if !got.isSwitch() {
		t.Errorf("Testing: Flag.isSwitch(); Expected: true; Got: false")
	}
	err := got.SetNArgs(5)
	if err != nil || got.nArgs != 5 || got.isSwitch() {
		t.Errorf("Testing: Flag.SetNargs(5); Expected: no error, Flag.nArgs==5, Flag.isSwitch()==false; Got: %v error, Flag.nArgs==%v, Flag.isSwitch()==%v", err, got.nArgs, got.isSwitch())
	}
}

func TestOptFlag(t *testing.T) {
	var testVar int = 100
	val := NewInt(&testVar)
	expected := &Flag{
		name:   "flag-name",
		value:  val,
		help:   "help message",
		nArgs:  1,
		defVal: val.String(),
	}
	got := NewOptFlag(expected.name, expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewOptFlag(); Expected: %#v; Got: %#v", expected, got)
	}
	if got.isSwitch() {
		t.Errorf("Testing: Flag.isSwitch(); Expected: false; Got: true")
	}
	err := got.SetNArgs(0)
	if err != nil || got.nArgs != 0 || !got.isSwitch() {
		t.Errorf("Testing: Flag.SetNargs(0); Expected: no error, Flag.nArgs==0, Flag.isSwitch()==true; Got: %v error, Flag.nArgs==%v, Flag.isSwitch()==%v", err, got.nArgs, got.isSwitch())
	}
}

func TestPosFlag(t *testing.T) {
	var testVar int = 100
	val := NewInt(&testVar)
	expected := &Flag{
		name:       "flag-name",
		value:      val,
		help:       "help message",
		nArgs:      1,
		positional: true,
		defVal:     val.String(),
	}
	got := NewPosFlag(expected.name, expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewPosFlag(); Expected: %#v; Got: %#v", expected, got)
	}
	if got.isSwitch() {
		t.Errorf("Testing: Flag.isSwitch(); Expected: false; Got: true")
	}
	err := got.SetNArgs(0)
	if err == nil {
		t.Errorf("Testing: Flag.SetNargs(0); Expected: error; Got: nil")
	}
	err = got.SetNArgs(5)
	if err != nil || got.nArgs != 5 || !got.positional {
		t.Errorf("Testing: Flag.SetNargs(5); Expected: no error, Flag.nArgs==5, Flag.positional==true; Got: %v error, Flag.nArgs==%v, Flag.positional==%v", err, got.nArgs, got.positional)
	}
}
