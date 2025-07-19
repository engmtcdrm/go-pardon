package examples

import (
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func PasswordValidate() {
	password := []byte{}
	passwordQuestion := pardon.NewPassword().
		Title("Enter your password:").
		Value(&password).
		Validate(func(input []byte) error {
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

	fmt.Printf("Entered password is %s%s%s\n", ansi.Green, string(password), ansi.Reset)

	os.Exit(0)
}

func containsSpecialChar(input []byte) bool {
	specialChars := "!@#$%^&*()-_=+[]{}|;:',.<>?/"
	for _, char := range input {
		if strings.ContainsRune(specialChars, rune(char)) {
			return true
		}
	}
	return false
}
