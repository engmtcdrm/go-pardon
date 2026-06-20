package pardon

import (
	"fmt"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/grapheme"
)

type Text struct {
	PromptBase[string, *Text]

	validateFn func(string) error

	// hide indicates whether the input should be hidden (e.g., for password
	// input).
	hide bool

	defaultValue           grapheme.ClusterSet
	promptDefault          string
	pendingValueClusterSet grapheme.ClusterSet
}

// NewPassword creates an InputPrompt for secure password input with masking.
func NewPassword(value *string) *Text {
	t := &Text{
		PromptBase: NewPromptBase[string, *Text](value),
		validateFn: func(s string) error { return nil },
		hide:       true,
	}
	t.PromptBase.Self = t
	return t
}

// NewQuestion creates a new InputPrompt for text input with a question mark
// icon.
func NewQuestion(value *string) *Text {
	t := &Text{
		PromptBase: NewPromptBase[string, *Text](value),
		validateFn: func(s string) error { return nil },
	}
	t.PromptBase.Self = t
	t.setDefaultValue()
	return t
}

// Ask displays the prompt and waits for input.
func (t *Text) Ask() error {
	if t.title.val == "" && t.title.fn == nil {
		return ErrNoTitle
	}

	if t.value == nil {
		return ErrNoValue
	}

	if err := t.ask(); err != nil {
		return err
	}

	return nil
}

// Hide sets whether the input should be hidden (e.g., for password input).
func (t *Text) Hide(hide bool) *Text {
	t.hide = hide
	return t
}

// ValidateFunc sets a validation function for the prompt input.
func (t *Text) ValidateFunc(fn func(string) error) *Text {
	if fn != nil {
		t.validateFn = fn
	}
	return t
}

// Value sets a default value for the prompt.
func (t *Text) Value(value *string) *Text {
	t.value = value
	t.setDefaultValue()
	return t
}

// ask handles the core logic of displaying the prompt, reading user input,
// validating it, and applying the answer transformation.
func (t *Text) ask() error {
	// Need to save cursor location so we can easily redraw the prompt after
	// user input without needing to recalculate cursor movements.
	fmt.Fprint(t.Out, saveCursor)

	if t.hide {
		fmt.Fprint(t.Out, ansi.HideCursor)
		defer func() {
			fmt.Fprint(t.Out, ansi.ShowCursor)
		}()
	}

	t.buildAndSetPrompt()
	fmt.Fprint(t.Out, t.getPromptWithDefaultValue())

	for {
		input, err := t.In.RawRead()
		if err != nil {
			return err
		}

		if done, err := t.processInput(input); done {
			return err
		}
	}
}

func (t *Text) setDefaultValue() {
	if t.value != nil {
		t.defaultValue = grapheme.ClusterSetFromString(*t.value)
	}
}

func (t *Text) getPromptWithDefaultValue() string {
	var prompt strings.Builder
	prompt.WriteString(ansi.ClearLineReset)
	prompt.WriteString(t.Prompt.String())

	if len(t.defaultValue) > 0 {
		t.promptDefault = fmt.Sprintf("%s%s%s%s",
			ansi.Dim,
			t.defaultValue.String(),
			ansi.Reset,
			ansi.CursorBackward(t.defaultValue.VisualLen()),
		)

		prompt.WriteString(t.promptDefault)
	}

	return prompt.String()
}

func (t *Text) printErrorMessage(err error) {
	errMsg := validationErrorMessage(err)

	var builder strings.Builder
	builder.WriteString(restoreCursor + ansi.ClearFromCursorToEndScreen)
	builder.WriteString(t.getPromptWithDefaultValue())
	builder.WriteString("\n")
	builder.WriteString(errMsg)
	builder.WriteString(restoreCursor)
	builder.WriteString(ansi.CursorHorizontalAbsolute(t.Prompt.VisualLen() + 1))
	fmt.Fprint(t.Out, builder.String())
}

// printFinalPromptLine handles printing the final prompt line after successful
// input.
func (t *Text) printFinalPromptLine() {
	var builder strings.Builder
	builder.WriteString(restoreCursor + ansi.ClearFromCursorToEndScreen)
	builder.WriteString(t.Prompt.String())

	// If the input is not hidden, We need to clear the line, then print the
	// prompt with the answer function applied.
	if !t.hide {
		t.answer.val = *t.value
		builder.WriteString(t.answer.Get())
	}

	builder.WriteString("\n")
	fmt.Fprint(t.Out, builder.String())
}

func (t *Text) handleEnter(_ grapheme.Cluster) (done bool, err error) {
	if err := t.validateFn(t.pendingValueClusterSet.String()); err != nil {
		t.printErrorMessage(err)
		t.PendingInputClusterSet = nil
		t.pendingValueClusterSet = nil
		return false, nil
	}

	// Only set value if input is not empty. Otherwise we will use the default
	// value.
	if len(t.pendingValueClusterSet) > 0 {
		*t.value = t.pendingValueClusterSet.String()
	}

	t.printFinalPromptLine()
	return true, nil
}

func (t *Text) handleDelete(_ grapheme.Cluster) (done bool, err error) {
	if len(t.pendingValueClusterSet) > 0 {
		t.pendingValueClusterSet = t.pendingValueClusterSet[:len(t.pendingValueClusterSet)-1]
		t.printInput("\b \b")
	}

	if len(t.pendingValueClusterSet) == 0 {
		fmt.Fprint(t.Out, t.getPromptWithDefaultValue())
	}

	t.PendingInputClusterSet = t.PendingInputClusterSet[1:]
	return false, nil
}

// Currently calls [Text.handleEnter].
func (t *Text) handleNewline(r grapheme.Cluster) (done bool, err error) {
	return t.handleEnter(r)
}

func (t *Text) processInput(input []byte) (done bool, err error) {
	if len(input) == 0 {
		return false, nil
	}

	t.PendingInputBytes = append(t.PendingInputBytes, input...)

	if needMoreInput := t.ConvertBytesToGraphemeSet(); needMoreInput {
		return false, nil
	}

	for len(t.PendingInputClusterSet) > 0 {
		r := t.PendingInputClusterSet[0]
		switch {
		case grapheme.Equal(r, grapheme.CtrlC):
			return true, ErrUserAborted
		case grapheme.Equal(r, grapheme.Enter):
			return t.handleEnter(r)
		case grapheme.Equal(r, grapheme.Newline):
			return t.handleNewline(r)
		case grapheme.Equal(r, grapheme.Delete), grapheme.Equal(r, grapheme.Backspace):
			return t.handleDelete(r)
		case grapheme.Equal(r, grapheme.Escape):
			doContinue, err := t.ProcessEscapeSequence(r, nil)
			if !doContinue {
				return false, err
			}
			continue
		}

		t.pendingValueClusterSet = append(t.pendingValueClusterSet, r)
		t.PendingInputClusterSet = t.PendingInputClusterSet[1:]
		fmt.Fprint(t.Out, ansi.ClearLineReset+t.Prompt.String())
		t.printInput(t.pendingValueClusterSet.String())
	}

	t.PendingInputClusterSet = nil

	return false, nil
}

func (t *Text) printInput(a ...any) {
	if !t.hide {
		fmt.Fprint(t.Out, a...)
	}
}
