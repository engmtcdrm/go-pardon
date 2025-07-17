package pardon

import (
	"fmt"
	"strings"
)

type Password struct {
	icon     evalVal[string]
	title    evalVal[string]
	value    *[]byte
	validate func([]byte) error
	tui      *tuiPrompt[[]byte]
}

func NewPassword() *Password {
	return &Password{
		icon:     evalVal[string]{val: Icons.Password, fn: nil},
		title:    evalVal[string]{val: "", fn: nil},
		value:    nil,
		validate: func(s []byte) error { return nil },
		tui: newTuiPrompt[[]byte]().DisplayInput(
			func(input []byte) string {
				return strings.Repeat("", len(input))
			},
		),
	}
}

func (p *Password) Title(title string) *Password {
	p.title.val = title
	p.title.fn = nil
	return p
}

func (p *Password) TitleFunc(fn func() string) *Password {
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

func (p *Password) IconFunc(fn func() string) *Password {
	p.icon.fn = fn
	return p
}

func (p *Password) Validate(fn func([]byte) error) *Password {
	p.validate = fn
	return p
}

func (p *Password) Ask() error {
	question := fmt.Sprintf("%s%s ", p.icon.Get(), p.title.Get())

	p.tui = p.tui.AppendInput(func(b []byte, c byte) []byte { return append(b, c) }).
		Validate(p.validate).
		RemoveLast(func(b []byte) []byte {
			if len(b) > 0 {
				return b[:len(b)-1]
			}
			return b
		}).
		ConvertInput(func(b []byte) []byte { return b })

	return p.tui.Display(question, p.value)
}
