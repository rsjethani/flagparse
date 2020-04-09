package main

import (
	"fmt"

	"github.com/rsjethani/flagparse"
)

func main() {
	var pos1 int
	var sw1 bool
	var opt1 string

	fs := flagparse.NewFlagSet()
	fs.Add("pos1", flagparse.NewPosFlag(flagparse.NewInt(&pos1), "pos1 help"))
	fs.Add("opt1", flagparse.NewOptFlag(flagparse.NewString(&opt1), "opt1 help"))
	fs.Add("sw1", flagparse.NewSwitchFlag(flagparse.NewBool(&sw1), "sw1 help"))

	fmt.Println("before parse: ", pos1, opt1, sw1)
	fs.Parse()
	fmt.Println("after parse: ", pos1, opt1, sw1)
}
