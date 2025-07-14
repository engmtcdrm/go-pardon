package gocliselect

import (
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"golang.org/x/term"
)

var (
	// NavigationKeys defines a map of specific byte keycodes related to
	// navigation functionality, such as up or down actions.
	NavigationKeys = map[byte]bool{
		KeyUp:   true,
		KeyDown: true,
	}
)

type Select[T comparable] struct {
	title           EvalVal[string]
	cursor          EvalVal[string]
	cursorPos       int
	scrollOffset    int
	items           []Item[T]
	ItemSelectColor func(...any) string
	value           *T
}

func NewSelect[T comparable]() *Select[T] {
	return &Select[T]{
		title:  EvalVal[string]{val: "", fn: nil},
		cursor: EvalVal[string]{val: ">", fn: nil},
		items:  make([]Item[T], 0),
	}
}

func (s *Select[T]) Title(title string) *Select[T] {
	s.title.val = title
	s.title.fn = nil
	return s
}

func (s *Select[T]) TitleFunc(fn func() string) *Select[T] {
	s.title.fn = fn
	return s
}

func (s *Select[T]) Cursor(cursor string) *Select[T] {
	s.cursor.val = cursor
	s.cursor.fn = nil
	return s
}

func (s *Select[T]) CursorFunc(fn func() string) *Select[T] {
	s.cursor.fn = fn
	return s
}

func (s *Select[T]) Items(items ...Item[T]) *Select[T] {
	s.items = items
	return s
}

func (s *Select[T]) Value(value *T) *Select[T] {
	s.value = value
	return s
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

	if len(s.items) == 0 {
		return ErrNoSelectItems
	}

	defer func() {
		fmt.Print(ansi.ShowCursor)
	}()

	fmt.Println(s.title.Get())

	s.renderSelectItems(false)

	fmt.Print(ansi.HideCursor)

	for {
		keyCode := getInput()

		switch keyCode {
		case KeyEscape:
			return nil
		case KeyCtrlC:
			return ErrUserAborted
		case KeyEnter, KeyCarriageReturn:
			*s.value = s.items[s.cursorPos].Value
			fmt.Println("\r")
			return nil
		case KeyUp:
			s.cursorPos = (s.cursorPos + len(s.items) - 1) % len(s.items)
			s.renderSelectItems(true)
		case KeyDown:
			s.cursorPos = (s.cursorPos + 1) % len(s.items)
			s.renderSelectItems(true)
		}
	}
}

// renderSelectItems prints the select item list.
// Setting redraw to true will re-render the options list with updated current selection.
func (s *Select[T]) renderSelectItems(redraw bool) {
	termHeight := 25 // Default height

	// Try to get terminal size, but don't fail if we can't
	if _, height, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		termHeight = height
	}

	termHeight = termHeight - 3 // Space for prompt and cursor movement
	selectSize := len(s.items)

	// Ensure scroll offset follows cursor movement
	if s.cursorPos < s.scrollOffset {
		s.scrollOffset = s.cursorPos
	} else if s.cursorPos >= s.scrollOffset+termHeight {
		s.scrollOffset = s.cursorPos - termHeight + 1
	}

	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Print(ansi.CursorUp(min(selectSize, termHeight)))
	}

	selectCursor := fmt.Sprintf("%s ", s.cursor.Get())

	// Render only visible select items
	for i := s.scrollOffset; i < min(s.scrollOffset+termHeight, selectSize); i++ {
		selectItem := s.items[i]
		cursor := strings.Repeat(" ", len(selectCursor))

		fmt.Print(ansi.ClearLine)

		if i == s.cursorPos {
			cursor = s.ItemSelectColor(selectCursor)
			fmt.Printf("\r%s%s\n", cursor, s.ItemSelectColor(selectItem.Key))
		} else {
			fmt.Printf("\r%s%s\n", cursor, selectItem.Key)
		}
	}
}
