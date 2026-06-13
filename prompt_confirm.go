package pardon

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/internal/keys"
)

// Confirm represents a yes/no confirmation prompt for user decisions.
type Confirm struct {
	terminal   *Terminal
	value      *bool
	icon       eval[string]
	title      eval[string]
	answer     eval[string]
	confirmKey []byte
	denyKey    []byte
	prompt     string
	promptOpts string
}

// NewConfirm creates a new Confirm prompt instance.
func NewConfirm(value *bool) *Confirm {
	return &Confirm{
		terminal:   NewTerminal(),
		value:      value,
		icon:       eval[string]{val: Icons.QuestionMark, fn: nil, defaultFn: defaultFuncs.iconFn},
		title:      eval[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		answer:     eval[string]{val: "", fn: nil, defaultFn: defaultFuncs.answerFn},
		confirmKey: []byte{keys.UpperY},
		denyKey:    []byte{keys.UpperN},
	}
}

// AnswerFunc sets a function to transform the final answer being displayed.
func (c *Confirm) AnswerFunc(fn func(string) string) *Confirm {
	c.answer.fn = fn
	return c
}

// Ask displays the confirmation prompt.
func (c *Confirm) Ask() error {
	if c.title.val == "" && c.title.fn == nil {
		return ErrNoTitle
	}

	if c.value == nil {
		return ErrNoValue
	}

	if err := c.ask(); err != nil {
		return err
	}

	return nil
}

// ConfirmKey sets the rune that represents the confirmation key (e.g., 'Y' for
// yes).
func (c *Confirm) ConfirmKey(key byte) *Confirm {
	c.confirmKey = []byte{key}
	return c
}

// DenyKey sets the rune that represents the denial key (e.g., 'N' for no).
func (c *Confirm) DenyKey(key byte) *Confirm {
	c.denyKey = []byte{key}
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

// Value sets a default value for the confirmation prompt.
func (c *Confirm) Value(value *bool) *Confirm {
	c.value = value
	return c
}

func (c *Confirm) ask() error {
	c.prompt = fmt.Sprintf("%s%s ", c.icon.Get(), c.title.Get())
	c.promptOpts = fmt.Sprintf("%s%s ", c.prompt, c.getPromptOptions())

	c.terminal.Print(c.promptOpts)

	for {
		input, err := c.terminal.RawRead()
		if err != nil {
			return err
		}

		if done, err := c.processInput(input); done {
			return err
		}
	}
}

func (c *Confirm) getPromptOptions() string {
	var confirmKey, denyKey []byte
	if *c.value {
		confirmKey = bytes.ToUpper(c.confirmKey)
		denyKey = bytes.ToLower(c.denyKey)
	} else {
		confirmKey = bytes.ToLower(c.confirmKey)
		denyKey = bytes.ToUpper(c.denyKey)
	}

	builder := strings.Builder{}
	builder.WriteString("[")

	if len(confirmKey) > 0 {
		builder.Write(confirmKey)
	} else {
		builder.WriteString("?")
	}

	builder.WriteString("/")

	if len(denyKey) > 0 {
		builder.Write(denyKey)
	} else {
		builder.WriteString("?")
	}

	builder.WriteString("]")

	return builder.String()
}

func (c *Confirm) getValueAsBytes() []byte {
	if *c.value {
		return c.confirmKey
	}

	return c.denyKey
}

func (c *Confirm) getValueAsString() string {
	if *c.value {
		return string(c.confirmKey)
	}

	return string(c.denyKey)
}

func (c *Confirm) printFinalPromptLine() {
	builder := strings.Builder{}
	builder.WriteString(ansi.ClearLineReset)
	c.answer.val = c.getValueAsString()
	promptAnswer := c.prompt + c.answer.Get()
	builder.WriteString(promptAnswer)
	builder.WriteString("\n" + ansi.ClearLineReset)
	c.terminal.Print(builder.String())
}

func (c *Confirm) processInput(input []byte) (done bool, err error) {
	if len(input) == 0 {
		return false, nil
	}

	// If user hit enter, use the current value of [Confirm.value] as the input
	if input[0] == keys.Enter || input[0] == keys.NewLine {
		input = c.getValueAsBytes()
	}

	switch {
	case bytes.Equal([]byte{keys.CtrlC}, input):
		return true, ErrUserAborted
	case bytes.EqualFold(c.confirmKey, input):
		*c.value = true
		c.printFinalPromptLine()
		return true, nil
	case bytes.EqualFold(c.denyKey, input):
		*c.value = false
		c.printFinalPromptLine()
		return true, nil
	}

	return false, nil
}
