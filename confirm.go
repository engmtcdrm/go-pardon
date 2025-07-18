package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-ansi"
)

// Confirm is a struct that represents a confirmation prompt.
type Confirm struct {
	icon         evalVal[string]
	title        evalVal[string]
	confirm      string
	confirmChars []byte
	deny         string
	denyChars    []byte
	enterChars   []byte
	exitChars    []byte
	value        *bool
	answerFn     func(string) string
}

func NewConfirm() *Confirm {
	return &Confirm{
		icon:         evalVal[string]{val: Icons.Alert, fn: nil, defaultFn: defaultFuncs.iconFn},
		title:        evalVal[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		confirm:      "Y",
		deny:         "N",
		confirmChars: []byte{keyYesUpper, keyYes},
		denyChars:    []byte{keyNoUpper, keyNo},
		enterChars:   []byte{keyEnter, keyCarriageReturn},
		exitChars:    []byte{keyCtrlC, keyEscape},
	}
}

// Set the title for the confirmation prompt.
func (c *Confirm) Title(title string) *Confirm {
	c.title.val = title
	c.title.fn = nil
	return c
}

func (c *Confirm) TitleFunc(fn func(string) string) *Confirm {
	c.title.fn = fn
	return c
}

func (c *Confirm) Icon(s string) *Confirm {
	c.icon.val = s
	c.icon.fn = nil
	return c
}

func (c *Confirm) IconFunc(fn func(string) string) *Confirm {
	c.icon.fn = fn
	return c
}

// Set the default value for the confirmation prompt.
func (c *Confirm) Value(value *bool) *Confirm {
	c.value = value
	return c
}

func (c *Confirm) formatFinalOutput(question string, result string) string {
	return fmt.Sprintf("%s\r%s %s\n", ansi.ClearToBegin, question, c.getAnswerFunc(result))
}

func (c *Confirm) AnswerFunc(fn func(string) string) *Confirm {
	c.answerFn = fn
	return c
}

func (c *Confirm) getAnswerFunc(text string) string {
	if c.answerFn != nil {
		return c.answerFn(text)
	}

	if defaultFuncs.answerFn != nil {
		return defaultFuncs.answerFn(text)
	}

	return text
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

	question := fmt.Sprintf("%s%s", c.icon.Get(), c.title.Get())
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
