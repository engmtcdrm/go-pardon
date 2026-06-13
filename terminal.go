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

type Terminal struct {
	// Hide indicates whether the input should be hidden (e.g., for password
	// input).
	Hide bool

	// Confirm indicates whether the input should be treated as a confirmation.
	Confirm bool

	// Out is the output writer for the terminal, typically [os.Stdout].
	Out io.Writer
	// In is the input reader for the terminal, typically [os.Stdin].
	In io.Reader

	// pending holds the bytes that have been read but not yet processed.
	pending []byte

	// result holds the runes that have been processed and are part of the final
	// input.
	result []rune

	CustomHandler func(t *Terminal, r rune) (done bool)
}

// NewTerminal creates a new Terminal instance with default settings for regular
// input. Output will go to [os.Stdout] and input will be read from [os.Stdin].
func NewTerminal() *Terminal {
	return &Terminal{
		Hide: false,
		Out:  os.Stdout,
		In:   os.Stdin,
	}
}

// NewHiddenTerminal creates a new Terminal instance configured for hidden input,
// such as for password prompts. Output will go to [os.Stdout] and input will be
// read from [os.Stdin].
func NewHiddenTerminal() *Terminal {
	input := NewTerminal()
	input.Hide = true

	return input
}

// GetTerminalHeight returns the height of the terminal in rows. If the terminal
// size cannot be determined, it returns a default height of 25 rows.
func (t *Terminal) GetTerminalHeight() int {
	termHeight := 25 // Default height

	f, ok := t.Out.(*os.File)
	if !ok {
		return termHeight
	}

	if _, height, err := term.GetSize(int(f.Fd())); err == nil {
		termHeight = height
	}

	return termHeight
}

// Print writes the given arguments to the terminal output.
func (t *Terminal) Print(a ...any) {
	fmt.Fprint(t.Out, a...)
}

// Printf formats according to a format specifier and writes to the terminal
// output.
func (t *Terminal) Printf(format string, a ...any) {
	fmt.Fprintf(t.Out, format, a...)
}

// Println writes the given arguments to the terminal output, followed by a
// newline.
func (t *Terminal) Println(a ...any) {
	fmt.Fprintln(t.Out, a...)
}

// RawRead reads input from the terminal in raw mode. It handles special keys
// like Enter, Backspace, etc., and returns the input as a slice of runes. If
// the input is interrupted (e.g., by Ctrl+C), it returns an error.
//
// [Terminal.Reset] must be called before invoking this function to ensure that
// any previous input does not interfere with the new input.
func (t *Terminal) RawRead() ([]rune, error) {
	inputFile, restoreTerminal, err := t.setTerminalToRawMode()
	if err != nil {
		return nil, err
	}
	defer restoreTerminal()

	return t.rawReadline(inputFile)
}

// RawRead2 reads input from the terminal in raw mode and returns the raw bytes.
//
// The caller is responsible for processing the bytes and handling special keys.
// As well as wrapping this call in a for loop to continue reading until the
// desired input is complete.
func (t *Terminal) RawRead2() ([]byte, error) {
	inputFile, restoreTerminal, err := t.setTerminalToRawMode()
	if err != nil {
		return nil, err
	}
	defer restoreTerminal()

	var buf [8]byte
	n, err := inputFile.Read(buf[:])
	if err != nil && err != io.EOF {
		return nil, err
	}

	if n == 0 {
		if err == io.EOF {
			return nil, nil
		}
	}

	return buf[:n], nil
}

// Reset clears the pending input and the result. This should be called prior
// to calling [Terminal.RawRead] to ensure that any previous input does not
// interfere with the new input. This is intentional to allow for easier testing
// of the [Terminal] struct.
func (t *Terminal) Reset() {
	t.pending = nil
	t.result = nil
}

// handleErase processes a backspace or delete key press by removing the last
// character
func (t *Terminal) handleErase() {
	t.pending = t.pending[1:]
	if len(t.result) > 0 {
		t.result = t.result[:len(t.result)-1]
		t.printInput("\b \b")
	}
}

