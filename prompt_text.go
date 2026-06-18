package pardon

import (
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/grapheme"
	"golang.org/x/term"
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

// ask handles the core logic of displaying the prompt, reading user input,
// validating it, and applying the answer transformation.
func (t *Text) ask() error {
	if t.hide {
		fmt.Fprint(t.Out, ansi.HideCursor)
		defer func() {
			fmt.Fprint(t.Out, ansi.ShowCursor)
		}()
	}

	t.Prompt = fmt.Sprintf("%s%s ", t.icon.Get(), t.title.Get())
	fmt.Fprint(t.Out, t.getPrompt())

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
	// t.defaultValue = grapheme.ClusterSetFromString(ansi.Strip(*t.value))
	t.defaultValue = grapheme.ClusterSetFromString(*t.value)
}

func (t *Text) getPrompt() string {
	var prompt strings.Builder
	prompt.WriteString(ansi.ClearLineReset)
	prompt.WriteString(t.Prompt)

	if len(t.defaultValue) > 0 {
		t.promptDefault = fmt.Sprintf("%s%s%s%s ",
			ansi.Dim,
			t.defaultValue.String(),
			ansi.Reset,
			ansi.CursorBackward(len(t.defaultValue)+1),
		)

		prompt.WriteString(t.promptDefault)
	}

	return prompt.String()
}

func (t *Text) getPromptLines(prompt string) (int, error) {
	promptLines := 1

	writer, ok := t.Out.(*os.File)
	if !ok {
		return 0, fmt.Errorf("unable to determine prompt lines: output writer is not a file")
	}

	fd := int(writer.Fd())
	width, _, err := term.GetSize(fd)
	if err != nil {
		return 0, err
	}

	if width == 0 {
		return 0, fmt.Errorf("unable to determine prompt lines: terminal width is 0")
	}

	promptCharCnt := len(ansi.Strip(prompt))

	// If prompt is wider than terminal, calculate number of lines it is so we
	// know how many lines it occupies.
	if promptCharCnt > width {
		promptLines = (promptCharCnt / width) + 1
	}

	return promptLines, nil
}

func (t *Text) printErrorMessage(err error) {
	builder := strings.Builder{}
	// Have to manually jump to the next line, otherwise the error message will
	// be printed on the same line as the prompt.
	builder.WriteByte('\n')

	// Write error message, move/clear the line above, then reprint the prompt.
	errMsg := validationErrorMessage(err)

	errMsgLines, err := t.getPromptLines(errMsg)
	if err != nil {
		panic(err)
	}

	builder.WriteString(errMsg)
	for i := 0; i < errMsgLines-1; i++ {
		builder.WriteString(ansi.CursorUp(1))
	}

	builder.WriteString(resetLineAbove())
	builder.WriteString(t.Prompt)
	fmt.Fprint(t.Out, builder.String())
}

// printFinalPromptLine handles printing the final prompt line after successful
// input.
func (t *Text) printFinalPromptLine() {
	builder := strings.Builder{}
	// If the input is not hidden, We need to clear the line, then print the
	// prompt with the answer function applied.
	if !t.hide {
		builder.WriteString(ansi.ClearLineReset)
		t.answer.val = *t.value
		promptAnswer := t.Prompt + t.answer.Get()
		builder.WriteString(promptAnswer)
	}

	// Regardless of being hidden or not we need to move to the next line and
	// clear it in case there are any validation error messages that are still
	// visible.
	builder.WriteString("\n" + ansi.ClearLineReset)
	fmt.Fprint(t.Out, builder.String())
}

func (t *Text) handleEnter(_ grapheme.Cluster) (done bool, err error) {
	if err := t.validateFn(t.pendingValueClusterSet.String()); err != nil {
		t.printErrorMessage(err)
		t.PendingInputClusterSet = nil
		t.pendingValueClusterSet = nil
		return false, nil
	}

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
		fmt.Fprint(t.Out, t.getPrompt())
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
		fmt.Fprint(t.Out, ansi.ClearLineReset+t.Prompt)
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
