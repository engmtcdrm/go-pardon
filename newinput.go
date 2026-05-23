package pardon

import (
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/engmtcdrm/go-pardon/internal/keys"
)

type Input struct {
	// Hide indicates whether the input should be hidden (e.g., for password
	// input).
	Hide bool

	writer  io.Writer
	reader  *os.File
	pending []byte
	result  []rune
}

func NewInput() *Input {
	return &Input{
		Hide:   false,
		writer: os.Stdout,
		reader: os.Stdin,
	}
}

// print writes the given arguments to the terminal if [Input.Hide] is false.
func (i *Input) print(a ...any) {
	if !i.Hide {
		fmt.Fprint(i.writer, a...)
	}
}

func (i *Input) rawReadline(f *os.File) (string, error) {
	i.reset()

	for {
		var buf [8]byte
		n, err := f.Read(buf[:])
		if err != nil && err != io.EOF {
			return "", err
		}

		if n == 0 {
			if err == io.EOF {
				break
			}
			continue
		}

		i.pending = append(i.pending, buf[:n]...)

		if returnString, done, err := i.processPending(); done {
			return returnString, err
		}
	}

	i.print("\n")
	return string(i.result), nil
}

// reset clears the pending input and the result. This is called prior to
// processing new input.
func (i *Input) reset() {
	i.pending = nil
	i.result = nil
}

func (i *Input) processPending() (returnString string, done bool, err error) {
	for len(i.pending) > 0 {
		switch i.pending[0] {
		case keys.CarriageReturn, keys.Enter:
			i.print("\n")
			return string(i.result), true, nil
		case keys.CtrlC:
			return "", true, ErrUserAborted
		case keys.Backspace, keys.Delete:
			i.handleErase()
			continue
		case keys.Escape:
			if doBreak := i.handleEscapeSequence(); doBreak {
				break
			}
			continue
		}

		r, size := utf8.DecodeRune(i.pending)
		if r == utf8.RuneError && size == 1 {
			if !utf8.FullRune(i.pending) {
				break
			}

			i.pending = i.pending[1:]
			continue
		}

		i.pending = i.pending[size:]

		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			i.result = append(i.result, r)
			i.print(string(r))
		}
	}

	return "", false, nil
}

// handleErase processes a backspace or delete key press by removing the last
// character
func (i *Input) handleErase() {
	i.pending = i.pending[1:]
	if len(i.result) > 0 {
		i.result = i.result[:len(i.result)-1]
		i.print("\b \b")
	}
}

func (i *Input) handleEscapeSequence() (doBreak bool) {
	if len(i.pending) < 2 {
		return true
	}

	if i.pending[1] != keys.LeftBracket && i.pending[1] != keys.CapitalO {
		i.pending = i.pending[1:]
		return false
	}

	k := 2
	for k < len(i.pending) {
		c := i.pending[k]
		if c >= 0x40 && c <= 0x7E {
			k++
			break
		}
		k++
	}

	if k > len(i.pending) {
		return true
	}

	i.pending = i.pending[k:]
	return false
}
