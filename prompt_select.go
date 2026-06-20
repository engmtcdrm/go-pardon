package pardon

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/grapheme"
	"github.com/mattn/go-runewidth"
)

// Select represents a multiple-choice selection prompt.
type Select[T comparable] struct {
	PromptBase[T, *Select[T]]

	cursor          eval[string]
	selectEval      eval[string] // cannot use select because it is a reserved keyword
	options         []Option[T]
	selectFn        func(string) string
	cursorPos       int
	cursorCharWidth int
	scrollOffset    int
}

// NewSelect creates a new Select prompt instance.
func NewSelect[T comparable](value *T) *Select[T] {
	s := &Select[T]{
		PromptBase: NewPromptBase[T, *Select[T]](value),
		cursor:     eval[string]{val: "> ", defaultFn: defaultFuncs.cursorFn},
		selectEval: eval[string]{val: "", defaultFn: defaultFuncs.selectFn},
		options:    make([]Option[T], 0),
	}
	s.PromptBase.Self = s
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

func (s *Select[T]) ask() error {
	// Need to save cursor location so we can easily redraw the prompt after
	// user input without needing to recalculate cursor movements.
	fmt.Fprint(s.Out, saveCursor)

	fmt.Fprint(s.Out, ansi.HideCursor)
	defer func() {
		fmt.Fprint(s.Out, ansi.ShowCursor)
	}()

	s.buildAndSetPrompt()
	s.render()

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

func (s *Select[T]) handleCtrlC(_ grapheme.Cluster) (done bool, err error) {
	var builder strings.Builder
	builder.WriteString(restoreCursor + ansi.ClearFromCursorToEndScreen)
	builder.WriteString(s.Prompt.String())
	fmt.Fprint(s.Out, builder.String())

	return true, ErrUserAborted
}

// handleEnter stores the selected value and prints the selected option as the
// final answer.
func (s *Select[T]) handleEnter(_ grapheme.Cluster) (done bool, err error) {
	*s.value = s.options[s.cursorPos].Value
	s.answer.val = s.options[s.cursorPos].Key
	s.printFinalPromptLine()
	return true, nil
}

// handleEscape processes escape sequences for navigating the options list.
func (s *Select[T]) handleEscape(seq grapheme.Cluster) (done bool, err error) {
	switch {
	case grapheme.Equal(seq, grapheme.UpArrow):
		s.cursorPos = (s.cursorPos + len(s.options) - 1) % len(s.options)
		s.render()
	case grapheme.Equal(seq, grapheme.DownArrow):
		s.cursorPos = (s.cursorPos + 1) % len(s.options)
		s.render()
	}
	return false, nil
}

// printFinalPromptLine handles printing the final prompt line.
func (s *Select[T]) printFinalPromptLine() {
	var builder strings.Builder
	builder.WriteString(restoreCursor + ansi.ClearFromCursorToEndScreen)
	builder.WriteString(s.Prompt.String())
	builder.WriteString(s.answer.Get())
	builder.WriteString("\n")
	fmt.Fprint(s.Out, builder.String())
}

func (s *Select[T]) processInput(input []byte) (done bool, err error) {
	if len(input) == 0 {
		return false, nil
	}

	s.PendingInputBytes = append(s.PendingInputBytes, input...)

	if needMoreInput := s.ConvertBytesToGraphemeSet(); needMoreInput {
		return false, nil
	}

	for len(s.PendingInputClusterSet) > 0 {
		r := s.PendingInputClusterSet[0]

		switch {
		case grapheme.Equal(r, grapheme.CtrlC):
			return s.handleCtrlC(r)
		case grapheme.Equal(r, grapheme.Enter), grapheme.Equal(r, grapheme.Newline):
			return s.handleEnter(r)
		case grapheme.Equal(r, grapheme.Escape):
			doContinue, err := s.ProcessEscapeSequence(r, s.handleEscape)
			if !doContinue {
				return false, err
			}
			continue
		}

		s.PendingInputClusterSet = s.PendingInputClusterSet[1:]
	}

	return false, nil
}

func (s *Select[T]) redraw(selectSize, termHeight int) {
	var output strings.Builder
	output.WriteString(restoreCursor + ansi.ClearFromCursorToEndScreen)
	output.WriteString(s.Prompt.String())
	output.WriteString("\n")

	selectCursor := s.cursor.Get()
	minSize := min(s.scrollOffset+termHeight, selectSize)

	for i := s.scrollOffset; i < minSize; i++ {
		selectedOption := s.options[i]
		cwidth := runewidth.StringWidth(ansi.Strip(selectCursor))
		cursor := strings.Repeat(" ", cwidth)

		if i != s.cursorPos {
			output.WriteString(cursor)
			output.WriteString(selectedOption.Key)
			if i != minSize-1 {
				output.WriteString("\n")
			}
			continue
		}

		s.selectEval.val = selectedOption.Key
		output.WriteString(selectCursor)
		output.WriteString(s.selectEval.Get())
		if i != minSize-1 {
			output.WriteString("\n")
		}
	}

	fmt.Fprint(s.Out, output.String())
}

// updateScrollOffset updates the scroll offset to ensure the selected option is
// visible within the terminal height.
func (s *Select[T]) updateScrollOffset(termHeight int) {
	if s.cursorPos < s.scrollOffset {
		s.scrollOffset = s.cursorPos
	} else if s.cursorPos >= s.scrollOffset+termHeight {
		s.scrollOffset = s.cursorPos - termHeight + 1
	}
}

// render displays the list of available options to the user.
func (s *Select[T]) render() {
	_, termHeight := s.GetTerminalSize()
	termHeight = termHeight - 2 // -2 for breathing room for prompt line
	selectSize := len(s.options)

	s.updateScrollOffset(termHeight)
	s.redraw(selectSize, termHeight)
}
