package new

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-ansi"
	"golang.org/x/term"
)

var (
	// ErrUserAborted is returned when the user cancels a prompt operation.
	ErrUserAborted = errors.New("user aborted")

	ErrNoPrompt = errors.New("prompt requires a prompt")
	ErrNoValue  = errors.New("value must be set")
)

type InputPrompt struct {
	answerFn   func(string) string
	validateFn func(string) error
	hidden     bool
	input      *Input
	prompt     string
	value      *string
}

// NewStringPrompt creates an InputPrompt for plaintext string input.
func NewStringPrompt() *InputPrompt {
	return &InputPrompt{
		answerFn:   func(s string) string { return s },
		validateFn: func(s string) error { return nil },
		hidden:     false,
		input:      NewInput(),
	}
}

// NewPasswordPrompt creates an InputPrompt for secure password input with masking.
func NewPasswordPrompt() *InputPrompt {
	inputPrompt := NewStringPrompt()
	inputPrompt.Hidden(true)
	return inputPrompt
}

func (p *InputPrompt) Hidden(hidden bool) *InputPrompt {
	p.hidden = hidden
	p.input.Hidden = hidden
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

func (p *InputPrompt) Display(prompt string, value *string) error {
	if prompt == "" {
		return ErrNoPrompt
	}

	if value == nil {
		return ErrNoValue
	}

	p.prompt = prompt
	p.value = value

	err := p.display()
	if err != nil {
		return err
	}

	return nil
}

func (p *InputPrompt) display() error {
	fmt.Print(p.prompt + " ")

	for {
		line, err := p.input.RawRead()
		if err != nil {
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
	// Have to manually jump to the next line if the input is hidden, otherwise
	// the error message will be printed on the same line as the prompt.
	if p.hidden {
		fmt.Println()
	}

	fmt.Print(buildValidationErrorMessage(err))
	fmt.Print(ansi.CursorUp(1) + ansi.ClearLineReset)
	fmt.Print(p.prompt + " ")
}

func (p *InputPrompt) printFinalPromptLine() {
	// If the input is not hidden, let's redraw the prompt and call the answer
	// function to look pretty.
	if !p.hidden {
		fmt.Print(ansi.CursorUp(1) + ansi.ClearLineReset)
		fmt.Print(p.prompt + " " + p.answerFn(*p.value))

	}

	fmt.Print("\n" + ansi.ClearLineReset)
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
