package pardon

import (
	"fmt"
)

type Question struct {
	icon     EvalVal[string]
	title    EvalVal[string]
	value    *string
	validate func(string) error
	tui      *TuiPrompt[string]
}

func NewQuestion() *Question {
	q := &Question{
		icon:     EvalVal[string]{val: questionMarkIcon, fn: nil},
		title:    EvalVal[string]{val: "", fn: nil},
		value:    nil,
		validate: func(s string) error { return nil },
		tui:      NewTuiPrompt[string](),
	}

	q.tui.Validate(q.validate).
		DisplayInput(func(s string) string { return s }).
		AppendInput(func(s string, b byte) string { return s + string(b) }).
		RemoveLast(func(s string) string {
			if len(s) > 0 {
				return s[:len(s)-1]
			}
			return s
		}).
		ConvertInput(func(s string) string { return s })

	return q
}

func (q *Question) Title(title string) *Question {
	q.title.val = title
	q.title.fn = nil
	return q
}

func (q *Question) TitleFunc(fn func() string) *Question {
	q.title.fn = fn
	return q
}

func (q *Question) Value(value *string) *Question {
	q.value = value
	return q
}

func (q *Question) Icon(s string) *Question {
	q.icon.val = s
	q.icon.fn = nil
	return q
}

func (q *Question) IconFunc(fn func() string) *Question {
	q.icon.fn = fn
	return q
}

func (q *Question) Validate(fn func(string) error) *Question {
	q.validate = fn
	q.tui = q.tui.Validate(fn)
	return q
}

func (q *Question) Ask() error {
	question := fmt.Sprintf("%s%s ", q.icon.Get(), q.title.Get())

	return q.tui.Display(question, q.value)
}
