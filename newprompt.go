package pardon

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"golang.org/x/term"
)

var (
	ErrNoPrompt = errors.New("prompt requires a prompt")
)

type InputPrompt struct {
	icon       eval[string]
	title      eval[string]
	answerFn   func(string) string
	validateFn func(string) error

	terminal *Terminal
	prompt   string
	value    *string
}

// NewStringPrompt creates an InputPrompt for plaintext string input.
func NewStringPrompt(value *string) *InputPrompt {
	return &InputPrompt{
		icon:       eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title:      eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		validateFn: func(s string) error { return nil },
		terminal:   NewTerminal(),
		value:      value,
	}
}

// NewPasswordPrompt creates an InputPrompt for secure password input with masking.
func NewPasswordPrompt(value *string) *InputPrompt {
	return &InputPrompt{
		icon:       eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title:      eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		validateFn: func(s string) error { return nil },
		terminal:   NewHiddenTerminal(),
		value:      value,
	}
}

func (p *InputPrompt) AnswerFunc(fn func(string) string) *InputPrompt {
	if fn != nil {
		p.answerFn = fn
	}
	return p
}

// Hide sets whether the input should be hidden (e.g., for password input).
func (p *InputPrompt) Hide(hide bool) *InputPrompt {
	p.terminal.Hide = hide
	return p
}

// Icon sets the prompt icon.
func (p *InputPrompt) Icon(s string) *InputPrompt {
	p.icon.val = s
	p.icon.fn = nil
	return p
}

// IconFunc sets a dynamic icon function.
func (p *InputPrompt) IconFunc(fn func(string) string) *InputPrompt {
	p.icon.fn = fn
	return p
}

// Title sets the prompt text.
func (p *InputPrompt) Title(title string) *InputPrompt {
	p.title.val = title
	return p
}

// TitleFunc sets a dynamic title function.
func (p *InputPrompt) TitleFunc(fn func(string) string) *InputPrompt {
	p.title.fn = fn
	return p
}

// Value sets a default input value.
func (p *InputPrompt) Value(value *string) *InputPrompt {
	p.value = value
	return p
}

// ValidateFunc sets a validation function for the prompt input.
func (p *InputPrompt) ValidateFunc(fn func(string) error) *InputPrompt {
	if fn != nil {
		p.validateFn = fn
	}
	return p
}

// Ask displays the prompt and waits for input.
func (p *InputPrompt) Ask() error {
	if p.title.val == "" && p.title.fn == nil {
		return ErrNoTitle
	}

	if p.value == nil {
		return ErrNoValue
	}

	p.prompt = fmt.Sprintf("%s%s ", p.icon.Get(), p.title.Get())

	if err := p.ask(); err != nil {
		return err
	}

	return nil
}

func (p *InputPrompt) callAnswerFunc(s string) string {
	if p.answerFn != nil {
		return p.answerFn(s)
	}

	if defaultFuncs.answerFn != nil {
		return defaultFuncs.answerFn(s)
	}

	return s
}

// ask handles the core logic of displaying the prompt, reading user input,
// validating it, and applying the answer transformation.
func (p *InputPrompt) ask() error {
	fmt.Fprint(p.terminal.Out, p.prompt)

	for {
		p.terminal.Reset()
		line, err := p.terminal.RawRead()
		if err != nil {
			if errors.Is(err, ErrUserAborted) {
				fmt.Fprint(p.terminal.Out, ansi.ClearLineReset+p.prompt)
				return err
			}
			return err
		}

		pendingValue := strings.TrimSpace(string(line))
		if err := p.validateFn(pendingValue); err != nil {
			p.printErrorMessage(err)
			continue
		}

		*p.value = pendingValue
		break
	}

	p.printFinalPromptLine()

	return nil
}

func (p *InputPrompt) printErrorMessage(err error) {
	builder := strings.Builder{}
	// Have to manually jump to the next line, otherwise the error message will
	// be printed on the same line as the prompt.
	builder.WriteByte('\n')

	// Write error message, move/clear the line above, then reprint the prompt.
	errMsg := validationErrorMessage(err)

	errMsgLines, err := p.getPromptLines(errMsg)
	if err != nil {
		panic(err)
	}

	builder.WriteString(errMsg)
	for i := 0; i < errMsgLines-1; i++ {
		builder.WriteString(ansi.CursorUp(1))
	}

	builder.WriteString(resetLineAbove())
	builder.WriteString(p.prompt)
	fmt.Fprint(p.terminal.Out, builder.String())
}

// printFinalPromptLine handles printing the final prompt line after successful
// input.
func (p *InputPrompt) printFinalPromptLine() {
	builder := strings.Builder{}
	// If the input is not hidden, We need to clear the line, then print the
	// prompt with the answer function applied.
	if !p.terminal.Hide {
		builder.WriteString(ansi.ClearLineReset)
		builder.WriteString(p.prompt + p.callAnswerFunc(*p.value))
	}

	// Regardless of being hidden or not we need to move to the next line and
	// clear it in case there are any validation error messages that are still
	// visible.
	builder.WriteString("\n" + ansi.ClearLineReset)
	fmt.Fprint(p.terminal.Out, builder.String())
}

func validationErrorMessage(err error) string {
	return fmt.Sprintf("%s%s* %v%s", ansi.ClearLineReset, ansi.RedBg, err, ansi.Reset)
}

func resetLineAbove() string {
	return ansi.CursorUp(1) + ansi.ClearLineReset
}

func (p *InputPrompt) getPromptLines(prompt string) (int, error) {
	promptLines := 1

	reader, ok := p.terminal.In.(*os.File)
	if !ok {
		return 0, fmt.Errorf("unable to determine prompt lines: input reader is not a file")
	}

	fd := int(reader.Fd())
	width, _, err := term.GetSize(fd)
	if err != nil {
		return 0, err
	}

	promptCharCnt := len(ansi.Strip(prompt))

	// If prompt is wider than terminal calculate number of lines it will take
	if promptCharCnt > width {
		promptLines = (promptCharCnt / width) + 1
	}

	return promptLines, nil
}