// handleEscapeSequence processes an escape sequence starting with the escape key.
func (t *Terminal) handleEscapeSequence() (doBreak bool) {
	if len(t.pending) < 2 {
		return true
	}

	if !validateEscapeSequence(t.pending[1]) {
		t.pending = t.pending[1:]
		return false
	}

	k := 2
	foundFinal := false
	for k < len(t.pending) {
		c := t.pending[k]
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

	t.pending = t.pending[k:]
	return false
}

// printInput writes the given arguments to the terminal if [Terminal.Hide] is
// false.
func (t *Terminal) printInput(a ...any) {
	if !t.Hide {
		t.Print(a...)
	}
}

// processPending processes the pending input bytes and updates the result.
func (t *Terminal) processPending() (returnRunes []rune, done bool, err error) {
	for len(t.pending) > 0 {
		switch t.pending[0] {
		case keys.NewLine, keys.Enter:
			t.result = append(t.result, rune(t.pending[0]))
			return t.result, true, nil
		case keys.CtrlC:
			return nil, true, ErrUserAborted
		case keys.Delete, keys.Backspace:
			t.handleErase()
			continue
		case keys.Escape:
			if doBreak := t.handleEscapeSequence(); doBreak {
				break
			}
			continue
		}

		r, size := utf8.DecodeRune(t.pending)
		if r == utf8.RuneError && size == 1 {
			if !utf8.FullRune(t.pending) {
				break
			}

			t.pending = t.pending[1:]
			continue
		}

		t.pending = t.pending[size:]

		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			t.result = append(t.result, r)

			if t.CustomHandler != nil {
				if done := t.CustomHandler(t, r); done {
					return t.result, true, nil
				}
			}

			if !t.Confirm {
				t.printInput(string(r))
			} else {
				return t.result, true, nil
			}
		}
	}

	return nil, false, nil
}

func (t *Terminal) rawReadline(f *os.File) ([]rune, error) {
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

		// idx := bytes.IndexAny(buf[:n], "\n\r")

		// If we found a newline/carriage return anywhere in the buffer,
		// move the file offset back so the next read will return any bytes
		// that follow the newline. The number of bytes to move back is
		// (idx+1 - n) which is <= 0.
		// if idx >= 0 {
		// 	_, _ = f.Seek(int64(idx+1-n), io.SeekCurrent)
		// }

		t.pending = append(t.pending, buf[:n]...)

		if returnRunes, done, err := t.processPending(); done {
			return returnRunes, err
		}
	}

	t.printInput("\n")
	return t.result, nil
}

// setTerminalToRawMode attempts to put the terminal into raw mode and returns
// the input file, file descriptor, old terminal state (if raw mode was set),
// and any error encountered.  If the input is not a terminal, it returns the
// file and a no-op restore function without error. The caller should defer the
// restore function to ensure that the terminal state is properly restored after
// raw input is processed.
func (t *Terminal) setTerminalToRawMode() (inputFile *os.File, restoreTerminal func(), err error) {
	inputFile, ok := t.In.(*os.File)
	if !ok {
		return nil, func() {
			// No cleanup needed since we didn't set raw mode
		}, fmt.Errorf("unable to read input: input reader is not a file")
	}

	// MakeRaw put the terminal connected to the given file descriptor
	// into raw mode
	fd := int(inputFile.Fd())

	// If the reader is not connected to a terminal (e.g., during tests
	// where we use PTYs or files), don't attempt to set raw mode and just
	// return the file and descriptor.
	if !term.IsTerminal(fd) {
		return inputFile, func() {
			// No cleanup needed since we didn't set raw mode
		}, nil
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return inputFile, func() {
			// No cleanup needed since we didn't set raw mode
		}, err
	}
	return inputFile, func() {
		term.Restore(fd, oldState)
	}, nil
}
