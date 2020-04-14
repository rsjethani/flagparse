package flagparse

import (
	"fmt"
)

type Flag struct {
	defVal     string
	nArgs      int
	positional bool
	value      Value
	usage      string
}

func newFlag(val Value, pos bool, usage string) *Flag {
	return &Flag{
		nArgs:      1,
		value:      val,
		usage:      usage,
		positional: pos,
		defVal:     val.String(),
	}
}

func (fl *Flag) isSwitch() bool {
	return !fl.positional && fl.nArgs == 0
}

func (fl *Flag) optToSwitch() {
	fl.defVal = ""
}

func (fl *Flag) switchToOpt() {
	fl.defVal = fl.value.String()
}

func (fl *Flag) SetNArgs(n int) error {
	if n == 0 {
		if fl.positional {
			return fmt.Errorf("nargs cannot be 0 for positional flag")
		}
		if fl.nArgs != 0 { // means this is an optional flag which needs to be converted to a switch
			fl.optToSwitch()
		}
	} else if fl.isSwitch() {
		fl.switchToOpt()
	}
	fl.nArgs = n
	return nil
}

func NewFlag(val Value, pos bool, usage string) *Flag {
	return newFlag(val, pos, usage)
}

func NewBoolFlag(val *bool, pos bool, usage string) *Flag {
	return newFlag(newBoolValue(val), pos, usage)
}

func NewBoolListFlag(val *[]bool, pos bool, usage string) *Flag {
	return newFlag(newBoolListValue(val), pos, usage)
}

func NewStringFlag(val *string, pos bool, usage string) *Flag {
	return newFlag(newStringValue(val), pos, usage)
}

func NewStringListFlag(val *[]string, pos bool, usage string) *Flag {
	return newFlag(newStringListValue(val), pos, usage)
}

func NewIntFlag(val *int, pos bool, usage string) *Flag {
	return newFlag(newIntValue(val), pos, usage)
}

func NewIntListFlag(val *[]int, pos bool, usage string) *Flag {
	return newFlag(newIntListValue(val), pos, usage)
}

func NewFloat64Flag(val *float64, pos bool, usage string) *Flag {
	return newFlag(newFloat64Value(val), pos, usage)
}

func NewFloat64ListFlag(val *[]float64, pos bool, usage string) *Flag {
	return newFlag(newFloat64ListValue(val), pos, usage)
}
