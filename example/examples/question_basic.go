package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func QuestionBasic() {
	favColor := ""
	question := pardon.NewQuestion().
		Title("What is your name?").
		Value(&favColor)

	if err := question.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Your name is %s%s%s\n", ansi.Green, favColor, ansi.Reset)

	os.Exit(0)
}
