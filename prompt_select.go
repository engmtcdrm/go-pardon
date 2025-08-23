package pardon

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/keys"
	"github.com/engmtcdrm/go-pardon/tui"
)

// Select represents a multiple-choice selection prompt.
type Select[T comparable] struct {
	icon         eval[string]
	title        eval[string]
	cursor       eval[string]
	cursorPos    int
	scrollOffset int
	options      []Option[T]
	answerFn     func(string) string
	selectFn     func(string) string
	value        *T
}

// NewSelect creates a new Select prompt instance.
func NewSelect[T comparable]() *Select[T] {
	return &Select[T]{
		icon:    eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title:   eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		cursor:  eval[string]{val: "> ", defaultFn: defaultFuncs.cursorFn},
		options: make([]Option[T], 0),
	}
}

// Title sets the prompt title text that will be displayed to the user.
func (sel *Select[T]) Title(title string) *Select[T] {
	sel.title.val = title
	sel.title.fn = nil
	return sel
}

// TitleFunc sets a function to dynamically format the prompt title.
func (sel *Select[T]) TitleFunc(fn func(string) string) *Select[T] {
	sel.title.fn = fn
	return sel
}

// Cursor sets the cursor symbol displayed next to the selected option.
func (sel *Select[T]) Cursor(cursor string) *Select[T] {
	sel.cursor.val = cursor
	sel.cursor.fn = nil
	return sel
}

// CursorFunc sets a function to dynamically format the cursor symbol.
func (sel *Select[T]) CursorFunc(fn func(string) string) *Select[T] {
	sel.cursor.fn = fn
	return sel
}

// Options sets the list of available options for selection.
func (sel *Select[T]) Options(options ...Option[T]) *Select[T] {
	if len(options) == 0 {
		return sel
	}

	sel.options = options
	return sel
}

// Value sets the pointer where the selected option's value will be stored.
func (sel *Select[T]) Value(value *T) *Select[T] {
	sel.value = value
	return sel
}

// Icon sets the icon displayed before the prompt title.
func (sel *Select[T]) Icon(icon string) *Select[T] {
	sel.icon.val = icon
	sel.icon.fn = nil
	return sel
}

// IconFunc sets a function to dynamically format the prompt icon.
func (sel *Select[T]) IconFunc(fn func(string) string) *Select[T] {
	sel.icon.fn = fn
	return sel
}

// AnswerFunc sets a function to format the final answer display.
func (sel *Select[T]) AnswerFunc(fn func(string) string) *Select[T] {
	sel.answerFn = fn
	return sel
}

// SelectFunc sets a function to format option text during selection.
func (sel *Select[T]) SelectFunc(fn func(string) string) *Select[T] {
	sel.selectFn = fn
	return sel
}

// setAnswerFunc configures the answer transformation priority:
// prompt-specific, global default, or the string itself.
func (sel *Select[T]) getSelectFunc(s string) string {
	if sel.selectFn != nil {
		return sel.selectFn(s)
	}

	if defaultFuncs.selectFn != nil {
		return defaultFuncs.selectFn(s)
	}

	return s
}

// getAnswerFunc returns the formatted text for the final answer display.
func (sel *Select[T]) getAnswerFunc(answer string) string {
	if sel.answerFn != nil {
		return sel.answerFn(answer)
	}

	if defaultFuncs.answerFn != nil {
		return defaultFuncs.answerFn(answer)
	}

	return answer
}

// Ask displays the select prompt and waits for user selection.
func (sel *Select[T]) Ask() error {
	if sel.title.val == "" && sel.title.fn == nil {
		return ErrNoTitle
	}

	if sel.value == nil {
		return ErrNoValue
	}

	if len(sel.options) == 0 {
		return ErrNoSelectOptions
	}

	defer func() {
		fmt.Print(ansi.ShowCursor)
	}()

	// Print the question
	fmt.Printf("%s%s\n", sel.icon.Get(), sel.title.Get())

	sel.renderOptions(false)
	fmt.Print(ansi.HideCursor)

	for {
		keyCode := tui.GetInput()

		switch keyCode {
		case keys.KeyCtrlC:
			return ErrUserAborted
		case keys.KeyEnter, keys.KeyCarriageReturn:
			*sel.value = sel.options[sel.cursorPos].Value
			visibleOptions := tui.Min(len(sel.options), tui.GetTerminalHeight()-3)
			tui.RenderClearAndReposition(visibleOptions+1, sel.icon.Get(), sel.title.Get(), sel.getAnswerFunc(sel.options[sel.cursorPos].Key))
			return nil
		case keys.KeyUp:
			sel.cursorPos = (sel.cursorPos + len(sel.options) - 1) % len(sel.options)
			sel.renderOptions(true)
		case keys.KeyDown:
			sel.cursorPos = (sel.cursorPos + 1) % len(sel.options)
			sel.renderOptions(true)
		}
	}
}

// renderOptions displays the list of available options to the user.
func (sel *Select[T]) renderOptions(redraw bool) {
	termHeight := tui.GetTerminalHeight()
	termHeight = termHeight - 3 // Space for prompt and cursor movement
	selectSize := len(sel.options)

	// Ensure scroll offset follows cursor movement
	if sel.cursorPos < sel.scrollOffset {
		sel.scrollOffset = sel.cursorPos
	} else if sel.cursorPos >= sel.scrollOffset+termHeight {
		sel.scrollOffset = sel.cursorPos - termHeight + 1
	}

	selectCursor := sel.cursor.Get()
	visibleLines := tui.Min(selectSize, termHeight)

	if redraw {
		// For terminal optimization: build entire output first, then write atomically
		var output strings.Builder

		// Move cursor up to start position
		output.WriteString(ansi.CursorUp(visibleLines))

		// Build all lines in memory first
		for i := sel.scrollOffset; i < tui.Min(sel.scrollOffset+termHeight, selectSize); i++ {
			selectedOption := sel.options[i]
			cursor := strings.Repeat(" ", utf8.RuneCountInString(ansi.StripCodes(selectCursor)))

			// Clear line and build content
			output.WriteString("\r")
			output.WriteString(ansi.ClearLine)

			if i == sel.cursorPos {
				cursor = sel.getSelectFunc(selectCursor)
				output.WriteString(cursor)
				output.WriteString(sel.getSelectFunc(selectedOption.Key))
			} else {
				output.WriteString(cursor)
				output.WriteString(selectedOption.Key)
			}
			output.WriteString("\n")
		}

		// Write everything at once to minimize flicker
		fmt.Print(output.String())
	} else {
		// Initial render without redraw
		for i := sel.scrollOffset; i < tui.Min(sel.scrollOffset+termHeight, selectSize); i++ {
			selectedOption := sel.options[i]
			cursor := strings.Repeat(" ", utf8.RuneCountInString(ansi.StripCodes(selectCursor)))

			if i == sel.cursorPos {
				cursor = sel.getSelectFunc(selectCursor)
				fmt.Printf("%s%s\n", cursor, sel.getSelectFunc(selectedOption.Key))
			} else {
				fmt.Printf("%s%s\n", cursor, selectedOption.Key)
			}
		}
	}
}
