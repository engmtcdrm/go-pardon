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

	// Confirm indicates whether the input should be treated as a confirmation.
	Confirm bool

	Writer io.Writer
	Reader io.Reader

	// pending holds the bytes that have been read but not yet processed.
	pending []byte

	// result holds the runes that have been processed and are part of the final
	// input.
	result []rune
}

func NewInput() *Input {
	return &Input{
		Hide:   false,
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
}

func NewHiddenInput() *Input {
	input := NewInput()
	input.Hide = true

	return input
}

func NewConfirmInput() *Input {
	input := NewInput()
	input.Confirm = true

	return input
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

// handleEscapeSequence processes an escape sequence starting with the escape key.
func (i *Input) handleEscapeSequence() (doBreak bool) {
	if len(i.pending) < 2 {
		return true
	}

	if !validateEscapeSequence(i.pending[1]) {
		i.pending = i.pending[1:]
		return false
	}

	k := 2
	foundFinal := false
	for k < len(i.pending) {
		c := i.pending[k]
		// The final byte of an escape sequence is in the range 0x40 to 0x7E.
		if c >= 0x40 && c <= 0x7E {
			k++
			foundFinal = true
			break
		}
		k++
	}

	// If we didn't find a final byte yet, the sequence is incomplete — wait
	// for more bytes by signalling the caller to break processing.
	if !foundFinal {
		return true
	}

	i.pending = i.pending[k:]
	return false
}

// print writes the given arguments to the terminal if [Input.Hide] is false.
func (i *Input) print(a ...any) {
	if !i.Hide {
		fmt.Fprint(i.Writer, a...)
	}
}

// processPending processes the pending input bytes and updates the result.
func (i *Input) processPending() (returnString string, done bool, err error) {
	for len(i.pending) > 0 {
		switch i.pending[0] {
		case keys.NewLine, keys.Enter:
			return string(i.result), true, nil
		case keys.CtrlC:
			return "", true, ErrUserAborted
		case keys.Delete, keys.Backspace:
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
			if !i.Confirm {
				i.print(string(r))
			} else {
				return string(i.result), true, nil
			}
		}
	}

	return "", false, nil
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
