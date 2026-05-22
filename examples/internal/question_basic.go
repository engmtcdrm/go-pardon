package internal

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func QuestionBasic() {
	name := ""
	question := pardon.NewQuestion(&name).
		Title("What is your name?")

	if err := question.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Your name is %s%s%s\n", ansi.Green, name, ansi.Reset)

	os.Exit(0)
}
