package flagparse_test

import (
	"fmt"

	"github.com/rsjethani/flagparse"
)

func Example_aPIApproach() {
	cfg := struct {
		Pos1 int
		Opt1 string
		Opt2 []float64
		Sw1  bool
	}{
		// default values for optional variables
		Opt1: "hello",
		Opt2: []float64{1.1, 2.2},
	}

	fs := flagparse.NewFlagSet()

	fl := flagparse.NewIntFlag(&cfg.Pos1, true, "pos1 usage")
	fs.Add(fl, "pos1")

	fl = flagparse.NewStringFlag(&cfg.Opt1, false, "opt1 usage")
	// both --opt1 and -o names can be used to set this flag
	fs.Add(fl, "--opt1", "-o")

	fl = flagparse.NewFloat64ListFlag(&cfg.Opt2, false, "opt2 usage")
	fl.SetNArgs(2)
	fs.Add(fl, "--opt2")

	fl = flagparse.NewBoolFlag(&cfg.Sw1, false, "sw1 usage")
	// make this optional flag a switch
	fl.SetNArgs(0)
	fs.Add(fl, "-s")

	fmt.Printf("\nBefore parsing: %+v", cfg)
	fs.CmdArgs = []string{"11", "--opt1", "hello world!", "--opt2", "33.33", "44.44", "-s"}
	fs.Parse()
	fmt.Printf("\nAfter parsing: %+v", cfg)
	// Output:
	// Before parsing: {Pos1:0 Opt1:hello Opt2:[1.1 2.2] Sw1:false}
	// After parsing: {Pos1:11 Opt1:hello world! Opt2:[33.33 44.44] Sw1:true}
}

func Example_structTagApproach() {
	cfg := struct {
		Pos1 int       `flagparse:"usage=pos1 usage"`
		Opt1 string    `flagparse:"name=--opt1:-o,usage=opt1 usage"`
		Opt2 []float64 `flagparse:"name=--opt2,usage=opt2 usage,nargs=2"`
		Sw1  bool      `flagparse:"name=-s,usage=sw1 usage,nargs=0"`
	}{
		// default values for optional variables
		Opt1: "hello",
		Opt2: []float64{1.1, 2.2},
	}
	fs, err := flagparse.NewFlagSetFrom(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\nBefore parsing: %+v", cfg)
	fs.CmdArgs = []string{"11", "--opt1", "hello world!", "--opt2", "33.33", "44.44", "-s"}
	fs.Parse()
	fmt.Printf("\nAfter parsing: %+v", cfg)
	// Output:
	// Before parsing: {Pos1:0 Opt1:hello Opt2:[1.1 2.2] Sw1:false}
	// After parsing: {Pos1:11 Opt1:hello world! Opt2:[33.33 44.44] Sw1:true}
}
