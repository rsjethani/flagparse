package main

import (
	"fmt"

	"github.com/rsjethani/flagparse"
)

func main() {
	config := struct {
		Salute   string  `flagparse:"help=Salutation for the employee"`
		Salary   float64 `flagparse:"positional,help=Employee salary"`
		FullName string  `flagparse:"positional,help=Full name of the employee"`
		EmpID    []int   `flagparse:"name=emp-id,help=Employee ID for new employee,nargs=3"`
		Intern   bool    `flagparse:"help=Is the new employee an intern,nargs=0"`
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
	mainSet.Parse()
	fmt.Printf("\nAFTER parsing: %+v\n", config)
}
