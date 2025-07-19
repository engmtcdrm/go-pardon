package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
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

	fmt.Printf("Entered password is %s%s%s\n", ansi.Green, string(password), ansi.Reset)

	os.Exit(0)
}
