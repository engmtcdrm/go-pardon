package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"golang.org/x/term"
)

// Min returns the smaller of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ClearCurrentLine clears the current line using standard ANSI codes
func ClearCurrentLine() {
	fmt.Printf("\r%s", ansi.ClearLine)
}

// RenderFinalAnswer renders the final answer after a prompt is completed
func RenderFinalAnswer(icon, title, answer string) {
	fmt.Printf("%s%s %s\n", icon, title, answer)
}

// RenderClearLines clears a specific number of lines from the current cursor position
func RenderClearLines(numLines int) {
	if numLines <= 0 {
		return
	}

	sequence := ansi.ClearLine + "\r\n"
	fmt.Print(strings.Repeat(sequence, numLines))
}

// RenderClearAndReposition clears lines and repositions cursor for final output
func RenderClearAndReposition(linesToErase int, icon, title, answer string) {
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

// GetTerminalHeight gets the terminal height with a fallback default
func GetTerminalHeight() int {
	termHeight := 25 // Default height
	if _, height, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		termHeight = height
	}
	return termHeight
}

// RenderFormattedOutput formats and renders output with ANSI clear
func RenderFormattedOutput(question, result string) string {
	return fmt.Sprintf("%s\r%s %s\n", ansi.ClearToBegin, question, result)
}
