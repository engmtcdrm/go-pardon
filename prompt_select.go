package pardon

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/mattn/go-runewidth"

	"github.com/engmtcdrm/go-pardon/internal/keys"
)

var (
	// navigationKeys defines a map of byte keycodes for navigation actions.
	// These keys are used for cursor movement and selection in interactive prompts.
	navigationKeys = map[byte]bool{
		keys.Up:    true,
		keys.Down:  true,
		keys.Left:  true,
		keys.Right: true,
	}

	// lastInputWasEscSeq tracks whether the previous input was part of an escape sequence.
	// This helps with proper handling of multi-byte terminal input sequences.
	lastInputWasEscSeq = false
)

// Select represents a multiple-choice selection prompt.
type Select[T comparable] struct {
	terminal     *Terminal
	icon         eval[string]
	title        eval[string]
	cursor       eval[string]
	answer       eval[string]
	selectEval   eval[string] // cannot use select because it is a reserved keyword
	options      []Option[T]
	selectFn     func(string) string
	cursorPos    int
	scrollOffset int
	value        *T
}

// NewSelect creates a new Select prompt instance.
func NewSelect[T comparable](value *T) *Select[T] {
	return &Select[T]{
		terminal:   NewTerminal(),
		icon:       eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title:      eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		cursor:     eval[string]{val: "> ", defaultFn: defaultFuncs.cursorFn},
		answer:     eval[string]{val: "", defaultFn: defaultFuncs.answerFn},
		selectEval: eval[string]{val: "", defaultFn: defaultFuncs.selectFn},
		options:    make([]Option[T], 0),
		value:      value,
	}
}

// AnswerFunc sets a function to format the final answer display.
func (s *Select[T]) AnswerFunc(fn func(string) string) *Select[T] {
	s.answer.fn = fn
	return s
}

// Ask displays the select prompt and waits for user selection.
func (s *Select[T]) Ask() error {
	if s.title.val == "" && s.title.fn == nil {
		return ErrNoTitle
	}

	if s.value == nil {
		return ErrNoValue
	}

	if len(s.options) == 0 {
		return ErrNoSelectOptions
	}

	defer func() {
		fmt.Fprint(s.terminal.Out, ansi.ShowCursor)
	}()

	// Print the question
	fmt.Fprintf(s.terminal.Out, "%s%s\n", s.icon.Get(), s.title.Get())

	s.renderOptions(false)
	fmt.Fprint(s.terminal.Out, ansi.HideCursor)

	for {
		keyCode, err := s.terminal.GetInput()
		if err != nil {
			return err
		}

		switch keyCode {
		case keys.CtrlC:
			return ErrUserAborted
		case keys.Enter, keys.NewLine:
			*s.value = s.options[s.cursorPos].Value
			s.answer.val = s.options[s.cursorPos].Key
			visibleOptions := min(len(s.options), s.terminal.GetTerminalHeight()-3)
			renderClearAndReposition(visibleOptions+1, s.icon.Get(), s.title.Get(), s.answer.Get())
			return nil
		case keys.Up:
			s.cursorPos = (s.cursorPos + len(s.options) - 1) % len(s.options)
			s.renderOptions(true)
		case keys.Down:
			s.cursorPos = (s.cursorPos + 1) % len(s.options)
			s.renderOptions(true)
		}
	}
}

// Cursor sets the cursor symbol displayed next to the selected option.
func (s *Select[T]) Cursor(cursor string) *Select[T] {
	s.cursor.val = cursor
	s.cursor.fn = nil
	return s
}

// CursorFunc sets a function to dynamically format the cursor symbol.
func (s *Select[T]) CursorFunc(fn func(string) string) *Select[T] {
	s.cursor.fn = fn
	return s
}

// Icon sets the icon displayed before the prompt title.
func (s *Select[T]) Icon(icon string) *Select[T] {
	s.icon.val = icon
	s.icon.fn = nil
	return s
}

// IconFunc sets a function to dynamically format the prompt icon.
func (s *Select[T]) IconFunc(fn func(string) string) *Select[T] {
	s.icon.fn = fn
	return s
}

// Options sets the list of available options for selection.
func (s *Select[T]) Options(options ...Option[T]) *Select[T] {
	if len(options) == 0 {
		return s
	}

	s.options = options
	return s
}

// SelectFunc sets a function to format option text during selection.
func (s *Select[T]) SelectFunc(fn func(string) string) *Select[T] {
	s.selectEval.fn = fn
	return s
}

// Title sets the prompt title text that will be displayed to the user.
func (s *Select[T]) Title(title string) *Select[T] {
	s.title.val = title
	return s
}

// TitleFunc sets a function to dynamically format the prompt title.
func (s *Select[T]) TitleFunc(fn func(string) string) *Select[T] {
	s.title.fn = fn
	return s
}

// Value sets the pointer where the selected option's value will be stored.
func (s *Select[T]) Value(value *T) *Select[T] {
	s.value = value
	return s
}

func (s *Select[T]) redraw(selectSize, termHeight int) {
	selectCursor := s.cursor.Get()
	visibleLines := min(selectSize, termHeight)

	// For terminal optimization: build entire output first, then write atomically
	var output strings.Builder

	// Move cursor up to start position
	output.WriteString(ansi.CursorUp(visibleLines))

	// Build all lines in memory first
	for i := s.scrollOffset; i < min(s.scrollOffset+termHeight, selectSize); i++ {
		selectedOption := s.options[i]
		cursor := strings.Repeat(" ", runewidth.StringWidth(ansi.Strip(selectCursor)))

		// Clear line and build content
		output.WriteString("\r")
		output.WriteString(ansi.ClearLine)

		if i != s.cursorPos {
			output.WriteString(cursor)
			output.WriteString(selectedOption.Key)
			output.WriteString("\n")
			continue
		}

		s.selectEval.val = selectedOption.Key
		output.WriteString(selectCursor)
		output.WriteString(s.selectEval.Get())
		output.WriteString("\n")
	}

	// Write everything at once to minimize flicker
	fmt.Fprint(s.terminal.Out, output.String())
}

func (s *Select[T]) updateScrollOffset(termHeight int) {
	if s.cursorPos < s.scrollOffset {
		s.scrollOffset = s.cursorPos
	} else if s.cursorPos >= s.scrollOffset+termHeight {
		s.scrollOffset = s.cursorPos - termHeight + 1
	}
}

// renderOptions displays the list of available options to the user.
func (s *Select[T]) renderOptions(redraw bool) {
	termHeight := s.terminal.GetTerminalHeight() - 3
	selectSize := len(s.options)

	s.updateScrollOffset(termHeight)

	if redraw {
		s.redraw(selectSize, termHeight)
	} else {
		selectCursor := s.cursor.Get()

		// Initial render without redraw
		for i := s.scrollOffset; i < min(s.scrollOffset+termHeight, selectSize); i++ {
			selectedOption := s.options[i]

			if i != s.cursorPos {
				cursor := strings.Repeat(" ", runewidth.StringWidth(ansi.Strip(selectCursor)))
				fmt.Fprintf(s.terminal.Out, "%s%s\n", cursor, selectedOption.Key)
				continue
			}

			s.selectEval.val = selectedOption.Key
			fmt.Fprintf(s.terminal.Out, "%s%s\n", selectCursor, s.selectEval.Get())
		}
	}
}
