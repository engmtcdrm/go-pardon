package pardon

import (
	"errors"
	"fmt"
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
		icon:       eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title:      eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
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
	inputPrompt.icon = eval[string]{val: Icons.Password, defaultFn: defaultFuncs.iconFn}
	inputPrompt.Hide(true)
	return inputPrompt
}

func (p *InputPrompt) AnswerFunc(fn func(string) string) *InputPrompt {
	if fn != nil {
		p.answerFn = fn
	}
	return p
}

func (p *InputPrompt) Hide(hide bool) *InputPrompt {
	p.hide = hide
	p.input.Hide = hide
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

// ask handles the core logic of displaying the prompt, reading user input,
// validating it, and applying the answer transformation.
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
	fmt.Print(builder.String())
}

func (p *InputPrompt) printFinalPromptLine() {
	builder := strings.Builder{}
	// If the input is not hidden, cursor will be on the next line due to user
	// pressing enter. We need to move the cursor up, clear, then print the
	// prompt with the answer function applied.
	if !p.hide {
		builder.WriteString(resetLineAbove())
		builder.WriteString(p.prompt + p.answerFn(*p.value))
	}

	// Regardless of being hidden or not we need to move to the next line and
	// clear it in case there are any validation error messages that are still
	// visible.
	builder.WriteString("\n" + ansi.ClearLineReset)
	fmt.Print(builder.String())
}

func validationErrorMessage(err error) string {
	return fmt.Sprintf("%s%s* %v%s", ansi.ClearLineReset, ansi.RedBg, err, ansi.Reset)
}

func resetLineAbove() string {
	return ansi.CursorUp(1) + ansi.ClearLineReset
}

func (p *InputPrompt) getPromptLines(prompt string) (int, error) {
	promptLines := 1

	fd := int(p.input.reader.Fd())
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
