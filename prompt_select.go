package pardon

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"

	"github.com/engmtcdrm/go-pardon/internal/keys"
)

// Select represents a multiple-choice selection prompt.
type Select[T comparable] struct {
	// Out is the output writer for the terminal, typically [os.Stdout].
	Out io.Writer

	// In is the terminal input reader.
	In *TerminalInput

	icon         eval[string]
	title        eval[string]
	cursor       eval[string]
	answer       eval[string]
	selectEval   eval[string] // cannot use select because it is a reserved keyword
	options      []Option[T]
	selectFn     func(string) string
	prompt       string
	cursorPos    int
	scrollOffset int
	value        *T
}

// NewSelect creates a new Select prompt instance.
func NewSelect[T comparable](value *T) *Select[T] {
	return &Select[T]{
		In:         NewTerminalInput(),
		Out:        os.Stdout,
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

	if err := s.ask(); err != nil {
		return err
	}

	return nil
}

func (s *Select[T]) ask() error {
	fmt.Fprint(s.Out, ansi.HideCursor)
	defer func() {
		fmt.Fprint(s.Out, ansi.ShowCursor)
	}()

	s.prompt = fmt.Sprintf("%s%s", s.icon.Get(), s.title.Get())
	fmt.Fprintln(s.Out, s.prompt)

	s.renderOptions(false)

	for {
		input, err := s.In.RawRead()
		if err != nil {
			return err
		}

		if done, err := s.processInput(input); done {
			return err
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

func (s *Select[T]) processInput(input []byte) (done bool, err error) {
	if len(input) == 0 {
		return false, nil
	}

	switch {
	case bytes.Equal(input, keys.CtrlC):
		return true, ErrUserAborted
	case bytes.Equal(input, keys.Enter), bytes.Equal(input, keys.Newline):
		*s.value = s.options[s.cursorPos].Value
		s.answer.val = s.options[s.cursorPos].Key
		visibleOptions := min(len(s.options), s.GetTerminalHeight()-3)
		renderClearAndReposition(visibleOptions+1, s.icon.Get(), s.title.Get(), s.answer.Get())
		return true, nil
	case bytes.Equal(keys.UpArrow, input):
		s.cursorPos = (s.cursorPos + len(s.options) - 1) % len(s.options)
		s.renderOptions(true)
	case bytes.Equal(keys.DownArrow, input):
		s.cursorPos = (s.cursorPos + 1) % len(s.options)
		s.renderOptions(true)
	}

	return false, nil
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
	fmt.Fprint(s.Out, output.String())
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
	termHeight := s.GetTerminalHeight() - 3
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
				fmt.Fprintf(s.Out, "%s%s\n", cursor, selectedOption.Key)
				continue
			}

			s.selectEval.val = selectedOption.Key
			fmt.Fprintf(s.Out, "%s%s\n", selectCursor, s.selectEval.Get())
		}
	}
}

// GetTerminalHeight returns the height of the terminal in rows. If the terminal
// size cannot be determined, it returns a default height of 25 rows.
func (s *Select[T]) GetTerminalHeight() int {
	termHeight := 25 // Default height

	f, ok := s.Out.(*os.File)
	if !ok {
		return termHeight
	}

	if _, height, err := term.GetSize(int(f.Fd())); err == nil {
		termHeight = height
	}

	return termHeight
}
