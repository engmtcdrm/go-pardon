package main

import (
	"fmt"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-eggy"

	pp "github.com/engmtcdrm/go-prettyprint"

	"github.com/engmtcdrm/go-pardon"
	"github.com/engmtcdrm/go-pardon/example/internal"
)

func main() {
	pardon.SetDefaultIconFunc(func(s string) string { return fmt.Sprintf("%s%s%s", ansi.Green, s, ansi.Reset) })
	pardon.SetDefaultSelectFunc(func(s string) string { return fmt.Sprintf("%s%s%s", ansi.Yellow, s, ansi.Reset) })
	pardon.SetDefaultCursorFunc(func(s string) string { return fmt.Sprintf("%s%s%s", ansi.Blue, s, ansi.Reset) })
	pardon.SetDefaultAnswerFunc(func(s string) string { return fmt.Sprintf("%s%s%s", ansi.Cyan, s, ansi.Reset) })

	ex := eggy.NewExamplePrompt(internal.AllExamples).
		Title(pp.Yellow("Examples of Pardon"))
	ex.Show()
}
