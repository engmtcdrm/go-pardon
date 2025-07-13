package gocliselect

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"golang.org/x/term"
)

var (
	ErrNoSelectItems = errors.New("select has no items to display")
	ErrUserAborted   = errors.New("user aborted")
)

type Select struct {
	Prompt          string
	Cursor          string
	CursorPos       int
	ScrollOffset    int
	SelectItems     []*SelectItem
	ItemSelectColor func(...any) string
}

type SelectItem struct {
	Key       string
	Value     any
	SubSelect *Select
}

func NewSelect(prompt string) *Select {
	return &Select{
		Prompt:      prompt,
		Cursor:      ">",
		SelectItems: make([]*SelectItem, 0),
	}
}

// AddItem will add a new select option to the select list
func (s *Select) AddItem(option string, id any) *Select {
	selectItem := &SelectItem{
		Key:   option,
		Value: id,
	}

	s.SelectItems = append(s.SelectItems, selectItem)
	return s
}

// renderSelectItems prints the select item list.
// Setting redraw to true will re-render the options list with updated current selection.
func (s *Select) renderSelectItems(redraw bool) {
	termHeight := 25 // Default height

	// Try to get terminal size, but don't fail if we can't
	if _, height, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		termHeight = height
	}

	termHeight = termHeight - 3 // Space for prompt and cursor movement
	selectSize := len(s.SelectItems)

	// Ensure scroll offset follows cursor movement
	if s.CursorPos < s.ScrollOffset {
		s.ScrollOffset = s.CursorPos
	} else if s.CursorPos >= s.ScrollOffset+termHeight {
		s.ScrollOffset = s.CursorPos - termHeight + 1
	}

	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Print(ansi.CursorUp(min(selectSize, termHeight)))
	}

	selectCursor := fmt.Sprintf("%s ", s.Cursor)

	// Render only visible select items
	for i := s.ScrollOffset; i < min(s.ScrollOffset+termHeight, selectSize); i++ {
		selectItem := s.SelectItems[i]
		cursor := strings.Repeat(" ", len(selectCursor))

		fmt.Print(ansi.ClearLine)

		if i == s.CursorPos {
			cursor = s.ItemSelectColor(selectCursor)
			fmt.Printf("\r%s%s\n", cursor, s.ItemSelectColor(selectItem.Key))
		} else {
			fmt.Printf("\r%s%s\n", cursor, selectItem.Key)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Display will display the current select options and awaits user selection
// It returns the users selected choice
func (s *Select) Display() (interface{}, error) {
	defer func() {
		fmt.Print(ansi.ShowCursor)
	}()

	if len(s.SelectItems) == 0 {
		return nil, ErrNoSelectItems
	}

	fmt.Println(s.Prompt)

	s.renderSelectItems(false)

	fmt.Print(ansi.HideCursor)

	for {
		keyCode := getInput()

		switch keyCode {
		case KeyEscape:
			return "", nil
		case KeyCtrlC:
			return "", ErrUserAborted
		case KeyEnter, KeyCarriageReturn:
			selectItem := s.SelectItems[s.CursorPos]
			fmt.Println("\r")
			return selectItem.Value, nil
		case KeyUp:
			s.CursorPos = (s.CursorPos + len(s.SelectItems) - 1) % len(s.SelectItems)
			s.renderSelectItems(true)
		case KeyDown:
			s.CursorPos = (s.CursorPos + 1) % len(s.SelectItems)
			s.renderSelectItems(true)
		}
	}
}

// getInput will read raw input from the terminal
// It returns the raw ASCII value inputted
func getInput() byte {
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
	if read == 3 {
		if _, ok := NavigationKeys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}
