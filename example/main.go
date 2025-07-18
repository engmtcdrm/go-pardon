package main

import (
	"fmt"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon"
	"github.com/engmtcdrm/go-pardon/example/examples"
)

type Color2 struct {
	Name string
	ID   int
	Sub  string
}

func main() {
	pardon.SetDefaultIconFunc(func(s string) string { return fmt.Sprintf("%s%s%s", ansi.Green, s, ansi.Reset) })
	pardon.SetDefaultSelectFunc(func(s string) string { return fmt.Sprintf("%s%s%s", ansi.Yellow, s, ansi.Reset) })
	pardon.SetDefaultCursorFunc(func(s string) string { return fmt.Sprintf("%s%s%s", ansi.Blue, s, ansi.Reset) })
	pardon.SetDefaultAnswerFunc(func(s string) string { return fmt.Sprintf("%s%s%s", ansi.Cyan, s, ansi.Reset) })

	funcMap := map[string]func(){}
	names := make([]pardon.Option[string], len(examples.AllExamples))

	for i, ex := range examples.AllExamples {
		funcMap[ex.Name] = ex.Fn
		names[i] = pardon.NewOption(fmt.Sprintf("%d. %s", i+1, ex.Name), ex.Name)
	}

	var selectedName string

	selectPrompt := pardon.NewSelect[string]().
		Title("Select an example:").
		Icon("").
		Options(names...).
		Value(&selectedName).
		AnswerFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Yellow, s, ansi.Reset)
		})

	if err := selectPrompt.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println()

	if fn, ok := funcMap[selectedName]; ok {
		fn()
	} else {
		fmt.Println("No function found for selection.")
	}
}
