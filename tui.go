package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-ansi"
)

type tuiPrompt[T any] struct {
	Question       string
	Value          *T
	validateFn     func(T) error
	displayInputFn func(T) string
	appendInputFn  func(T, byte) T
	removeLastFn   func(T) T
	convertInputFn func(T) T
}

func newTuiPrompt[T any]() *tuiPrompt[T] {
	return &tuiPrompt[T]{
		validateFn:     func(input T) error { return nil },
		displayInputFn: func(input T) string { return fmt.Sprintf("%v", input) },
		appendInputFn:  func(input T, key byte) T { return input },
		removeLastFn:   func(input T) T { return input },
		convertInputFn: func(input T) T { return input },
	}
}

func (t *tuiPrompt[T]) Validate(fn func(T) error) *tuiPrompt[T] {
	if fn != nil {
		t.validateFn = fn
	}
	return t
}

func (t *tuiPrompt[T]) DisplayInput(fn func(T) string) *tuiPrompt[T] {
	if fn != nil {
		t.displayInputFn = fn
	}
	return t
}

func (t *tuiPrompt[T]) AppendInput(fn func(T, byte) T) *tuiPrompt[T] {
	if fn != nil {
		t.appendInputFn = fn
	}
	return t
}

func (t *tuiPrompt[T]) RemoveLast(fn func(T) T) *tuiPrompt[T] {
	if fn != nil {
		t.removeLastFn = fn
	}
	return t
}

func (t *tuiPrompt[T]) ConvertInput(fn func(T) T) *tuiPrompt[T] {
	if fn != nil {
		t.convertInputFn = fn
	}
	return t
}

func (t *tuiPrompt[T]) Display(prompt string, value *T) error {
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
		fmt.Print(t.displayInputFn(input))
	}

	redraw()

	for {
		keyCode := getInput()
		switch keyCode {
		case keyEnter, keyCarriageReturn:
			val := t.convertInputFn(input)
			if err := t.validateFn(val); err != nil {
				lastError = err.Error()
				showError = true
				redraw()
				continue
			}
			*value = val
			fmt.Printf("\r%s%s\n%s", prompt, t.displayInputFn(input), ansi.ClearLine)
			return nil
		case keyCtrlC:
			fmt.Printf("\n%s", ansi.ClearLine)
			return ErrUserAborted
		case keyBackspace:
			input = t.removeLastFn(input)
			showError = false
		case keyUp, keyDown:
			showError = false
		default:
			input = t.appendInputFn(input, keyCode)
			showError = false
		}
		redraw()
	}
}
