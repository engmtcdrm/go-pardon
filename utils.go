package pardon

import (
	"strings"

	"github.com/engmtcdrm/go-ansi"
)

const (
	saveCursor    = "\x1b7"
	restoreCursor = "\x1b8"
)

func validationErrorMessage(err error) string {
	var builder strings.Builder
	builder.WriteString(ansi.ClearLineReset)
	builder.WriteString(ansi.RedBg)
	builder.WriteString("* ")
	builder.WriteString(err.Error())
	builder.WriteString(ansi.Reset)

	return builder.String()
}

// zeroParent returns the zero value for the generic type P. This is used to
// initialize the Self field in the BasePrompt struct to a zero value of the
// concrete parent type.
func zeroParent[P any]() P {
	var p P
	return p
}
