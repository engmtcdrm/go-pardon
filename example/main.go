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
	items := []gocliselect.SelectItem[Color]{}

	items = append(items, gocliselect.SelectItem[Color]{Key: "Red", Value: Color{Name: "Red", ID: 1}})
	items = append(items, gocliselect.SelectItem[Color]{Key: "Blue", Value: Color{Name: "Blue", ID: 2}})
	items = append(items, gocliselect.SelectItem[Color]{Key: "Green", Value: Color{Name: "Green", ID: 3}})
	items = append(items, gocliselect.SelectItem[Color]{Key: "Yellow", Value: Color{Name: "Yellow", ID: 4}})
	items = append(items, gocliselect.SelectItem[Color]{Key: "Red", Value: Color{Name: "Red", ID: 5}})
	items = append(items, gocliselect.SelectItem[Color]{Key: "Blue", Value: Color{Name: "Blue", ID: 6}})
	items = append(items, gocliselect.SelectItem[Color]{Key: "Green", Value: Color{Name: "Green", ID: 7}})
	items = append(items, gocliselect.SelectItem[Color]{Key: "Yellow", Value: Color{Name: "Yellow", ID: 8}})

	menu := gocliselect.NewSelect[Color]().
		Title(pp.Cyan("Choose a color:")).
		Items(items...)

	menu.ItemSelectColor = pp.Yellow

	result, err := menu.Ask()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option: %v\n", result)

	items2 := []gocliselect.SelectItem[Color2]{}
	items2 = append(items2, gocliselect.NewSelectItem("Red", Color2{Name: "Red", ID: 1, Sub: "Sub Red"}))
	items2 = append(items2, gocliselect.NewSelectItem("Blue", Color2{Name: "Blue", ID: 2, Sub: "Sub Blue"}))
	items2 = append(items2, gocliselect.NewSelectItem("Green", Color2{Name: "Green", ID: 3, Sub: "Sub Green"}))
	items2 = append(items2, gocliselect.NewSelectItem("Yellow", Color2{Name: "Yellow", ID: 4, Sub: "Sub Yellow"}))

	menu2 := gocliselect.NewSelect[Color2]().
		TitleFunc(func() string {
			return pp.Cyan("Choose a color with sub:")
		}).
		Items(items2...)

	menu2.ItemSelectColor = pp.Green

	result2, err := menu2.Ask()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option with sub: %v\n", result2)

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
		fmt.Println("Confirmed!")
	} else {
		fmt.Println("Cancelled!")
	}

	favColor := ""
	question := gocliselect.NewQuestion().
		TitleFunc(func() string {
			return pp.Cyan("What is your favorite color?")
		}).
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

	fmt.Printf("Your favorite color is: %s\n", favColor)

}
