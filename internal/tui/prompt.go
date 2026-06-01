package tui

import (
	"os"

	"github.com/engmtcdrm/go-pardon/internal/keys"
	"golang.org/x/term"
)

var (
	// navigationKeys defines a map of byte keycodes for navigation actions.
	// These keys are used for cursor movement and selection in interactive prompts.
	navigationKeys = map[byte]bool{
		keys.Up:    true,
		keys.Down:  true,
		keys.Left:  true,
		keys.Right: true,
	}

	// lastInputWasEscSeq tracks whether the previous input was part of an escape sequence.
	// This helps with proper handling of multi-byte terminal input sequences.
	lastInputWasEscSeq = false

	// inputBuffer provides buffering for handling paste operations and multi-byte input.
	inputBuffer []byte
)

// GetInput reads raw keyboard input from the terminal.
// Handles buffered input, raw mode, and ANSI escape sequences.
func GetInput() byte {
	// If we have buffered input from a paste operation, return it first
	if len(inputBuffer) > 0 {
		result := inputBuffer[0]
		inputBuffer = inputBuffer[1:]
		lastInputWasEscSeq = false
		return result
	}

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

	// Read input - use a larger buffer to handle paste operations
	readBytes := make([]byte, 4096) // Increased to 4KB to handle larger pastes
	read, err := os.Stdin.Read(readBytes)
	if err != nil {
		// Handle read error, it might be due to signal interruption
		return 0
	}

	// If we read more than 3 bytes, it's likely a paste operation
	if read > 3 {
		// Buffer all characters except the first one
		inputBuffer = append(inputBuffer, readBytes[1:read]...)
		lastInputWasEscSeq = false
		return readBytes[0]
	}

	// Handle escape sequences (arrow keys)
	if read == 3 && readBytes[0] == keys.Escape && readBytes[1] == keys.LeftBracket {
		// This is a proper ANSI escape sequence (ESC[X)
		if _, ok := navigationKeys[readBytes[2]]; ok {
			lastInputWasEscSeq = true
			return readBytes[2]
		}
		// If it's an escape sequence but not a navigation key, ignore it
		lastInputWasEscSeq = false
		return 0
	}

	// For any other input (1, 2, or 3 bytes that aren't escape sequences),
	// return the first byte which contains the actual character
	lastInputWasEscSeq = false
	if read > 0 {
		return readBytes[0]
	}

	return 0
}
