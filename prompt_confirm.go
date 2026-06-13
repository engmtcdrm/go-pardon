package pardon

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/internal/runekeys"
)

// Confirm represents a yes/no confirmation prompt for user decisions.
type Confirm struct {
	terminal   *Terminal
	value      *bool
	icon       eval[string]
	title      eval[string]
	answer     eval[string]
	confirmKey rune
	denyKey    rune
	prompt     string
	promptOpts string
}

// NewConfirm creates a new Confirm prompt instance.
func NewConfirm(value *bool) *Confirm {
	return &Confirm{
		terminal:   NewConfirmTerminal(),
		value:      value,
		icon:       eval[string]{val: Icons.QuestionMark, fn: nil, defaultFn: defaultFuncs.iconFn},
		title:      eval[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		answer:     eval[string]{val: "", fn: nil, defaultFn: defaultFuncs.answerFn},
		confirmKey: runekeys.UpperY,
		denyKey:    runekeys.UpperN,
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

	c.terminal.CustomHandler = func(t *Terminal, r rune) (done bool) {
		return c.processLine([]rune{r})
	}

	if err := c.ask(); err != nil {
		return err
	}

	return nil
}

// ConfirmKey sets the rune that represents the confirmation key (e.g., 'Y' for
// yes).
func (c *Confirm) ConfirmKey(r rune) *Confirm {
	c.confirmKey = r
	return c
}

// DenyKey sets the rune that represents the denial key (e.g., 'N' for no).
func (c *Confirm) DenyKey(r rune) *Confirm {
	c.denyKey = r
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
		c.terminal.Reset()
		line, err := c.terminal.RawRead()
		if err != nil {
			if errors.Is(err, ErrUserAborted) {
				c.terminal.Print(ansi.ClearLineReset + c.prompt)
				return err
			}

			return err
		}

		if done := c.processLine(line); done {
			break
		}
	}

	c.printFinalPromptLine()

	return nil
}

func (c *Confirm) equal(a rune, b rune) bool {
	return unicode.ToLower(a) == unicode.ToLower(b)
}

func (c *Confirm) getPromptOptions() string {
	var confirmKey, denyKey rune
	if *c.value {
		confirmKey = unicode.ToUpper(c.confirmKey)
		denyKey = unicode.ToLower(c.denyKey)
	} else {
		confirmKey = unicode.ToLower(c.confirmKey)
		denyKey = unicode.ToUpper(c.denyKey)
	}

	return fmt.Sprintf("[%c/%c]", confirmKey, denyKey)
}

func (c *Confirm) getValueAsRunes() []rune {
	if *c.value {
		return []rune{c.confirmKey}
	}

	return []rune{c.denyKey}
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

func (c *Confirm) processLine(line []rune) (done bool) {
	if len(line) == 0 {
		return false
	}

	// If user hit enter, use the current value of [Confirm.value] as the input
	if line[0] == runekeys.Enter || line[0] == runekeys.NewLine {
		line = c.getValueAsRunes()
	}

	switch {
	case c.equal(line[0], c.confirmKey):
		*c.value = true
		return true
	case c.equal(line[0], c.denyKey):
		*c.value = false
		return true
	}

	return false
}
