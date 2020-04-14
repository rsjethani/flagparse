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
	fs.Add("pos1", flagparse.NewIntFlag(&pos1, true, "pos1 usage"))
	fs.Add("opt1", flagparse.NewStringFlag(&opt1, false, "opt1 usage"))

	sw1Flag := flagparse.NewBoolFlag(&sw1, false, "sw1 usage")
	sw1Flag.SetNArgs(0)
	fs.Add("sw1", sw1Flag)

	fmt.Println("before parse: ", pos1, opt1, sw1)
	fs.Parse()
	fmt.Println("after parse: ", pos1, opt1, sw1)
}
