package tui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/keys"
	"golang.org/x/term"
)

var (
	ErrUserAborted = errors.New("user aborted")
)

// InputPrompt handles text-based input prompts (Question, Password, etc.)
type InputPrompt[T any] struct {
	appendInputFn  func(T, byte) T
	displayInputFn func(T) string
	removeLastFn   func(T) T
	answerFn       func(string) string
	validateFn     func(T) error
}

// NewStringPrompt creates a prompt for string input (Question type)
func NewStringPrompt() *InputPrompt[string] {
	return &InputPrompt[string]{
		appendInputFn:  func(s string, b byte) string { return s + string(b) },
		displayInputFn: func(s string) string { return s },
		removeLastFn: func(s string) string {
			if len(s) > 0 {
				return s[:len(s)-1]
			}
			return s
		},
		answerFn:   func(s string) string { return s },
		validateFn: func(s string) error { return nil },
	}
}

// NewPasswordPrompt creates a prompt for password input ([]byte type)
func NewPasswordPrompt() *InputPrompt[[]byte] {
	return &InputPrompt[[]byte]{
		appendInputFn: func(b []byte, c byte) []byte { return append(b, c) },
		displayInputFn: func(b []byte) string {
			if len(b) == 0 {
				return ""
			}
			return strings.Repeat("", len(b))
		},
		removeLastFn: func(b []byte) []byte {
			if len(b) > 0 {
				return b[:len(b)-1]
			}
			return b
		},
		answerFn:   func(s string) string { return s },
		validateFn: func(b []byte) error { return nil },
	}
}

func (p *InputPrompt[T]) Validate(fn func(T) error) *InputPrompt[T] {
	if fn != nil {
		p.validateFn = fn
	}
	return p
}

func (p *InputPrompt[T]) DisplayInput(fn func(T) string) *InputPrompt[T] {
	if fn != nil {
		p.displayInputFn = fn
	}
	return p
}

func (p *InputPrompt[T]) AnswerFunc(fn func(string) string) *InputPrompt[T] {
	if fn != nil {
		p.answerFn = fn
	}
	return p
}

func (p *InputPrompt[T]) Display(prompt string, value *T) error {
	input := *value
	var lastError string
	showError := false

	redraw := func() {
		// Clear the current line
		fmt.Printf("\r%s", ansi.ClearLine)
		fmt.Printf("%s", prompt)

		// Clear the next line
		fmt.Printf("\n%s", ansi.ClearLine)

		// If there's an error, display it
		if showError && lastError != "" {
			fmt.Printf("%s* %v%s", ansi.Red, lastError, ansi.Reset)
		}

		// Display the current input
		fmt.Printf("%s\r%s", ansi.CursorUp(1), ansi.ClearLine)
		fmt.Print(prompt)
		fmt.Print(p.displayInputFn(input))
	}

	redraw()

	for {
		keyCode := GetInput()
		switch keyCode {
		case keys.KeyEnter, keys.KeyCarriageReturn:
			if err := p.validateFn(input); err != nil {
				lastError = err.Error()
				showError = true
				redraw()
				continue
			}
			*value = input
			fmt.Printf("\r%s%s\n%s", prompt, p.answerFn(p.displayInputFn(input)), ansi.ClearLine)
			return nil
		case keys.KeyCtrlC:
			fmt.Printf("\n%s", ansi.ClearLine)
			return ErrUserAborted
		case keys.KeyBackspace:
			input = p.removeLastFn(input)
			showError = false
		case keys.KeyUp, keys.KeyDown:
			// Only treat as navigation keys if they came from escape sequences
			if lastInputWasEscapeSequence {
				showError = false
			} else {
				// Treat as regular input (A=65, B=66)
				input = p.appendInputFn(input, keyCode)
				showError = false
			}
		default:
			input = p.appendInputFn(input, keyCode)
			showError = false
		}
		redraw()
	}
}

var (
	// navigationKeys defines a map of specific byte keycodes related to
	// navigation functionality, such as up or down actions.
	navigationKeys = map[byte]bool{
		keys.KeyUp:   true,
		keys.KeyDown: true,
	}

	// Track whether the last input was from an escape sequence
	lastInputWasEscapeSequence = false
)

// GetInput will read raw input from the terminal
// It returns the raw ASCII value inputted
func GetInput() byte {
	// Use stdin file descriptor for cross-platform compatibility
	fd := int(os.Stdin.Fd())

	// Save the original terminal state
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// Fallback if raw mode fails - read normally
		readBytes := make([]byte, 1)
		_, readErr := os.Stdin.Read(readBytes)
		if readErr != nil {
			return 0
		}
		return readBytes[0]
	}
	defer term.Restore(fd, oldState)

	// Read input
	readBytes := make([]byte, 3)
	read, err := os.Stdin.Read(readBytes)
	if err != nil {
		// Handle read error, it might be due to signal interruption
		return 0
	}

	// Arrow keys are prefixed with the ANSI escape code which take up the first two bytes.
	// The third byte is the key specific value we are looking for.
	// For example the left arrow key is '<esc>[A' while the right is '<esc>[C'
	// See: https://en.wikipedia.org/wiki/ANSI_escape_code
	if read == 3 && readBytes[0] == keys.KeyEscape && readBytes[1] == keys.KeyLeftBracket {
		// This is a proper ANSI escape sequence (ESC[X)
		if _, ok := navigationKeys[readBytes[2]]; ok {
			lastInputWasEscapeSequence = true
			return readBytes[2]
		}
		// If it's an escape sequence but not a navigation key, ignore it
		lastInputWasEscapeSequence = false
		return 0
	}

	// For any other input (1, 2, or 3 bytes that aren't escape sequences),
	// return the first byte which contains the actual character
	lastInputWasEscapeSequence = false
	if read > 0 {
		return readBytes[0]
	}

	return 0
}
