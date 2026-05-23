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

	// Hide indicates whether the input should be hidden (e.g., for password
	// input).
	hide bool

	input  *Input
	prompt string
	value  *string
}

// NewStringPrompt creates an InputPrompt for plaintext string input.
func NewStringPrompt(value *string) *InputPrompt {
	return &InputPrompt{
		answerFn:   func(s string) string { return s },
		validateFn: func(s string) error { return nil },
		hide:       false,
		input:      NewInput(),
		value:      value,
	}
}

// NewPasswordPrompt creates an InputPrompt for secure password input with masking.
func NewPasswordPrompt(value *string) *InputPrompt {
	inputPrompt := NewStringPrompt(value)
	inputPrompt.Hide(true)
	return inputPrompt
}

func (p *InputPrompt) Hide(hide bool) *InputPrompt {
	p.hide = hide
	p.input.Hide = hide
	return p
}

// Title sets the question text.
func (p *InputPrompt) Title(title string) *InputPrompt {
	p.title.val = title
	return p
}

// TitleFunc sets a dynamic title function.
func (p *InputPrompt) TitleFunc(fn func(string) string) *InputPrompt {
	p.title.fn = fn
	return p
}

func (p *InputPrompt) Validate(fn func(string) error) *InputPrompt {
	if fn != nil {
		p.validateFn = fn
	}
	return p
}

func (p *InputPrompt) AnswerFunc(fn func(string) string) *InputPrompt {
	if fn != nil {
		p.answerFn = fn
	}
	return p
}

// Value sets a default input value.
func (p *InputPrompt) Value(value *string) *InputPrompt {
	p.value = value
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

func (p *InputPrompt) Ask() error {
	if p.title.val == "" && p.title.fn == nil {
		return ErrNoTitle
	}

	if p.value == nil {
		return ErrNoValue
	}

	p.prompt = fmt.Sprintf("%s%s ", p.icon.Get(), p.title.Get())

	err := p.ask()
	if err != nil {
		return err
	}

	return nil
}

func (p *InputPrompt) ask() error {
	fmt.Print(p.prompt)

	for {
		line, err := p.input.RawRead()
		if err != nil {
			if errors.Is(err, ErrUserAborted) {
				fmt.Print(ansi.ClearLineReset + p.prompt)
				return err
			}
			return err
		}

		pendingValue := strings.TrimSpace(line)
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
	// Have to manually jump to the next line if the input is hidden, otherwise
	// the error message will be printed on the same line as the prompt.
	if p.hide {
		builder.WriteByte('\n')
	}

	builder.WriteString(buildValidationErrorMessage(err))
	builder.WriteString(ansi.CursorUp(1) + ansi.ClearLineReset)
	builder.WriteString(p.prompt)
	fmt.Print(builder.String())
}

func (p *InputPrompt) printFinalPromptLine() {
	builder := strings.Builder{}
	// If the input is not hidden, let's redraw the prompt and call the answer
	// function to look pretty.
	if !p.hide {
		builder.WriteString(ansi.CursorUp(1) + ansi.ClearLineReset)
		builder.WriteString(p.prompt + p.answerFn(*p.value))
	}

	builder.WriteString("\n" + ansi.ClearLineReset)
	fmt.Print(builder.String())
}

func buildValidationErrorMessage(err error) string {
	return fmt.Sprintf("%s%s* %s%s", ansi.ClearLineReset, ansi.Red, err.Error(), ansi.Reset)
}

func getPromptLines(prompt string) (int, error) {
	promptLines := 1

	fd := int(os.Stdin.Fd())
	width, _, err := term.GetSize(fd)
	if err != nil {
		return 0, err
	}

	promptCharCnt := len(ansi.StripCodes(prompt))

	// If prompt is wider than terminal calculate number of lines it will take
	if promptCharCnt > width {
		promptLines = (promptCharCnt / width) + 1
	}

	return promptLines, nil
}
