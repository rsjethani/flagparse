package flagparse_test

import (
	"fmt"

	"github.com/rsjethani/flagparse"
)

func Example_aPIApproach() {
	cfg := struct {
		Pos1 int
		Pos2 []float64
		Opt1 int
		Opt2 []string
		Sw1  bool
	}{Opt1: -1}

	fs := flagparse.NewFlagSet()

	fl := flagparse.NewIntFlag(&cfg.Pos1, true, "pos1 usage")
	fs.Add("pos1", fl)

	fl = flagparse.NewFloat64ListFlag(&cfg.Pos2, true, "pos2 usage")
	fl.SetNArgs(2)
	fs.Add("pos2", fl)

	fl = flagparse.NewIntFlag(&cfg.Opt1, false, "opt1 usage")
	fs.Add("opt1", fl)

	fl = flagparse.NewStringListFlag(&cfg.Opt2, false, "opt2 usage")
	fl.SetNArgs(2)
	fs.Add("opt2", fl)

	fl = flagparse.NewBoolFlag(&cfg.Sw1, false, "sw1 usage")
	// make this optional flag a switch
	fl.SetNArgs(0)
	fs.Add("sw1", fl)

	fmt.Printf("\nBefore parsing: %+v", cfg)
	fs.CmdArgs = []string{"11", "1.1", "2.2", "--opt1", "22", "--opt2", "hello", "world", "--sw1"}
	fs.Parse()
	fmt.Printf("\nAfter parsing: %+v", cfg)
	// Output:
	// Before parsing: {Pos1:0 Pos2:[] Opt1:-1 Opt2:[] Sw1:false}
	// After parsing: {Pos1:11 Pos2:[1.1 2.2] Opt1:22 Opt2:[hello world] Sw1:true}

}

func Example_structTagApproach() {
	cfg := struct {
		Pos1 int       `flagparse:"positional,usage=pos1 usage"`
		Pos2 []float64 `flagparse:"positional,usage=pos2 usage,nargs=2"`
		Opt1 int       `flagparse:"name=opt1-name,usage=opt1 usage"`
		Opt2 []string  `flagparse:"usage=opt2 usage,nargs=2"`
		Sw1  bool      `flagparse:"usage=sw1 usage,nargs=0"`
	}{Opt1: -1}
	fs, err := flagparse.NewFlagSetFrom(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\nBefore parsing: %+v", cfg)
	fs.CmdArgs = []string{"11", "1.1", "2.2", "--opt1-name", "22", "--opt2", "hello", "world", "--sw1"}
	fs.Parse()
	fmt.Printf("\nAfter parsing: %+v", cfg)
	// Output:
	// Before parsing: {Pos1:0 Pos2:[] Opt1:-1 Opt2:[] Sw1:false}
	// After parsing: {Pos1:11 Pos2:[1.1 2.2] Opt1:22 Opt2:[hello world] Sw1:true}
}
