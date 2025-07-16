package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/pardon"
)

func SelectStruct() {
	type Color struct {
		Name string
		ID   int
	}

	selectedColor := Color{}
	colors := []pardon.Option[Color]{}
	colors = append(colors, pardon.Option[Color]{Key: "Red", Value: Color{Name: "Red", ID: 1}})
	colors = append(colors, pardon.Option[Color]{Key: "Blue", Value: Color{Name: "Blue", ID: 2}})
	colors = append(colors, pardon.Option[Color]{Key: "Green", Value: Color{Name: "Green", ID: 3}})
	colors = append(colors, pardon.Option[Color]{Key: "Yellow", Value: Color{Name: "Yellow", ID: 4}})

	menu := pardon.NewSelect[Color]().
		Title("Choose a color:").
		Options(colors...).
		Value(&selectedColor)

	if err := menu.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option: %v\n", selectedColor)
	os.Exit(0)
}
