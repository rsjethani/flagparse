package flagparse

import (
	"fmt"
)

type Flag struct {
	value      Value
	help       string
	positional bool
	nArgs      int // TODO: convert to string for patterns like '*', '+' etc.
}

func NewPosFlag(value Value, help string) *Flag {
	return &Flag{
		nArgs:      1,
		value:      value,
		help:       help,
		positional: true,
	}
}

func NewOptFlag(value Value, help string) *Flag {
	return &Flag{
		nArgs:      1,
		value:      value,
		help:       help,
		positional: false,
	}
}

func NewSwitchFlag(value Value, help string) *Flag {
	return &Flag{
		nArgs:      0,
		value:      value,
		help:       help,
		positional: false,
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
