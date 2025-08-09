package examples

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
)

func SelectKitchensink() {
	var selectedColor int
	colors := []pardon.Option[int]{}
	colors = append(colors, pardon.Option[int]{Key: "Red", Value: 1})
	colors = append(colors, pardon.Option[int]{Key: "Blue", Value: 2})
	colors = append(colors, pardon.Option[int]{Key: "Green", Value: 3})
	colors = append(colors, pardon.Option[int]{Key: "Yellow", Value: 4})

	selectPrompt := pardon.NewSelect[int]().
		Title("Choose a color:").
		TitleFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Green, s, ansi.Reset)
		}).
		Icon("??? ").
		IconFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Blue, s, ansi.Reset)
		}).
		AnswerFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.CyanBg, s, ansi.Reset)
		}).
		Cursor("» ").
		CursorFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s%s", ansi.Magenta, ansi.YellowBg, s, ansi.Reset)
		}).
		SelectFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s%s", ansi.RedBg, ansi.Cyan, s, ansi.Reset)
		}).
		Options(colors...).
		Value(&selectedColor)

	if err := selectPrompt.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Selected option: %v\n", selectedColor)
	os.Exit(0)
}
