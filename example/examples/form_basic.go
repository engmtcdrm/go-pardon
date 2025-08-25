package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-pardon"
)

func FormBasic() {
	continueFlag := true
	age := ""

	f := pardon.NewForm(
		pardon.NewConfirm(&continueFlag).
			Title("Are you sure you want to proceed?"),
		pardon.NewQuestion(&age).
			Title("How old are you?"),
	)

	if err := f.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if continueFlag {
		fmt.Println("Proceeding!")
	} else {
		fmt.Println("Stopping!")
	}

	fmt.Println("You are", age, "years old.")

	os.Exit(0)
}
