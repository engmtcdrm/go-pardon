package pardon

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type TerminalInput struct {
	// Reader is the input reader for the terminal, typically [os.Stdin].
	Reader io.Reader
}

// NewTerminalInput creates a new [TerminalInput] instance with input read from
// [os.Stdin].
func NewTerminalInput() TerminalInput {
	return TerminalInput{Reader: os.Stdin}
}

// RawRead reads input from the terminal in raw mode and returns the raw bytes.
//
// The caller is responsible for processing the bytes and handling special keys.
// As well as wrapping this call in a for loop to continue reading until the
// desired input is complete.
func (t TerminalInput) RawRead() ([]byte, error) {
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
func (t TerminalInput) setTerminalToRawMode() (inputFile *os.File, restoreTerminal func(), err error) {
	inputFile, ok := t.Reader.(*os.File)
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
