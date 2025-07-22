package pardon

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/keys"
	"github.com/engmtcdrm/go-pardon/tui"
)

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

func NewSelect[T comparable]() *Select[T] {
	return &Select[T]{
		icon:    eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title:   eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		cursor:  eval[string]{val: "> ", defaultFn: defaultFuncs.cursorFn},
		options: make([]Option[T], 0),
	}
}

func (s *Select[T]) Title(title string) *Select[T] {
	s.title.val = title
	s.title.fn = nil
	return s
}

func (s *Select[T]) TitleFunc(fn func(string) string) *Select[T] {
	s.title.fn = fn
	return s
}

func (s *Select[T]) Cursor(cursor string) *Select[T] {
	s.cursor.val = cursor
	s.cursor.fn = nil
	return s
}

func (s *Select[T]) CursorFunc(fn func(string) string) *Select[T] {
	s.cursor.fn = fn
	return s
}

func (s *Select[T]) Options(options ...Option[T]) *Select[T] {
	if len(options) == 0 {
		return s
	}

	s.options = options
	return s
}

func (s *Select[T]) Value(value *T) *Select[T] {
	s.value = value
	return s
}

func (s *Select[T]) Icon(icon string) *Select[T] {
	s.icon.val = icon
	s.icon.fn = nil
	return s
}

func (s *Select[T]) IconFunc(fn func(string) string) *Select[T] {
	s.icon.fn = fn
	return s
}

// AnswerFunc allows customization of the selected option's display format.
func (s *Select[T]) AnswerFunc(fn func(string) string) *Select[T] {
	s.answerFn = fn
	return s
}

func (s *Select[T]) getSelectFunc(text string) string {
	if s.selectFn != nil {
		return s.selectFn(text)
	}

	if defaultFuncs.selectFn != nil {
		return defaultFuncs.selectFn(text)
	}

	return text
}

func (s *Select[T]) getAnswerFunc(answer string) string {
	if s.answerFn != nil {
		return s.answerFn(answer)
	}

	if defaultFuncs.answerFn != nil {
		return defaultFuncs.answerFn(answer)
	}

	return answer
}

// Ask will display the current select options and awaits user selection
// It returns the users selected choice
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
		fmt.Print(ansi.ShowCursor)
	}()

	// Print the question
	fmt.Printf("%s%s\n", s.icon.Get(), s.title.Get())

	s.renderOptions(false)
	fmt.Print(ansi.HideCursor)

	for {
		keyCode := tui.GetInput()

		switch keyCode {
		case keys.KeyCtrlC:
			return ErrUserAborted
		case keys.KeyEnter, keys.KeyCarriageReturn:
			*s.value = s.options[s.cursorPos].Value
			visibleOptions := tui.Min(len(s.options), tui.GetTerminalHeight()-3)
			tui.RenderClearAndReposition(visibleOptions+1, s.icon.Get(), s.title.Get(), s.getAnswerFunc(s.options[s.cursorPos].Key))
			return nil
		case keys.KeyUp:
			s.cursorPos = (s.cursorPos + len(s.options) - 1) % len(s.options)
			s.renderOptions(true)
		case keys.KeyDown:
			s.cursorPos = (s.cursorPos + 1) % len(s.options)
			s.renderOptions(true)
		}
	}
}

// Setting redraw to true will re-render the options list with updated current selection.
func (s *Select[T]) renderOptions(redraw bool) {
	termHeight := tui.GetTerminalHeight()
	termHeight = termHeight - 3 // Space for prompt and cursor movement
	selectSize := len(s.options)

	// Ensure scroll offset follows cursor movement
	if s.cursorPos < s.scrollOffset {
		s.scrollOffset = s.cursorPos
	} else if s.cursorPos >= s.scrollOffset+termHeight {
		s.scrollOffset = s.cursorPos - termHeight + 1
	}

	selectCursor := s.cursor.Get()
	visibleLines := tui.Min(selectSize, termHeight)

	if redraw {
		// For terminal optimization: build entire output first, then write atomically
		var output strings.Builder

		// Move cursor up to start position
		output.WriteString(ansi.CursorUp(visibleLines))

		// Build all lines in memory first
		for i := s.scrollOffset; i < tui.Min(s.scrollOffset+termHeight, selectSize); i++ {
			selectedOption := s.options[i]
			cursor := strings.Repeat(" ", utf8.RuneCountInString(ansi.StripCodes(selectCursor)))

			// Clear line and build content
			output.WriteString("\r")
			output.WriteString(ansi.ClearLine)

			if i == s.cursorPos {
				cursor = s.getSelectFunc(selectCursor)
				output.WriteString(cursor)
				output.WriteString(s.getSelectFunc(selectedOption.Key))
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
		for i := s.scrollOffset; i < tui.Min(s.scrollOffset+termHeight, selectSize); i++ {
			selectedOption := s.options[i]
			cursor := strings.Repeat(" ", utf8.RuneCountInString(ansi.StripCodes(selectCursor)))

			if i == s.cursorPos {
				cursor = s.getSelectFunc(selectCursor)
				fmt.Printf("%s%s\n", cursor, s.getSelectFunc(selectedOption.Key))
			} else {
				fmt.Printf("%s%s\n", cursor, selectedOption.Key)
			}
		}
	}
}
