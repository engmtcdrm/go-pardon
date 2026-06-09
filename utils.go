package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/internal/keys"
)

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func resetLineAbove() string {
	return ansi.CursorUp(1) + ansi.ClearLineReset
}

func validationErrorMessage(err error) string {
	return fmt.Sprintf("%s%s* %v%s", ansi.ClearLineReset, ansi.RedBg, err, ansi.Reset)
}

// https://www.climagic.org/mirrors/VT100_Escape_Codes.html
func validateEscapeSequence(b byte) bool {
	switch b {
	case keys.LeftBracket, keys.UpperO:
		return true
	default:
		return false
	}
}
