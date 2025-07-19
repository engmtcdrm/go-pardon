package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func PasswordKitchesink() {
	password := []byte{}
	passwordQuestion := pardon.NewPassword().
		Title("Enter your password:").
		TitleFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Green, s, ansi.Reset)
		}).
		Icon("§ ").
		IconFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Red, s, ansi.Reset)
		}).
		AnswerFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.BlueBg, s, ansi.Reset)
		}).
		Validate(func(input []byte) error {
			if len(input) < 8 {
				return fmt.Errorf("password must be at least 8 characters long")
			}
			if !containsSpecialChar(input) {
				return fmt.Errorf("password must contain at least one special character")
			}
			return nil
		}).
		Value(&password)

	if err := passwordQuestion.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entered password is %s%s%s\n", ansi.Green, string(password), ansi.Reset)

	os.Exit(0)
}
