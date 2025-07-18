package pardon

import (
	"fmt"
)

type Question struct {
	icon     evalVal[string]
	title    evalVal[string]
	value    *string
	answerFn func(string) string
	tui      *tuiPrompt[string]
}

func NewQuestion() *Question {
	q := &Question{
		icon:  evalVal[string]{val: Icons.QuestionMark, fn: nil, defaultFn: defaultFuncs.iconFn},
		title: evalVal[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		value: nil,
		tui:   newTuiPrompt[string](),
	}

	q.tui.Validate(func(s string) error { return nil }).
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

func (q *Question) TitleFunc(fn func(string) string) *Question {
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

func (q *Question) IconFunc(fn func(string) string) *Question {
	q.icon.fn = fn
	return q
}

func (q *Question) AnswerFunc(fn func(string) string) *Question {
	q.answerFn = fn
	return q
}

func (q *Question) Validate(fn func(string) error) *Question {
	q.tui = q.tui.Validate(fn)
	return q
}

func (q *Question) setAnswerFunc() {
	if q.answerFn != nil {
		q.tui.AnswerFunc(q.answerFn)
		return
	}

	if defaultFuncs.answerFn != nil {
		q.tui.AnswerFunc(defaultFuncs.answerFn)
		return
	}

	q.tui.AnswerFunc(func(input string) string { return input })
}

func (q *Question) Ask() error {
	question := fmt.Sprintf("%s%s ", q.icon.Get(), q.title.Get())
	q.setAnswerFunc()

	return q.tui.Display(question, q.value)
}
