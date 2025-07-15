package gocliselect

import (
	"os"

	"golang.org/x/term"
)

const (
	questionMarkIcon = "[?] "
	passwordIcon     = "🔒 "
)

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getInput will read raw input from the terminal
// It returns the raw ASCII value inputted
func getInput() byte {
	// Use stdin file descriptor for cross-platform compatibility
	fd := int(os.Stdin.Fd())

	// Save the original terminal state
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// Fallback if raw mode fails - read normally
		readBytes := make([]byte, 1)
		_, readErr := os.Stdin.Read(readBytes)
		if readErr != nil {
			return 0
		}
		return readBytes[0]
	}
	defer term.Restore(fd, oldState)

	// Read input
	readBytes := make([]byte, 3)
	read, err := os.Stdin.Read(readBytes)
	if err != nil {
		// Handle read error, it might be due to signal interruption
		return 0
	}

	// Arrow keys are prefixed with the ANSI escape code which take up the first two bytes.
	// The third byte is the key specific value we are looking for.
	// For example the left arrow key is '<esc>[A' while the right is '<esc>[C'
	// See: https://en.wikipedia.org/wiki/ANSI_escape_code
	if read == 3 {
		if _, ok := NavigationKeys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}

func containsChar(chars []byte, key byte) bool {
	for _, c := range chars {
		if c == key {
			return true
		}
	}
	return false
}
