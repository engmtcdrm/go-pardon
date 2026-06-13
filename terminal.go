package pardon

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	// Hide indicates whether the input should be hidden (e.g., for password
	// input).
	Hide bool

	// Out is the output writer for the terminal, typically [os.Stdout].
	Out io.Writer
	// In is the input reader for the terminal, typically [os.Stdin].
	In io.Reader
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

// PrintInput writes the given arguments to th>e terminal if [Terminal.Hide] is
// false.
func (t *Terminal) PrintInput(a ...any) {
	if !t.Hide {
		t.Print(a...)
	}
}

// RawRead reads input from the terminal in raw mode and returns the raw bytes.
//
// The caller is responsible for processing the bytes and handling special keys.
// As well as wrapping this call in a for loop to continue reading until the
// desired input is complete.
func (t *Terminal) RawRead() ([]byte, error) {
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
