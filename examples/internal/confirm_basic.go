package internal

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-pardon"
)

func ConfirmBasic() {
	continueFlag := true

	confirm := pardon.NewConfirm(&continueFlag).
		Title("Are you sure you want to proceed?")

	if err := confirm.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if continueFlag {
		fmt.Println("Proceeding!")
	} else {
		fmt.Println("Stopping!")
	}

	os.Exit(0)
}
