package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/pardon"
)

func PasswordBasic() {
	password := []byte{}
	passwordQuestion := pardon.NewPassword().
		Title("Enter your password:").
		Value(&password)

	if err := passwordQuestion.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entered password is '%s'\n", string(password))

	os.Exit(0)
}
