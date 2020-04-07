package flagparse

import (
	"reflect"
	"testing"
)

func Test_NewSwitchFlag(t *testing.T) {
	var testVar int = 100
	val := NewInt(&testVar)
	expected := &Flag{
		nArgs:      0,
		defVal:     "",
		value:      val,
		positional: false,
		name:       "flag-name",
		help:       "help message",
	}
	got := NewSwitchFlag(expected.name, expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewSwitchFlag(%q, %#v, %q); Expected: %#v; Got: %#v", expected.name, expected.value, expected.help, expected, got)
	}
}

func Test_NewOptFlag(t *testing.T) {
	var testVar int = 100
	val := NewInt(&testVar)
	expected := &Flag{
		nArgs:      1,
		defVal:     val.String(),
		value:      val,
		positional: false,
		name:       "flag-name",
		help:       "help message",
	}
	got := NewOptFlag(expected.name, expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewOptFlag(%q, %#v, %q); Expected: %#v; Got: %#v", expected.name, expected.value, expected.help, expected, got)
	}
}

func Test_NewPosFlag(t *testing.T) {
	var testVar int = 100
	val := NewInt(&testVar)
	expected := &Flag{
		nArgs:      1,
		defVal:     val.String(),
		value:      val,
		positional: true,
		name:       "flag-name",
		help:       "help message",
	}
	got := NewPosFlag(expected.name, expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewPosFlag(%q, %#v, %q); Expected: %#v; Got: %#v", expected.name, expected.value, expected.help, expected, got)
	}
}

func Test_SetNArgs(t *testing.T) {

}
