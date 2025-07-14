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

type Color2 struct {
	Name string
	ID   int
	Sub  string
}

func main() {
	selectedItem := Color{}
	items := []gocliselect.Item[Color]{}
	items = append(items, gocliselect.Item[Color]{Key: "Red", Value: Color{Name: "Red", ID: 1}})
	items = append(items, gocliselect.Item[Color]{Key: "Blue", Value: Color{Name: "Blue", ID: 2}})
	items = append(items, gocliselect.Item[Color]{Key: "Green", Value: Color{Name: "Green", ID: 3}})
	items = append(items, gocliselect.Item[Color]{Key: "Yellow", Value: Color{Name: "Yellow", ID: 4}})
	items = append(items, gocliselect.Item[Color]{Key: "Red", Value: Color{Name: "Red", ID: 5}})
	items = append(items, gocliselect.Item[Color]{Key: "Blue", Value: Color{Name: "Blue", ID: 6}})
	items = append(items, gocliselect.Item[Color]{Key: "Green", Value: Color{Name: "Green", ID: 7}})
	items = append(items, gocliselect.Item[Color]{Key: "Yellow", Value: Color{Name: "Yellow", ID: 8}})

	menu := gocliselect.NewSelect[Color]().
		Title(pp.Cyan("Choose a color:")).
		Items(items...).
		Value(&selectedItem)

	menu.ItemSelectColor = pp.Yellow

	if err := menu.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option: %v\n\n", selectedItem)

	selectedItem2 := Color2{}

	items2 := []gocliselect.Item[Color2]{}
	items2 = append(items2, gocliselect.NewSelectItem("Red", Color2{Name: "Red", ID: 1, Sub: "Sub Red"}))
	items2 = append(items2, gocliselect.NewSelectItem("Blue", Color2{Name: "Blue", ID: 2, Sub: "Sub Blue"}))
	items2 = append(items2, gocliselect.NewSelectItem("Green", Color2{Name: "Green", ID: 3, Sub: "Sub Green"}))
	items2 = append(items2, gocliselect.NewSelectItem("Yellow", Color2{Name: "Yellow", ID: 4, Sub: "Sub Yellow"}))

	menu2 := gocliselect.NewSelect[Color2]().
		Title("Choose a color with sub:").
		Items(items2...).
		Value(&selectedItem2)

	menu2.ItemSelectColor = pp.Green

	if err := menu2.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option with sub: %v\n\n", selectedItem2)

	continueFlag := true

	confirm := gocliselect.NewConfirm().
		TitleFunc(func() string {
			return pp.Cyan("Are you sure you want to proceed?")
		}).
		QuestionMarkFunc(func() string {
			return pp.Cyan("[?]")
		}).
		Value(&continueFlag)

	if err := confirm.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if continueFlag {
		fmt.Print("Confirmed!\n\n")
	} else {
		fmt.Print("Cancelled!\n\n")
	}

	favColor := ""
	question := gocliselect.NewQuestion().
		Title("What is your favorite color?").
		QuestionMarkFunc(func() string {
			return pp.Cyan("[?]")
		}).
		Value(&favColor).
		Validate(func(input string) error {
			if input == "" {
				return fmt.Errorf("input cannot be empty")
			}
			return nil
		})

	if err := question.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Your favorite color is: %s\n\n", favColor)

	password := []byte{}
	passwordQuestion := gocliselect.NewPassword().
		Title("Enter your password:").
		Value(&password)

	if err := passwordQuestion.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Your password is: %s\n\n", string(password))

}
