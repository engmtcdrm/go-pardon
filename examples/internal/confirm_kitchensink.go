package internal

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func ConfirmKitchensink() {
	continueFlag := true

	confirm := pardon.NewConfirm(&continueFlag).
		Title("Are you sure you want to proceed?").
		TitleFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.MagentaBg, s, ansi.Reset)
		}).
		Icon("✔ ").
		IconFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Red, s, ansi.Reset)
		}).
		AnswerFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.BlueBg, s, ansi.Reset)
		})

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
