package pardon

import "github.com/engmtcdrm/go-pardon/internal/keys"

// https://www.climagic.org/mirrors/VT100_Escape_Codes.html
func validateEscapeSequence(b byte) bool {
	switch b {
	case keys.LeftBracket, keys.CapitalO:
		return true
	default:
		return false
	}
}
