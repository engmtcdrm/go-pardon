package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/pardon"
)

func ConfirmBasic() {
	continueFlag := true

	confirm := pardon.NewConfirm().
		Title("Are you sure you want to proceed?").
		Value(&continueFlag)

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
