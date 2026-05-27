package pardon

import (
	"errors"
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/internal/keys"
)

// Confirm represents a yes/no confirmation prompt for user decisions.
type Confirm struct {
	icon       eval[string]
	title      eval[string]
	confirm    string
	deny       string
	value      *bool
	answerFn   func(string) string
	input      *Input
	prompt     string
	promptOpts string
}

// NewConfirm creates a new Confirm prompt instance.
func NewConfirm(value *bool) *Confirm {
	return &Confirm{
		icon:    eval[string]{val: Icons.QuestionMark, fn: nil, defaultFn: defaultFuncs.iconFn},
		title:   eval[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		confirm: "Y",
		deny:    "N",
		value:   value,
		input:   NewConfirmInput(),
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

	c.prompt = fmt.Sprintf("%s%s ", c.icon.Get(), c.title.Get())
	c.promptOpts = fmt.Sprintf("%s%s ", c.prompt, options)

	if err := c.ask(); err != nil {
		return err
	}

	return nil
}

func (c *Confirm) ask() error {
	fmt.Print(c.promptOpts)

outer:
	for {
		line, err := c.input.RawRead()
		if err != nil {
			if errors.Is(err, ErrUserAborted) {
				fmt.Print(ansi.ClearLineReset + c.promptOpts)
				return err
			}
			return err
		}

		// If user hit enter, user the default value
		if line == "" {
			line = c.getDefaultValue()
		}

		pendingValue := strings.TrimSpace(line)
		switch pendingValue[0] {
		case keys.UpperY, keys.LowerY:
			*c.value = true
			break outer
		case keys.UpperN, keys.LowerN:
			*c.value = false
			break outer
		default:
			continue
		}
	}

	c.printFinalPromptLine()

	return nil
}

func (c *Confirm) printFinalPromptLine() {
	builder := strings.Builder{}
	builder.WriteString(ansi.ClearLineReset)

	if *c.value {
		builder.WriteString(c.prompt + c.setAnswerFunc(c.confirm))
	} else {
		builder.WriteString(c.prompt + c.setAnswerFunc(c.deny))
	}

	builder.WriteString("\n" + ansi.ClearLineReset)
	fmt.Print(builder.String())
}

func (c *Confirm) getDefaultValue() string {
	if *c.value {
		return c.confirm
	}

	return c.deny
}
