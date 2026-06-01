package internal

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func PasswordBasic() {
	var password string
	passwordQuestion := pardon.NewPassword(&password).
		Title("Enter your password:")

	if err := passwordQuestion.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entered password is %s%s%s\n", ansi.Green, password, ansi.Reset)

	os.Exit(0)
}
