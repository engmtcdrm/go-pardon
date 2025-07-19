package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-pardon/tui"
)

type Password struct {
	icon     evalVal[string]
	title    evalVal[string]
	value    *[]byte
	answerFn func(string) string
	tui      *tui.InputPrompt[[]byte]
}

func NewPassword() *Password {
	return &Password{
		icon:  evalVal[string]{val: Icons.Password, defaultFn: defaultFuncs.iconFn},
		title: evalVal[string]{val: "", defaultFn: defaultFuncs.titleFn},
		value: nil,
		tui:   tui.NewPasswordPrompt(),
	}
}

func (p *Password) Title(title string) *Password {
	p.title.val = title
	p.title.fn = nil
	return p
}

func (p *Password) TitleFunc(fn func(string) string) *Password {
	p.title.fn = fn
	return p
}

func (p *Password) Value(value *[]byte) *Password {
	p.value = value
	return p
}

func (p *Password) Icon(s string) *Password {
	p.icon.val = s
	p.icon.fn = nil
	return p
}

func (p *Password) IconFunc(fn func(string) string) *Password {
	p.icon.fn = fn
	return p
}

func (p *Password) Validate(fn func([]byte) error) *Password {
	p.tui.Validate(fn)
	return p
}

func (p *Password) AnswerFunc(fn func(string) string) *Password {
	p.answerFn = fn
	return p
}

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

func (p *Password) Ask() error {
	question := fmt.Sprintf("%s%s ", p.icon.Get(), p.title.Get())
	p.setAnswerFunc()

	return p.tui.Display(question, p.value)
}
