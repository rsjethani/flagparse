package flagparse

import (
	"fmt"
)

type Flag struct {
	defVal     string
	nArgs      int
	positional bool
	value      Value
	help       string
}

func NewPosFlag(value Value, help string) *Flag {
	fl := NewSwitchFlag(value, help)
	fl.nArgs = 1
	fl.positional = true
	fl.defVal = value.String()
	return fl
}

func NewOptFlag(value Value, help string) *Flag {
	fl := NewSwitchFlag(value, help)
	fl.nArgs = 1
	fl.defVal = value.String()
	return fl
}

func NewSwitchFlag(value Value, help string) *Flag {
	return &Flag{
		value: value,
		help:  help,
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
