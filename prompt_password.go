package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-pardon/tui"
)

// Password represents a password input prompt that securely collects sensitive information.
type Password struct {
	icon     eval[string]
	title    eval[string]
	value    *[]byte
	answerFn func(string) string
	tui      *tui.InputPrompt[[]byte]
}

// NewPassword creates a new Password prompt instance.
func NewPassword() *Password {
	return &Password{
		icon:  eval[string]{val: Icons.Password, defaultFn: defaultFuncs.iconFn},
		title: eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		value: nil,
		tui:   tui.NewPasswordPrompt(),
	}
}

// Title sets a static title for the password prompt.
func (p *Password) Title(title string) *Password {
	p.title.val = title
	p.title.fn = nil
	return p
}

// TitleFunc sets a dynamic title function for the password prompt.
func (p *Password) TitleFunc(fn func(string) string) *Password {
	p.title.fn = fn
	return p
}

// Value sets a default value for the password prompt.
func (p *Password) Value(value *[]byte) *Password {
	p.value = value
	return p
}

// Icon sets a static icon for the password prompt.
func (p *Password) Icon(s string) *Password {
	p.icon.val = s
	p.icon.fn = nil
	return p
}

// IconFunc sets a dynamic icon function for the password prompt.
func (p *Password) IconFunc(fn func(string) string) *Password {
	p.icon.fn = fn
	return p
}

// Validate sets a validation function for the password prompt.
func (p *Password) Validate(fn func([]byte) error) *Password {
	p.tui.Validate(fn)
	return p
}

// AnswerFunc sets a function to transform the final answer before returning.
func (p *Password) AnswerFunc(fn func(string) string) *Password {
	p.answerFn = fn
	return p
}

// setAnswerFunc configures the answer transformation priority:
// prompt-specific, global default, or identity function.
func (p *Password) setAnswerFunc() {
	if p.answerFn != nil {
		p.tui.AnswerFunc(p.answerFn)
		return
	}

	if defaultFuncs.answerFn != nil {
		p.tui.AnswerFunc(defaultFuncs.answerFn)
		return
	}

	p.tui.AnswerFunc(func(input string) string { return input })
}

// Ask displays the password prompt.
func (p *Password) Ask() error {
	question := fmt.Sprintf("%s%s ", p.icon.Get(), p.title.Get())
	p.setAnswerFunc()

	return p.tui.Display(question, p.value)
}
