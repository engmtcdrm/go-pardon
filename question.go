package gocliselect

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Question struct {
	questionMark EvalVal[string]
	title        EvalVal[string]
	value        *string
	validate     func(string) error
}

func NewQuestion() *Question {
	return &Question{
		questionMark: EvalVal[string]{val: questionMark, fn: nil},
		title:        EvalVal[string]{val: "", fn: nil},
		value:        nil,
		validate:     func(s string) error { return nil },
	}
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

func (q *Question) QuestionMark(s string) *Question {
	q.questionMark.val = s
	q.questionMark.fn = nil
	return q
}

func (q *Question) QuestionMarkFunc(fn func() string) *Question {
	q.questionMark.fn = fn
	return q
}

func (q *Question) Validate(fn func(string) error) *Question {
	q.validate = fn
	return q
}

func (q *Question) Ask() error {
	if q.title.val == "" && q.title.fn == nil {
		return ErrNoTitle
	}

	if q.value == nil {
		return ErrNoValue
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	defer signal.Stop(sigChan)

	fmt.Printf("%s %s ", q.questionMark.Get(), q.title.Get())

	answerChan := make(chan string, 1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if err := q.validate(answer); err != nil {
				fmt.Printf("Error: %v\nTry again: ", err)
				continue
			}
			answerChan <- answer
			break
		}
	}()

	select {
	case <-sigChan:
		fmt.Println("")
		return ErrUserAborted
	case answer := <-answerChan:
		*q.value = answer
		return nil
	}
}
