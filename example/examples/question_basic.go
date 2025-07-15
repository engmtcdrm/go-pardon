package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/gocliselect"
)

func QuestionBasic() {
	favColor := ""
	question := gocliselect.NewQuestion().
		Title("What is your favorite color?").
		Value(&favColor)

	if err := question.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entered favorite color is '%s'\n", favColor)

	os.Exit(0)
}
