package pardon

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/grapheme"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
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

func (s *Select[T]) ask() error {
	fmt.Fprint(s.Out, ansi.HideCursor)
	defer func() {
		fmt.Fprint(s.Out, ansi.ShowCursor)
	}()

	s.buildAndSetPrompt()
	fmt.Fprintln(s.Out, s.Prompt)

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
			return true, ErrUserAborted
		case grapheme.Equal(r, grapheme.Enter), grapheme.Equal(r, grapheme.Newline):
			*s.value = s.options[s.cursorPos].Value
			s.answer.val = s.options[s.cursorPos].Key
			_, termHeight := s.GetTerminalSize()
			visibleOptions := min(len(s.options), termHeight-3)
			renderClearAndReposition(visibleOptions+1, s.icon.Get(), s.title.Get(), s.answer.Get())
			return true, nil
		case grapheme.Equal(r, grapheme.Escape):
			doContinue, err := s.ProcessEscapeSequence(r, s.escapeSequenceHandler)
			if !doContinue {
				return false, err
			}
			continue
		}

		s.PendingInputClusterSet = s.PendingInputClusterSet[1:]
	}

	return false, nil
}

func (s *Select[T]) escapeSequenceHandler(seq grapheme.Cluster) (done bool, err error) {
	switch {
	case grapheme.Equal(seq, grapheme.UpArrow):
		s.cursorPos = (s.cursorPos + len(s.options) - 1) % len(s.options)
		s.renderOptions(true)
	case grapheme.Equal(seq, grapheme.DownArrow):
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
		cwidth := runewidth.StringWidth(ansi.Strip(selectCursor))
		cursor := strings.Repeat(" ", cwidth)
		uniseg.StringWidth(ansi.Strip(selectCursor))

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
	_, termHeight := s.GetTerminalSize()
	termHeight = termHeight - 3
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
				cwidth := runewidth.StringWidth(ansi.Strip(selectCursor))
				cursor := strings.Repeat(" ", cwidth)
				fmt.Fprintf(s.Out, "%s%s\n", cursor, selectedOption.Key)
				continue
			}

			s.selectEval.val = selectedOption.Key
			fmt.Fprintf(s.Out, "%s%s\n", selectCursor, s.selectEval.Get())
		}
	}
}
