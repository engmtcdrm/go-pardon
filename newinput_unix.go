//go:build linux || darwin || freebsd
// +build linux darwin freebsd

package pardon

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// RawRead reads input from the terminal in raw mode. It handles special keys
// like Enter, Backspace, etc., and returns the input as a string. If the input
// is interrupted (e.g., by Ctrl+C), it returns an error.
func (i *Input) RawRead() (string, error) {
	reader, ok := i.Reader.(*os.File)
	if !ok {
		return "", fmt.Errorf("unable to read input: input reader is not a file")
	}

	// MakeRaw put the terminal connected to the given file descriptor
	// into raw mode
	fd := int(reader.Fd())
	if !term.IsTerminal(fd) {
		return "", fmt.Errorf("file descriptor %d is not a terminal", fd)
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	return i.rawReadline(reader)
}
