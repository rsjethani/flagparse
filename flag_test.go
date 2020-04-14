package flagparse

import (
	"reflect"
	"testing"
)

func Test_SetNArgs_WithPositionalFlag(t *testing.T) {
	testVar := 100

	posFlag := NewIntFlag(&testVar, true, "")
	if err := posFlag.SetNArgs(0); err == nil {
		t.Errorf("Testing: Flag.SetNArgs(0); Expected: error; Got: no error")
	}

	expected := *posFlag
	expected.nArgs = 10
	if err := posFlag.SetNArgs(10); err != nil || !reflect.DeepEqual(*posFlag, expected) {
		t.Errorf("Testing: Flag.SetNArgs(10); Expected: no error and %#v; Got: error %v, %#v", expected, err, *posFlag)
	}
}

func Test_SetNArgs_WithOptionalFlag(t *testing.T) {
	testVar := 100

	optFlag := NewIntFlag(&testVar, false, "")
	expected := *optFlag
	expected.nArgs = 0
	expected.optToSwitch()
	if err := optFlag.SetNArgs(0); err != nil || !reflect.DeepEqual(*optFlag, expected) {
		t.Errorf("Testing: Flag.SetNArgs(0); Expected: no error and %#v; Got: error %v, %#v", expected, err, *optFlag)
	}

	expected = *optFlag
	expected.nArgs = 10
	expected.switchToOpt()
	if err := optFlag.SetNArgs(10); err != nil || !reflect.DeepEqual(*optFlag, expected) {
		t.Errorf("Testing: Flag.SetNArgs(10); Expected: no error and %#v; Got: error %v, %#v", expected, err, *optFlag)
	}
}
