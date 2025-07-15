package pardon

import (
	"fmt"
	"strings"
)

type Password struct {
	icon     EvalVal[string]
	title    EvalVal[string]
	value    *[]byte
	validate func([]byte) error
	tui      *TuiPrompt[[]byte]
}

func NewPassword() *Password {
	return &Password{
		icon:     EvalVal[string]{val: passwordIcon, fn: nil},
		title:    EvalVal[string]{val: "", fn: nil},
		value:    nil,
		validate: func(s []byte) error { return nil },
		tui: NewTuiPrompt[[]byte]().DisplayInput(
			func(input []byte) string {
				return strings.Repeat("", len(input))
			},
		),
	}
}

func (q *Password) Title(title string) *Password {
	q.title.val = title
	q.title.fn = nil
	return q
}

func (q *Password) TitleFunc(fn func() string) *Password {
	q.title.fn = fn
	return q
}

func (q *Password) Value(value *[]byte) *Password {
	q.value = value
	return q
}

func (q *Password) Icon(s string) *Password {
	q.icon.val = s
	q.icon.fn = nil
	return q
}

func (q *Password) IconFunc(fn func() string) *Password {
	q.icon.fn = fn
	return q
}

func (q *Password) Validate(fn func([]byte) error) *Password {
	q.validate = fn
	return q
}

func (p *Password) Ask() error {
	question := fmt.Sprintf("%s %s ", p.icon.Get(), p.title.Get())

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
