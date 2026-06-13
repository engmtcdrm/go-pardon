package pardon

import (
	"errors"
	"fmt"
	"testing"

	"github.com/engmtcdrm/go-ansi"
	"github.com/stretchr/testify/assert"
)

// Tests for [resetLineAbove] function.
func Test_resetLineAbove(t *testing.T) {
	expected := ansi.CursorUp(1) + ansi.ClearLineReset
	result := resetLineAbove()
	assert.Equal(t, expected, result, "resetLineAbove did not return expected string")
}

// Tests for [validationErrorMessage] function.
func Test_validationErrorMessage(t *testing.T) {
	t.Run("should return expected string", func(t *testing.T) {
		err := errors.New("Invalid input")
		expected := fmt.Sprintf("%s%s* %v%s", ansi.ClearLineReset, ansi.RedBg, err, ansi.Reset)

		result := validationErrorMessage(err)
		assert.Equal(t, expected, result, "validationErrorMessage did not return expected string")
	})
}
