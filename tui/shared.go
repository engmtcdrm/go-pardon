package tui

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-ansi"
	"golang.org/x/term"
)

// Shared rendering utilities

// Min returns the smaller of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RenderFinalAnswer renders the final answer after a prompt is completed
func RenderFinalAnswer(icon, title, answer string) {
	fmt.Printf("%s%s %s\n", icon, title, answer)
}

// RenderClearLines clears a specific number of lines from the current cursor position
func RenderClearLines(numLines int) {
	for i := 0; i < numLines; i++ {
		fmt.Printf("%s\r\n", ansi.ClearLine)
	}
}

// RenderClearAndReposition clears lines and repositions cursor for final output
func RenderClearAndReposition(linesToErase int, icon, title, answer string) {
	// Move cursor up to question line
	fmt.Print(ansi.CursorUp(linesToErase))

	// Clear lines
	RenderClearLines(linesToErase)

	// Move cursor back up to question line
	fmt.Print(ansi.CursorUp(linesToErase))

	// Print final answer
	RenderFinalAnswer(icon, title, answer)
	fmt.Print(ansi.ShowCursor)
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
