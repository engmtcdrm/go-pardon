package pardon

import (
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/engmtcdrm/go-pardon/internal/keys"
	"golang.org/x/term"
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

// RawRead reads input from the terminal in raw mode. It handles special keys
// like Enter, Backspace, etc., and returns the input as a slice of runes. If
// the input is interrupted (e.g., by Ctrl+C), it returns an error.
func (i *Input) RawRead() ([]rune, error) {
	reader, ok := i.Reader.(*os.File)
	if !ok {
		return nil, fmt.Errorf("unable to read input: input reader is not a file")
	}

	// MakeRaw put the terminal connected to the given file descriptor
	// into raw mode
	fd := int(reader.Fd())

	// If the reader is not connected to a terminal (e.g., during tests
	// where we use PTYs or files), fall back to rawReadline which reads
	// directly from the provided file without attempting to put the
	// descriptor into raw mode.
	if !term.IsTerminal(fd) {
		return i.rawReadline(reader)
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	defer term.Restore(fd, oldState)

	return i.rawReadline(reader)
}

// Reset clears the pending input and the result. This should be called prior
// to calling [Input.RawRead] to ensure that any previous input does not
// interfere with the new input. This is intentional to allow for easier testing
// of the [Input] struct.
func (i *Input) Reset() {
	i.pending = nil
	i.result = nil
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
func (i *Input) processPending() (returnRunes []rune, done bool, err error) {
	for len(i.pending) > 0 {
		switch i.pending[0] {
		case keys.NewLine, keys.Enter:
			return i.result, true, nil
		case keys.CtrlC:
			return nil, true, ErrUserAborted
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
				return i.result, true, nil
			}
		}
	}

	return nil, false, nil
}

func (i *Input) rawReadline(f *os.File) ([]rune, error) {
	for {
		var buf [8]byte
		n, err := f.Read(buf[:])
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n == 0 {
			if err == io.EOF {
				break
			}
			continue
		}

		i.pending = append(i.pending, buf[:n]...)

		if returnRunes, done, err := i.processPending(); done {
			return returnRunes, err
		}
	}

	i.print("\n")
	return i.result, nil
}
