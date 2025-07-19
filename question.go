package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-pardon/tui"
)

type Question struct {
	icon     evalVal[string]
	title    evalVal[string]
	value    *string
	answerFn func(string) string
	tui      *tui.InputPrompt[string]
}

func NewQuestion() *Question {
	return &Question{
		icon:  evalVal[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title: evalVal[string]{val: "", defaultFn: defaultFuncs.titleFn},
		value: nil,
		tui:   tui.NewStringPrompt(),
	}
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
	q.tui.Validate(fn)
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
