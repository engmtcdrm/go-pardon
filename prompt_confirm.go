package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-pardon/keys"
	"github.com/engmtcdrm/go-pardon/tui"
)

// Confirm is a struct that represents a confirmation prompt.
type Confirm struct {
	icon     evalVal[string]
	title    evalVal[string]
	confirm  string
	deny     string
	value    *bool
	answerFn func(string) string
}

func NewConfirm() *Confirm {
	return &Confirm{
		icon:    evalVal[string]{val: Icons.Alert, fn: nil, defaultFn: defaultFuncs.iconFn},
		title:   evalVal[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		confirm: "Y",
		deny:    "N",
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

func (c *Confirm) formatFinalOutput(question string, answer string) string {
	return tui.RenderFormattedOutput(question, c.getAnswerFunc(answer))
}

func (c *Confirm) AnswerFunc(fn func(string) string) *Confirm {
	c.answerFn = fn
	return c
}

func (c *Confirm) getAnswerFunc(s string) string {
	if c.answerFn != nil {
		return c.answerFn(s)
	}

	if defaultFuncs.answerFn != nil {
		return defaultFuncs.answerFn(s)
	}

	return s
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
		keyCode := tui.GetInput()

		switch keyCode {
		case keys.KeyYesUpper, keys.KeyYes:
			*c.value = true
			fmt.Print(c.formatFinalOutput(question, c.confirm))
			return nil
		case keys.KeyNoUpper, keys.KeyNo:
			*c.value = false
			fmt.Print(c.formatFinalOutput(question, c.deny))
			return nil
		case keys.KeyEnter, keys.KeyCarriageReturn:
			if *c.value {
				fmt.Print(c.formatFinalOutput(question, c.confirm))
			} else {
				fmt.Print(c.formatFinalOutput(question, c.deny))
			}
			return nil
		case keys.KeyCtrlC, keys.KeyEscape:
			fmt.Println("")
			return ErrUserAborted
		}
	}
}
