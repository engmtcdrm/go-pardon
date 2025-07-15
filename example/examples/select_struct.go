package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/gocliselect"
)

func SelectStruct() {
	type Color struct {
		Name string
		ID   int
	}

	selectedColor := Color{}
	colors := []gocliselect.Option[Color]{}
	colors = append(colors, gocliselect.Option[Color]{Key: "Red", Value: Color{Name: "Red", ID: 1}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Blue", Value: Color{Name: "Blue", ID: 2}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Green", Value: Color{Name: "Green", ID: 3}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Yellow", Value: Color{Name: "Yellow", ID: 4}})

	menu := gocliselect.NewSelect[Color]().
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
