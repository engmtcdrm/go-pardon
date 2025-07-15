package main

import (
	"fmt"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/gocliselect"
	"github.com/engmtcdrm/gocliselect/example/examples"
)

type Color2 struct {
	Name string
	ID   int
	Sub  string
}

func main() {
	funcMap := map[string]func(){}
	names := make([]gocliselect.Option[string], len(examples.AllExamples))

	for i, ex := range examples.AllExamples {
		funcMap[ex.Name] = ex.Fn
		names[i] = gocliselect.NewOption(ex.Name, ex.Name)
	}

	var selectedName string

	selectPrompt := gocliselect.NewSelect[string]().
		Title("Select an example:").
		Options(names...).
		Value(&selectedName).
		SelectFunc(func(s string) string {
			return fmt.Sprintf("%s%s%s", ansi.Green, s, ansi.Reset)
		})

	if err := selectPrompt.Ask(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if fn, ok := funcMap[selectedName]; ok {
		fn()
	} else {
		fmt.Println("No function found for selection.")
	}
}
