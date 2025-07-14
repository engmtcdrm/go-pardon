package gocliselect

import (
	"fmt"
)

// Confirm is a struct that represents a confirmation prompt.
type Confirm struct {
	questionMark EvalVal[string]
	title        EvalVal[string]
	confirm      string
	confirmChars []byte
	deny         string
	denyChars    []byte
	exitChars    []byte
	value        *bool
}

func NewConfirm() *Confirm {
	return &Confirm{
		questionMark: EvalVal[string]{val: questionMark, fn: nil},
		title:        EvalVal[string]{val: "", fn: nil},
		confirm:      "Y",
		confirmChars: []byte{KeyYesUpper, KeyYes},
		deny:         "N",
		denyChars:    []byte{KeyNoUpper, KeyNo},
		exitChars:    []byte{KeyCtrlC, KeyEscape},
	}
}

// Set the title for the confirmation prompt.
func (c *Confirm) Title(title string) *Confirm {
	c.title.val = title
	c.title.fn = nil
	return c
}

func (c *Confirm) TitleFunc(fn func() string) *Confirm {
	c.title.fn = fn
	return c
}

func (c *Confirm) QuestionMark(s string) *Confirm {
	c.questionMark.val = s
	c.questionMark.fn = nil
	return c
}

func (c *Confirm) QuestionMarkFunc(fn func() string) *Confirm {
	c.questionMark.fn = fn
	return c
}

// Set the default value for the confirmation prompt.
func (c *Confirm) Value(value *bool) *Confirm {
	c.value = value
	return c
}

func (c *Confirm) Ask() error {
	if c.title.val == "" && c.title.fn == nil {
		return ErrNoTitle
	}

	if c.value == nil {
		return ErrNoValue
	}

	options := "[y/N]"
	if *c.value {
		options = "[Y/n]"
	}

	// Display the confirmation prompt
	fmt.Printf("%s %s %s: ", c.questionMark.Get(), c.title.Get(), options)

	// Capture user input
	for {
		keyCode := getInput()

		switch {
		case containsChar(c.confirmChars, keyCode):
			*c.value = true
			fmt.Println(c.confirm)
			return nil
		case containsChar(c.denyChars, keyCode):
			*c.value = false
			fmt.Println(c.deny)
			return nil
		case keyCode == KeyEnter:
			if *c.value {
				fmt.Println(c.confirm)
			} else {
				fmt.Println(c.deny)
			}
			return nil
		case containsChar(c.exitChars, keyCode):
			fmt.Println("")
			return ErrUserAborted
		}
	}
}
