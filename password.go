package gocliselect

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/engmtcdrm/go-ansi"
	"golang.org/x/term"
)

type Password struct {
	passwordIcon EvalVal[string]
	title        EvalVal[string]
	value        *[]byte
	validate     func([]byte) error
}

func NewPassword() *Password {
	return &Password{
		passwordIcon: EvalVal[string]{val: passwordIcon, fn: nil},
		title:        EvalVal[string]{val: "", fn: nil},
		value:        nil,
		validate:     func(s []byte) error { return nil },
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

func (q *Password) PasswordIcon(s string) *Password {
	q.passwordIcon.val = s
	q.passwordIcon.fn = nil
	return q
}

func (q *Password) PasswordIconFunc(fn func() string) *Password {
	q.passwordIcon.fn = fn
	return q
}

func (q *Password) Validate(fn func([]byte) error) *Password {
	q.validate = fn
	return q
}

func (q *Password) Ask() error {
	if q.title.val == "" && q.title.fn == nil {
		return ErrNoTitle
	}

	if q.value == nil {
		return ErrNoValue
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	defer signal.Stop(sigChan)

	fmt.Print(ansi.HideCursor)
	fmt.Printf("%s%s", q.passwordIcon.Get(), q.title.Get())

	answerChan := make(chan []byte, 1)
	go func() {
		answer, _ := term.ReadPassword(int(syscall.Stdin))
		if err := q.validate(answer); err != nil {
			fmt.Printf("Error: %v\nTry again: ", err)
		}
		answerChan <- answer
	}()

	select {
	case <-sigChan:
		fmt.Print(ansi.ShowCursor)
		fmt.Println("")
		return ErrUserAborted
	case answer := <-answerChan:
		fmt.Print(ansi.ShowCursor)
		fmt.Println("")
		*q.value = answer
		return nil
	}
}
