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
	// ErrUserAborted is returned when the user cancels a prompt operation.
	ErrUserAborted = errors.New("user aborted")
)

var (
	// navigationKeys defines a map of byte keycodes for navigation actions.
	// These keys are used for cursor movement and selection in interactive prompts.
	navigationKeys = map[byte]bool{
		keys.KeyUp:    true,
		keys.KeyDown:  true,
		keys.KeyLeft:  true,
		keys.KeyRight: true,
	}

	// lastInputWasEscSeq tracks whether the previous input was part of an escape sequence.
	// This helps with proper handling of multi-byte terminal input sequences.
	lastInputWasEscSeq = false

	// inputBuffer provides buffering for handling paste operations and multi-byte input.
	inputBuffer []byte
)

// InputPrompt provides a generic framework for text-based input prompts.
type InputPrompt[T any] struct {
	appendInputFn  func(T, byte) T
	displayInputFn func(T) string
	removeLastFn   func(T) T
	answerFn       func(string) string
	validateFn     func(T) error
}

// NewStringPrompt creates an InputPrompt for plaintext string input.
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

// NewPasswordPrompt creates an InputPrompt for secure password input with masking.
func NewPasswordPrompt() *InputPrompt[[]byte] {
	return &InputPrompt[[]byte]{
		appendInputFn:  func(b []byte, c byte) []byte { return append(b, c) },
		displayInputFn: func(b []byte) string { return "" }, // Mask all input
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

// Validate sets a validation function for the input prompt.
func (p *InputPrompt[T]) Validate(fn func(T) error) *InputPrompt[T] {
	if fn != nil {
		p.validateFn = fn
	}
	return p
}

// DisplayInput sets a function to customize how input is displayed during typing.
func (p *InputPrompt[T]) DisplayInput(fn func(T) string) *InputPrompt[T] {
	if fn != nil {
		p.displayInputFn = fn
	}
	return p
}

// AnswerFunc sets a function to transform the final answer before returning.
func (p *InputPrompt[T]) AnswerFunc(fn func(string) string) *InputPrompt[T] {
	if fn != nil {
		p.answerFn = fn
	}
	return p
}

// getActualInputLength gets the actual length of input for display calculations
// This is used for incremental updates where display length != actual length
func getActualInputLength[T any](input T) int {
	switch v := any(input).(type) {
	case []byte:
		return len(v)
	case string:
		return len(v)
	default:
		return 0
	}
}

// Display displays the input prompt and handles user input
func (p *InputPrompt[T]) Display(prompt string, value *T) error {
	input := *value
	var lastError string
	showError := false
	var lastInputLen int
	var wasShowingError bool

	fmt.Print(prompt)

	redraw := func() {
		currentInput := p.displayInputFn(input)
		currentInputLen := len(currentInput)
		actualInputLen := getActualInputLength(input)

		if showError && lastError != "" {
			displayErr := fmt.Sprintf(
				"\r%s%s%s\n\r%s%s* %s%s%s\r%s%s",
				ansi.ClearLine,
				prompt,
				currentInput,
				ansi.ClearLine,
				ansi.Red,
				lastError,
				ansi.Reset,
				ansi.CursorUp(1),
				prompt,
				currentInput,
			)
			fmt.Print(displayErr)
			wasShowingError = true
		} else {
			// If we were showing an error and now we're not, clear the error line
			if wasShowingError {
				clearErr := fmt.Sprintf(
					"\r%s%s%s\n\r%s%s\r%s%s",
					ansi.ClearLine,
					prompt,
					currentInput,
					ansi.ClearLine,
					ansi.CursorUp(1),
					prompt,
					currentInput,
				)
				fmt.Print(clearErr)
				wasShowingError = false
			} else {
				// No error - use incremental update approach
				if currentInputLen == 0 && actualInputLen > 0 {
					// This is a password field - no visual updates needed
					// The cursor stays in place for security
				} else if currentInputLen < lastInputLen {
					// Input got shorter (backspace)
					extraChars := lastInputLen - currentInputLen
					backspaceSeq := fmt.Sprintf("%s%s%s",
						strings.Repeat("\b", extraChars),
						strings.Repeat(" ", extraChars),
						strings.Repeat("\b", extraChars))
					fmt.Print(backspaceSeq)
				} else if currentInputLen > lastInputLen {
					// Input got longer - just append the new characters
					newChars := currentInput[lastInputLen:]
					fmt.Print(newChars)
				}
				// For same length or other cases, we don't need to do anything
				// since the cursor is already in the right position
			}
		}

		lastInputLen = currentInputLen
	}

	// Initial display of any existing input
	if len(p.displayInputFn(input)) > 0 {
		fmt.Print(p.displayInputFn(input))
		lastInputLen = len(p.displayInputFn(input))
	}

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
			// Clear any error lines and display final output
			var finalOutput string

			if showError {
				finalOutput = fmt.Sprintf(
					"\r%s\n%s\r%s%s\n",
					ansi.ClearLine,
					ansi.ClearLine,
					prompt,
					p.answerFn(p.displayInputFn(input)),
				)
			} else {
				finalOutput = fmt.Sprintf(
					"\r%s%s%s\n",
					ansi.ClearLine,
					prompt,
					p.answerFn(p.displayInputFn(input)),
				)
			}
			fmt.Print(finalOutput)
			return nil
		case keys.KeyCtrlC:
			// Clear any error lines
			var finalOutput string

			if showError {
				finalOutput = fmt.Sprintf("\r%s\n%s\r", ansi.ClearLine, ansi.ClearLine)
			} else {
				finalOutput = fmt.Sprintf("\r%s", ansi.ClearLine)
			}
			fmt.Print(finalOutput)
			return ErrUserAborted
		case keys.KeyBackspace:
			input = p.removeLastFn(input)
			showError = false
			redraw()
		case keys.KeyUp, keys.KeyDown, keys.KeyLeft, keys.KeyRight:
			// Only treat as navigation keys if they came from escape sequences
			if lastInputWasEscSeq {
				showError = false
				// Don't redraw for navigation keys to prevent flicker
				continue
			} else {
				// Treat as regular input (A=65, B=66, C=67, D=68)
				input = p.appendInputFn(input, keyCode)
				showError = false
				redraw()
			}
		default:
			// Filter out control characters (0-31), allow all others
			if keyCode >= 32 {
				// Printable ASCII (32-126) and extended characters (128+)
				input = p.appendInputFn(input, keyCode)
				showError = false
				redraw()
			}
		}
	}
}

