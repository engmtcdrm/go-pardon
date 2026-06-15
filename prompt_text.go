package pardon

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/keys"
	"golang.org/x/term"
)

type Text struct {
	BaseInputParser
	// Out is the output writer for the terminal, typically [os.Stdout].
	Out io.Writer

	// In is the terminal input reader.
	In *TerminalInput

	icon       eval[string]
	title      eval[string]
	answer     eval[string]
	validateFn func(string) error

	// hide indicates whether the input should be hidden (e.g., for password
	// input).
	hide bool

	prompt               string
	pendingValueRuneKeys keys.Keys
	value                *string
}

// NewPassword creates an InputPrompt for secure password input with masking.
func NewPassword(value *string) *Text {
	return &Text{
		icon:       eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title:      eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		answer:     eval[string]{val: "", defaultFn: defaultFuncs.answerFn},
		validateFn: func(s string) error { return nil },
		In:         NewTerminalInput(),
		Out:        os.Stdout,
		hide:       true,
		value:      value,
	}
}

// NewQuestion creates a new InputPrompt for text input with a question mark
// icon.
func NewQuestion(value *string) *Text {
	return &Text{
		icon:       eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title:      eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		answer:     eval[string]{val: "", defaultFn: defaultFuncs.answerFn},
		validateFn: func(s string) error { return nil },
		In:         NewTerminalInput(),
		Out:        os.Stdout,
		value:      value,
	}
}

// AnswerFunc sets a function to format the final answer being displayed.
func (t *Text) AnswerFunc(fn func(string) string) *Text {
	if fn != nil {
		t.answer.fn = fn
	}
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

// Icon sets the prompt icon.
func (t *Text) Icon(s string) *Text {
	t.icon.val = s
	t.icon.fn = nil
	return t
}

// IconFunc sets a dynamic icon function.
func (t *Text) IconFunc(fn func(string) string) *Text {
	t.icon.fn = fn
	return t
}

// Title sets the prompt text.
func (t *Text) Title(title string) *Text {
	t.title.val = title
	t.title.fn = nil
	return t
}

// TitleFunc sets a dynamic title function.
func (t *Text) TitleFunc(fn func(string) string) *Text {
	t.title.fn = fn
	return t
}

// ValidateFunc sets a validation function for the prompt input.
func (t *Text) ValidateFunc(fn func(string) error) *Text {
	if fn != nil {
		t.validateFn = fn
	}
	return t
}

// Value sets a default input value.
func (t *Text) Value(value *string) *Text {
	t.value = value
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

	t.prompt = fmt.Sprintf("%s%s ", t.icon.Get(), t.title.Get())
	fmt.Fprint(t.Out, t.prompt)

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
	builder.WriteString(t.prompt)
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
		promptAnswer := t.prompt + t.answer.Get()
		builder.WriteString(promptAnswer)
	}

	// Regardless of being hidden or not we need to move to the next line and
	// clear it in case there are any validation error messages that are still
	// visible.
	builder.WriteString("\n" + ansi.ClearLineReset)
	fmt.Fprint(t.Out, builder.String())
}

func (t *Text) handleEnter(_ keys.Key) (done bool, err error) {
	if err := t.validateFn(t.pendingValueRuneKeys.String()); err != nil {
		t.printErrorMessage(err)
		t.pendingInputRuneKeys = nil
		t.pendingValueRuneKeys = nil
		return false, nil
	}

	*t.value = t.pendingValueRuneKeys.String()

	t.printFinalPromptLine()
	return true, nil
}

func (t *Text) handleDelete(_ keys.Key) (done bool, err error) {
	if len(t.pendingValueRuneKeys) > 0 {
		t.pendingValueRuneKeys = t.pendingValueRuneKeys[:len(t.pendingValueRuneKeys)-1]
		t.printInput("\b \b")
	}

	t.pendingInputRuneKeys = t.pendingInputRuneKeys[1:]
	return false, nil
}

// Currently calls [Text.handleEnter].
func (t *Text) handleNewline(r keys.Key) (done bool, err error) {
	return t.handleEnter(r)
}

func (t *Text) processInput(input []byte) (done bool, err error) {
	if len(input) == 0 {
		return false, nil
	}

	t.pendingInput = append(t.pendingInput, input...)

	if needMoreInput := t.parseInputToRuneKeys(); needMoreInput {
		return false, nil
	}

	// Reset pending escape sequence before processing new input.
	t.pendingEscSequence = keys.Keys{}

	for len(t.pendingInputRuneKeys) > 0 {
		r := t.pendingInputRuneKeys[0]
		switch {
		case equal(r, keys.CtrlC):
			return true, ErrUserAborted
		case equal(r, keys.Enter):
			return t.handleEnter(r)
		case equal(r, keys.Newline):
			return t.handleNewline(r)
		case equal(r, keys.Delete), equal(r, keys.Backspace):
			return t.handleDelete(r)
		case equal(r, keys.Escape):
			doContinue, err := t.processEscapeSequence(r, nil)
			if !doContinue {
				return false, err
			}
			continue
		}

		t.pendingValueRuneKeys = append(t.pendingValueRuneKeys, r)
		t.pendingInputRuneKeys = t.pendingInputRuneKeys[1:]
		t.printInput(string(r))
	}

	t.pendingInputRuneKeys = nil

	return false, nil
}

func (t *Text) printInput(a ...any) {
	if !t.hide {
		fmt.Fprint(t.Out, a...)
	}
}
