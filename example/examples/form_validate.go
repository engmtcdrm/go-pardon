package examples

import (
	"fmt"
	"os"
	"strconv"

	"github.com/engmtcdrm/go-pardon"
)

func FormValidate() {
	continueFlag := true
	age := ""

	f := pardon.NewForm(
		pardon.NewConfirm().
			Title("Are you sure you want to proceed?").
			Value(&continueFlag),
		pardon.NewQuestion().
			Title("How old are you?").
			Validate(func(answer string) error {
				if answer == "" {
					return fmt.Errorf("age cannot be empty")
				}

				answerInt, err := strconv.Atoi(answer)
				if err != nil {
					return fmt.Errorf("invalid age")
				}

				if answerInt < 0 {
					return fmt.Errorf("age cannot be negative")
				}

				return nil
			}).
			Value(&age),
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