// GetInput reads raw keyboard input from the terminal.
// Handles buffered input, raw mode, and ANSI escape sequences.
func GetInput() byte {
	// If we have buffered input from a paste operation, return it first
	if len(inputBuffer) > 0 {
		result := inputBuffer[0]
		inputBuffer = inputBuffer[1:]
		lastInputWasEscSeq = false
		return result
	}

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

	// Read input - use a larger buffer to handle paste operations
	readBytes := make([]byte, 4096) // Increased to 4KB to handle larger pastes
	read, err := os.Stdin.Read(readBytes)
	if err != nil {
		// Handle read error, it might be due to signal interruption
		return 0
	}

	// If we read more than 3 bytes, it's likely a paste operation
	if read > 3 {
		// Buffer all characters except the first one
		inputBuffer = append(inputBuffer, readBytes[1:read]...)
		lastInputWasEscSeq = false
		return readBytes[0]
	}

	// Handle escape sequences (arrow keys)
	if read == 3 && readBytes[0] == keys.KeyEscape && readBytes[1] == keys.KeyLeftBracket {
		// This is a proper ANSI escape sequence (ESC[X)
		if _, ok := navigationKeys[readBytes[2]]; ok {
			lastInputWasEscSeq = true
			return readBytes[2]
		}
		// If it's an escape sequence but not a navigation key, ignore it
		lastInputWasEscSeq = false
		return 0
	}

	// For any other input (1, 2, or 3 bytes that aren't escape sequences),
	// return the first byte which contains the actual character
	lastInputWasEscSeq = false
	if read > 0 {
		return readBytes[0]
	}

	return 0
}
