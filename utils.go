package pardon

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
)

// equal reports whether a and b are the same length and contain the same runes.
// A nil argument is equivalent to an empty slice.
func equal(a []rune, b []rune) bool {
	return string(a) == string(b)
}

// equalFold reports whether a and b, interpreted as UTF-8 strings,
// are equal under simple Unicode case-folding, which is a more general
// form of case-insensitivity.
func equalFold(a []rune, b ...rune) bool {
	return strings.EqualFold(string(a), string(b))
}

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
