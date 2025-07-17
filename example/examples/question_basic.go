package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-pardon"
)

func QuestionBasic() {
	favColor := ""
	question := pardon.NewQuestion().
		Title("What is your favorite color?").
		Value(&favColor).
		Validate(func(input string) error {
			if input == "" {
				return fmt.Errorf("color cannot be empty")
			}
			return nil
		})

	if err := question.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entered favorite color is '%s'\n", favColor)

	os.Exit(0)
}
