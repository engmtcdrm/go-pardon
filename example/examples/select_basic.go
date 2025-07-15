package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/pardon"
)

func SelectBasic() {
	var selectedColor int
	colors := []pardon.Option[int]{}
	colors = append(colors, pardon.Option[int]{Key: "Red", Value: 1})
	colors = append(colors, pardon.Option[int]{Key: "Blue", Value: 2})
	colors = append(colors, pardon.Option[int]{Key: "Green", Value: 3})
	colors = append(colors, pardon.Option[int]{Key: "Yellow", Value: 4})

	selectPrompt := pardon.NewSelect[int]().
		Title("Choose a color:").
		Options(colors...).
		Value(&selectedColor)

	if err := selectPrompt.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option: %v\n", selectedColor)
	os.Exit(0)
}
