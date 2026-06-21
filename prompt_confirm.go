package pardon

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/grapheme"
)

// Confirm represents a yes/no confirmation prompt for user decisions.
type Confirm struct {
	PromptBase[bool, *Confirm]

	confirmKeyCluster grapheme.Cluster
	denyKeyCluster    grapheme.Cluster
	promptOpts        string
}

// NewConfirm creates a new Confirm prompt instance.
func NewConfirm(value *bool) *Confirm {
	c := &Confirm{
		PromptBase:        NewPromptBase[bool, *Confirm](value),
		confirmKeyCluster: grapheme.New('Y'),
		denyKeyCluster:    grapheme.New('N'),
	}
	c.PromptBase.Self = c
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

// ConfirmKey sets the [grapheme.Cluster] that represents the confirmation key
// (e.g., 'Y' for yes).
func (c *Confirm) ConfirmKey(key grapheme.Cluster) *Confirm {
	c.confirmKeyCluster = key
	return c
}

// DenyKey sets the [grapheme.Cluster] that represents the denial key (e.g., 'N'
// for no).
func (c *Confirm) DenyKey(key grapheme.Cluster) *Confirm {
	c.denyKeyCluster = key
	return c
}

func (c *Confirm) ask() error {
	// Need to save cursor location so we can easily redraw the prompt after
	// user input without needing to recalculate cursor movements.
	fmt.Fprint(c.Out, ansi.SaveCursorPos)

	c.buildAndSetPrompt()
	c.promptOpts = fmt.Sprintf("%s%s ", c.Prompt.String(), c.getPromptOptions())

	fmt.Fprint(c.Out, c.promptOpts)

	for {
		input, err := c.In.RawRead()
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
		confirmKey = bytes.ToUpper(c.confirmKeyCluster.Bytes())
		denyKey = bytes.ToLower(c.denyKeyCluster.Bytes())
	} else {
		confirmKey = bytes.ToLower(c.confirmKeyCluster.Bytes())
		denyKey = bytes.ToUpper(c.denyKeyCluster.Bytes())
	}

	var builder strings.Builder
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

func (c *Confirm) getValueAsCluster() grapheme.Cluster {
	if *c.value {
		return c.confirmKeyCluster
	}

	return c.denyKeyCluster
}

func (c *Confirm) getValueAsString() string {
	if *c.value {
		return string(c.confirmKeyCluster)
	}

	return string(c.denyKeyCluster)
}

func (c *Confirm) printFinalPromptLine() {
	c.answer.val = c.getValueAsString()
	promptAnswer := c.Prompt.String() + c.answer.Get()

	var builder strings.Builder
	builder.WriteString(ansi.RestoreCursorPos + ansi.ClearFromCursorToEndScreen)
	builder.WriteString(promptAnswer)
	builder.WriteString("\n")
	fmt.Fprint(c.Out, builder.String())
}

func (c *Confirm) processInput(input []byte) (done bool, err error) {
	if len(input) == 0 {
		return false, nil
	}

	c.PendingInputBytes = append(c.PendingInputBytes, input...)

	if needMoreInput := c.ConvertBytesToGraphemeSet(); needMoreInput {
		return false, nil
	}

	for len(c.PendingInputClusterSet) > 0 {
		r := c.PendingInputClusterSet[0]
		// If user hit enter, use the current value of [Confirm.value] as the input
		if grapheme.Equal(r, grapheme.Enter) || grapheme.Equal(r, grapheme.Newline) {
			r = c.getValueAsCluster()
		}

		switch {
		case grapheme.Equal(r, grapheme.CtrlC):
			return true, ErrUserAborted
		case grapheme.EqualFold(r, c.confirmKeyCluster):
			*c.value = true
			c.printFinalPromptLine()
			return true, nil
		case grapheme.EqualFold(r, c.denyKeyCluster):
			*c.value = false
			c.printFinalPromptLine()
			return true, nil
		}

		c.PendingInputClusterSet = c.PendingInputClusterSet[1:]
	}

	return false, nil
}
