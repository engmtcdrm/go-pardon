package main

import (
	"fmt"

	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/gocliselect"
)

type Color struct {
	Name string
	ID   int
}

func main() {
	menu := gocliselect.NewSelect(pp.Cyan("Choose a color:"))

	menu.ItemSelectColor = pp.Yellow

	menu.AddItem("Red", Color{Name: "Red", ID: 1})
	menu.AddItem("Blue", Color{Name: "Blue", ID: 2})
	menu.AddItem("Green", Color{Name: "Green", ID: 3})
	menu.AddItem("Yellow", Color{Name: "Yellow", ID: 4})
	menu.AddItem("Red", Color{Name: "Red", ID: 5})
	menu.AddItem("Blue", Color{Name: "Blue", ID: 6})
	menu.AddItem("Green", Color{Name: "Green", ID: 7})
	menu.AddItem("Yellow", Color{Name: "Yellow", ID: 8})

	result, err := menu.Display()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option: %v\n", result)
}
