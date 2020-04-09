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
		help:       "help message",
	}
	got := NewSwitchFlag(expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewSwitchFlag(%#v, %q); Expected: %#v; Got: %#v", expected.value, expected.help, expected, got)
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
		help:       "help message",
	}
	got := NewOptFlag(expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewOptFlag(%#v, %q); Expected: %#v; Got: %#v", expected.value, expected.help, expected, got)
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
		help:       "help message",
	}
	got := NewPosFlag(expected.value, expected.help)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Testing: NewPosFlag(%#v, %q); Expected: %#v; Got: %#v", expected.value, expected.help, expected, got)
	}
}

func Test_SetNArgs_WithPositionalFlag(t *testing.T) {
	testVar := 100
	testVal := NewInt(&testVar)

	posFlag := NewPosFlag(testVal, "help")
	if err := posFlag.SetNArgs(0); err == nil {
		t.Errorf("Testing: Flag.SetNArgs(0); Expected: error; Got: no error")
	}

	posFlag = NewPosFlag(testVal, "help")
	expected := *posFlag
	expected.nArgs = 10
	if err := posFlag.SetNArgs(10); err != nil || !reflect.DeepEqual(*posFlag, expected) {
		t.Errorf("Testing: Flag.SetNArgs(10); Expected: no error and %#v; Got: error %v, %#v", expected, err, *posFlag)
	}
}

func Test_SetNArgs_WithOptionalFlag(t *testing.T) {
	testVar := 100
	testVal := NewInt(&testVar)

	optFlag := NewOptFlag(testVal, "help")
	expected := *optFlag
	expected.nArgs = 0
	expected.optToSwitch()
	if err := optFlag.SetNArgs(0); err != nil || !reflect.DeepEqual(*optFlag, expected) {
		t.Errorf("Testing: Flag.SetNArgs(0); Expected: no error and %#v; Got: error %v, %#v", expected, err, *optFlag)
	}

	optFlag = NewOptFlag(testVal, "help")
	expected = *optFlag
	expected.nArgs = 10
	if err := optFlag.SetNArgs(10); err != nil || !reflect.DeepEqual(*optFlag, expected) {
		t.Errorf("Testing: Flag.SetNArgs(10); Expected: no error and %#v; Got: error %v, %#v", expected, err, *optFlag)
	}
}

func Test_SetNArgs_WithSwitchFlag(t *testing.T) {
	testVar := 100
	testVal := NewInt(&testVar)

	swFlag := NewSwitchFlag(testVal, "help")
	expected := *swFlag
	if err := swFlag.SetNArgs(0); err != nil || !reflect.DeepEqual(*swFlag, expected) {
		t.Errorf("Testing: Flag.SetNArgs(0); Expected: no error and %#v; Got: error %v, %#v", expected, err, *swFlag)
	}

	swFlag = NewSwitchFlag(testVal, "help")
	expected = *swFlag
	expected.nArgs = 10
	expected.switchToOpt()
	if err := swFlag.SetNArgs(10); err != nil || !reflect.DeepEqual(*swFlag, expected) {
		t.Errorf("Testing: Flag.SetNArgs(10); Expected: no error and %#v; Got: error %v, %#v", expected, err, *swFlag)
	}
}
