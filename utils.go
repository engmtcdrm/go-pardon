package pardon

import (
	"fmt"

	"github.com/engmtcdrm/go-ansi"
)

const (
	saveCursor    = "\x1b7"
	restoreCursor = "\x1b8"
)

func resetLineAbove() string {
	return ansi.CursorUp(1) + ansi.ClearLineReset
}

func validationErrorMessage(err error) string {
	return fmt.Sprintf("%s%s* %v%s", ansi.ClearLineReset, ansi.RedBg, err, ansi.Reset)
}

// zeroParent returns the zero value for the generic type P. This is used to
// initialize the Self field in the BasePrompt struct to a zero value of the
// concrete parent type.
func zeroParent[P any]() P {
	var p P
	return p
}
