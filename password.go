package pardon

import (
	"fmt"
	"strings"
)

type Password struct {
	icon     evalVal[string]
	title    evalVal[string]
	value    *[]byte
	answerFn func(string) string
	tui      *tuiPrompt[[]byte]
}

func NewPassword() *Password {
	p := &Password{
		icon:  evalVal[string]{val: Icons.Password, fn: nil, defaultFn: defaultFuncs.iconFn},
		title: evalVal[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		value: nil,
		tui:   newTuiPrompt[[]byte](),
	}

	p.tui.Validate(func(s []byte) error { return nil }).
		DisplayInput(func(input []byte) string { return strings.Repeat("", len(input)) }).
		AppendInput(func(b []byte, c byte) []byte { return append(b, c) }).
		RemoveLast(func(b []byte) []byte {
			if len(b) > 0 {
				return b[:len(b)-1]
			}
			return b
		}).
		ConvertInput(func(b []byte) []byte { return b })

	return p
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
