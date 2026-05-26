package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-pardon/internal/keys"
	"github.com/engmtcdrm/go-pardon/internal/tui"
)

// Confirm represents a yes/no confirmation prompt for user decisions.
type Confirm struct {
	icon     eval[string]
	title    eval[string]
	confirm  string
	deny     string
	value    *bool
	answerFn func(string) string
}

// NewConfirm creates a new Confirm prompt instance.
func NewConfirm(value *bool) *Confirm {
	return &Confirm{
		icon:    eval[string]{val: Icons.QuestionMark, fn: nil, defaultFn: defaultFuncs.iconFn},
		title:   eval[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		confirm: "Y",
		deny:    "N",
		value:   value,
	}
}

// Title sets a static title for the confirmation prompt.
func (c *Confirm) Title(title string) *Confirm {
	c.title.val = title
	return c
}

// TitleFunc sets a dynamic title function for the confirmation prompt.
func (c *Confirm) TitleFunc(fn func(string) string) *Confirm {
	c.title.fn = fn
	return c
}

// Icon sets a static icon for the confirmation prompt.
func (c *Confirm) Icon(s string) *Confirm {
	c.icon.val = s
	c.icon.fn = nil
	return c
}

// IconFunc sets a dynamic icon function for the confirmation prompt.
func (c *Confirm) IconFunc(fn func(string) string) *Confirm {
	c.icon.fn = fn
	return c
}

// Value sets a default value for the confirmation prompt.
func (c *Confirm) Value(value *bool) *Confirm {
	c.value = value
	return c
}

// AnswerFunc sets a function to transform the final answer before returning.
func (c *Confirm) AnswerFunc(fn func(string) string) *Confirm {
	c.answerFn = fn
	return c
}

// formatFinalOutput formats the final confirmation display after user selection.
func (c *Confirm) formatFinalOutput(question string, answer string) string {
	return tui.RenderFormattedOutput(question, c.setAnswerFunc(answer))
}

// setAnswerFunc configures the answer transformation priority:
// prompt-specific, global default, or the string itself.
func (c *Confirm) setAnswerFunc(s string) string {
	if c.answerFn != nil {
		return c.answerFn(s)
	}

	if defaultFuncs.answerFn != nil {
		return defaultFuncs.answerFn(s)
	}

	return s
}

// Ask displays the confirmation prompt.
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

	question := fmt.Sprintf("%s%s", c.icon.Get(), c.title.Get())
	question_opt := fmt.Sprintf("%s %s ", question, options)

	// Display the confirmation prompt
	fmt.Print(question_opt)

	// Capture user input
	for {
		keyCode := tui.GetInput()

		switch keyCode {
		case keys.UpperY, keys.LowerY:
			*c.value = true
			fmt.Print(c.formatFinalOutput(question, c.confirm))
			return nil
		case keys.UpperN, keys.LowerN:
			*c.value = false
			fmt.Print(c.formatFinalOutput(question, c.deny))
			return nil
		case keys.Enter, keys.NewLine:
			if *c.value {
				fmt.Print(c.formatFinalOutput(question, c.confirm))
			} else {
				fmt.Print(c.formatFinalOutput(question, c.deny))
			}
			return nil
		case keys.CtrlC, keys.Escape:
			fmt.Println()
			return ErrUserAborted
		}
	}
}
