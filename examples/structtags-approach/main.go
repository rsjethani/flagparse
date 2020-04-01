package main

import (
	"fmt"

	"github.com/rsjethani/flagparse"
)

func main() {
	config := struct {
		Salute   string  `flagparse:"help=Salutation for the employee"`
		Salary   float64 `flagparse:"help=Employee salary,type=pos"`
		FullName string  `flagparse:"help=Full name of the employee,type=pos"`
		EmpID    []int   `flagparse:"name=emp-id,help=Employee ID for new employee,nargs=3"`
		IsIntern bool    `flagparse:"name=is-intern,help=Is the new employee an intern,type=switch"`
	}{
		EmpID:  []int{100},
		Salute: "Mr.",
	}

	mainSet, err := flagparse.NewFlagSetFrom(&config)
	if err != nil {
		fmt.Println(err)
		return
	}
	mainSet.Desc = "CLI for managing employee database"

	fmt.Printf("\nBEFORE parsing: %+v\n", config)
	// mainSet.CmdArgs = []string{"3.4", "asd", "--salute", "XXX", "--is-intern", "--emp-id", "88888", "345", "33"}
	mainSet.Parse()
	fmt.Printf("\nAFTER parsing: %+v\n", config)
}
