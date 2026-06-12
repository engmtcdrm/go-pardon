package pardon

import (
	"github.com/engmtcdrm/go-pardon/internal/keys"
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
)
