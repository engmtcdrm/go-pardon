package pardon

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/keys"
	"golang.org/x/term"
)

type Text struct {
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

	prompt       string
	pendingValue []byte
	value        *string
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

func (t *Text) processInput(input []byte) (done bool, err error) {
	if len(input) == 0 {
		return false, nil
	}

	switch {
	case bytes.Equal(input, keys.CtrlC):
		return true, ErrUserAborted
	case bytes.Equal(input, keys.Enter), bytes.Equal(input, keys.Newline):
		if err := t.validateFn(string(t.pendingValue)); err != nil {
			t.printErrorMessage(err)
			return false, nil
		}

		*t.value = string(t.pendingValue)

		t.printFinalPromptLine()
		return true, nil
	case bytes.Equal(input, keys.Delete), bytes.Equal(input, keys.Backspace):
		if len(t.pendingValue) > 0 {
			// Remove last UTF-8 rune (safe for multi-byte characters).
			_, size := utf8.DecodeLastRune(t.pendingValue)
			if size <= 0 {
				size = 1
			}
			t.pendingValue = t.pendingValue[:len(t.pendingValue)-size]
			t.printInput("\b \b")
		}
		return false, nil
	}

	t.pendingValue = append(t.pendingValue, input...)
	t.printInput(string(input))

	return false, nil
}

func (t *Text) printInput(a ...any) {
	if !t.hide {
		fmt.Fprint(t.Out, a...)
	}
}
