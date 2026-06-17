package pardon

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
)

// renderClearAndReposition clears lines and renders final answer.
// Minimizes screen flicker by batching terminal operations.
func renderClearAndReposition(linesToErase int, icon, title, answer string) {
	var output strings.Builder

	// Move cursor up to question line
	output.WriteString(ansi.CursorUp(linesToErase))

	sequence := "\r" + ansi.ClearLine

	// Clear all lines in one pass to reduce flickering
	for i := range linesToErase {
		output.WriteString(sequence)
		if i < linesToErase-1 {
			output.WriteString("\n")
		}
	}

	// Move cursor back up to question line
	output.WriteString(ansi.CursorUp(linesToErase - 1))

	// Print final answer
	output.WriteString("\r")
	output.WriteString(icon)
	output.WriteString(title)
	output.WriteString(" ")
	output.WriteString(answer)
	output.WriteString("\n")
	output.WriteString(ansi.ShowCursor)

	// Write everything at once to minimize flicker
	fmt.Print(output.String())
}

func resetLineAbove() string {
	return ansi.CursorUp(1) + ansi.ClearLineReset
}

func validationErrorMessage(err error) string {
	return fmt.Sprintf("%s%s* %v%s", ansi.ClearLineReset, ansi.RedBg, err, ansi.Reset)
}
