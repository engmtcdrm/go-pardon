package internal

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func ReusingPrompt() {
	favColor := ""
	question := pardon.NewQuestion(&favColor).
		Title("What is your favorite color?").
		TitleFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Green, s, ansi.Reset)
		}).
		Icon("??? ").
		IconFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Blue, s, ansi.Reset)
		}).
		AnswerFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.CyanBg, s, ansi.Reset)
		}).
		Validate(func(input string) error {
			if input == "" {
				return fmt.Errorf("color cannot be empty")
			}

			validColors := []string{"red", "green", "blue", "yellow", "purple", "orange"}
			for _, color := range validColors {
				if input == color {
					return nil
				}
			}
			return fmt.Errorf("invalid color: %s, must be one of: %v", input, validColors)
		})

	if err := question.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	favColor2 := ""
	question = question.
		Value(&favColor2).
		Title("What is your second favorite color?")
		// TitleFunc(func(s string) string {
		// 	return fmt.Sprintf("%s%s%s", ansi.Green, s, ansi.Reset)
		// })

	if err := question.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entered favorite color is '%s%s%s'\n", ansi.Green, favColor, ansi.Reset)
	fmt.Printf("Entered second favorite color is '%s%s%s'\n", ansi.Green, favColor2, ansi.Reset)

	os.Exit(0)
}
