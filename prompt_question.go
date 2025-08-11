package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-pardon/tui"
)

// Question represents a text input prompt for user questions.
type Question struct {
	icon     eval[string]
	title    eval[string]
	value    *string
	answerFn func(string) string
	tui      *tui.InputPrompt[string]
}

// NewQuestion creates a new Question prompt instance.
func NewQuestion() *Question {
	return &Question{
		icon:  eval[string]{val: Icons.QuestionMark, defaultFn: defaultFuncs.iconFn},
		title: eval[string]{val: "", defaultFn: defaultFuncs.titleFn},
		value: nil,
		tui:   tui.NewStringPrompt(),
	}
}

// Title sets the question text.
func (q *Question) Title(title string) *Question {
	q.title.val = title
	q.title.fn = nil
	return q
}

// TitleFunc sets a dynamic title function.
func (q *Question) TitleFunc(fn func(string) string) *Question {
	q.title.fn = fn
	return q
}

// Value sets a default input value.
func (q *Question) Value(value *string) *Question {
	q.value = value
	return q
}

// Icon sets the prompt icon.
func (q *Question) Icon(s string) *Question {
	q.icon.val = s
	q.icon.fn = nil
	return q
}

// IconFunc sets a dynamic icon function.
func (q *Question) IconFunc(fn func(string) string) *Question {
	q.icon.fn = fn
	return q
}

// AnswerFunc sets a function to transform the final answer.
func (q *Question) AnswerFunc(fn func(string) string) *Question {
	q.answerFn = fn
	return q
}

// Validate sets input validation.
func (q *Question) Validate(fn func(string) error) *Question {
	q.tui.Validate(fn)
	return q
}

// setAnswerFunc configures the answer transformation priority:
// prompt-specific, global default, or identity function.
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

// Ask displays the question prompt and waits for input.
func (q *Question) Ask() error {
	question := fmt.Sprintf("%s%s ", q.icon.Get(), q.title.Get())
	q.setAnswerFunc()

	return q.tui.Display(question, q.value)
}
