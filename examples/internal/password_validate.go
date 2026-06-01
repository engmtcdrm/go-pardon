package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func PasswordValidate() {
	var password string
	passwordQuestion := pardon.NewPassword(&password).
		Title("Enter your password:").
		ValidateFunc(func(input string) error {
			if len(input) < 8 {
				return fmt.Errorf("password must be at least 8 characters long")
			}
			if !containsSpecialChar(input) {
				return fmt.Errorf("password must contain at least one special character")
			}
			return nil
		})

	if err := passwordQuestion.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entered password is %s%s%s\n", ansi.Green, password, ansi.Reset)

	os.Exit(0)
}

func containsSpecialChar(input string) bool {
	specialChars := "!@#$%^&*()-_=+[]{}|;:',.<>?/"
	return strings.ContainsAny(input, specialChars)
}
