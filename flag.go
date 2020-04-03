package flagparse

import (
	"fmt"
)

type Flag struct {
	// TODO: convert to string for patterns like '*', '+' etc.
	nArgs      int
	positional bool
	value      Value
	name       string
	help       string
}

func NewPosFlag(name string, value Value, help string) *Flag {
	return &Flag{
		name:       name,
		nArgs:      1,
		value:      value,
		help:       help,
		positional: true,
	}
}

func NewOptFlag(name string, value Value, help string) *Flag {
	return &Flag{
		name:       name,
		nArgs:      1,
		value:      value,
		help:       help,
		positional: false,
	}
}

func NewSwitchFlag(name string, value Value, help string) *Flag {
	fl := NewOptFlag(name, value, help)
	fl.nArgs = 0
	return fl
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
