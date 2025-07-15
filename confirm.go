package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-ansi"
)

// Confirm is a struct that represents a confirmation prompt.
type Confirm struct {
	questionMark EvalVal[string]
	title        EvalVal[string]
	confirm      string
	confirmChars []byte
	deny         string
	denyChars    []byte
	enterChars   []byte
	exitChars    []byte
	value        *bool
}

func NewConfirm() *Confirm {
	return &Confirm{
		questionMark: EvalVal[string]{val: questionMarkIcon, fn: nil},
		title:        EvalVal[string]{val: "", fn: nil},
		confirm:      "Y",
		deny:         "N",

		confirmChars: []byte{KeyYesUpper, KeyYes},
		denyChars:    []byte{KeyNoUpper, KeyNo},
		enterChars:   []byte{KeyEnter, KeyCarriageReturn},
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

func (c *Confirm) formatFinalOutput(question string, result string) string {
	return fmt.Sprintf("%s\r%s %s\n", ansi.ClearToBegin, question, result)
}

func (c *Confirm) Ask() error {
	if c.title.val == "" && c.title.fn == nil {
		return ErrNoTitle
	}

	if c.value == nil {
		return ErrNoValue
	}

	options := "(y/N)"
	if *c.value {
		options = "(Y/n)"
	}

	question := fmt.Sprintf("%s %s", c.questionMark.Get(), c.title.Get())
	question_opt := fmt.Sprintf("%s %s ", question, options)

	// Display the confirmation prompt
	fmt.Print(question_opt)

	// Capture user input
	for {
		keyCode := getInput()

		switch {
		case containsChar(c.confirmChars, keyCode):
			*c.value = true
			fmt.Print(c.formatFinalOutput(question, c.confirm))
			return nil
		case containsChar(c.denyChars, keyCode):
			*c.value = false
			fmt.Print(c.formatFinalOutput(question, c.deny))
			return nil
		case containsChar(c.enterChars, keyCode):
			if *c.value {
				fmt.Print(c.formatFinalOutput(question, c.confirm))
			} else {
				fmt.Print(c.formatFinalOutput(question, c.deny))
			}
			return nil
		case containsChar(c.exitChars, keyCode):
			fmt.Println("")
			return ErrUserAborted
		}
	}
}
