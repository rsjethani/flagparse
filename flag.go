package flagparse

import (
	"fmt"
)

type Flag struct {
	// TODO: convert to string for patterns like '*', '+' etc.
	defVal     string
	nArgs      int
	positional bool
	value      Value
	name       string
	help       string
}

func NewPosFlag(name string, value Value, help string) *Flag {
	fl := NewSwitchFlag(name, value, help)
	fl.nArgs = 1
	fl.positional = true
	fl.defVal = value.String()
	return fl
}

func NewOptFlag(name string, value Value, help string) *Flag {
	fl := NewSwitchFlag(name, value, help)
	fl.nArgs = 1
	fl.defVal = value.String()
	return fl
}

func NewSwitchFlag(name string, value Value, help string) *Flag {
	return &Flag{
		name:  name,
		value: value,
		help:  help,
	}
}

func (fl *Flag) isSwitch() bool {
	return !fl.positional && fl.nArgs == 0
}

func (fl *Flag) SetNArgs(n int) error {
	if n == 0 && fl.positional {
		return fmt.Errorf("nargs cannot be 0 for positional flag")
	}
	fl.nArgs = n
	return nil
}
