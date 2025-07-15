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
	selectedColor := Color{}
	colors := []gocliselect.Option[Color]{}
	colors = append(colors, gocliselect.Option[Color]{Key: "Red", Value: Color{Name: "Red", ID: 1}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Blue", Value: Color{Name: "Blue", ID: 2}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Green", Value: Color{Name: "Green", ID: 3}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Yellow", Value: Color{Name: "Yellow", ID: 4}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Red", Value: Color{Name: "Red", ID: 5}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Blue", Value: Color{Name: "Blue", ID: 6}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Green", Value: Color{Name: "Green", ID: 7}})
	colors = append(colors, gocliselect.Option[Color]{Key: "Yellow", Value: Color{Name: "Yellow", ID: 8}})

	menu := gocliselect.NewSelect[Color]().
		Title(pp.Cyan("Choose a color:")).
		Options(colors...).
		Value(&selectedColor)

	menu.OptionSelectColor = pp.Yellow

	if err := menu.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option: %v\n\n", selectedColor)

	selectedColor2 := Color2{}

	colors2 := []gocliselect.Option[Color2]{}
	colors2 = append(colors2, gocliselect.NewOption("Red", Color2{Name: "Red", ID: 1, Sub: "Sub Red"}))
	colors2 = append(colors2, gocliselect.NewOption("Blue", Color2{Name: "Blue", ID: 2, Sub: "Sub Blue"}))
	colors2 = append(colors2, gocliselect.NewOption("Green", Color2{Name: "Green", ID: 3, Sub: "Sub Green"}))
	colors2 = append(colors2, gocliselect.NewOption("Yellow", Color2{Name: "Yellow", ID: 4, Sub: "Sub Yellow"}))

	menu2 := gocliselect.NewSelect[Color2]().
		Title("Choose a color with sub:").
		Options(colors2...).
		Value(&selectedColor2)

	menu2.OptionSelectColor = pp.Green

	if err := menu2.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option with sub: %v\n\n", selectedColor2)

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
