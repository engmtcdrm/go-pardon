package pardon

import (
	"errors"
	"fmt"
	"testing"

	"github.com/engmtcdrm/go-ansi"
	"github.com/stretchr/testify/assert"
)

// Tests for [validationErrorMessage] function.
func Test_validationErrorMessage(t *testing.T) {
	t.Run("should return expected string", func(t *testing.T) {
		err := errors.New("Invalid input")
		expected := fmt.Sprintf("%s%s* %v%s", ansi.ClearLineReset, ansi.RedBg, err, ansi.Reset)

		result := validationErrorMessage(err)
		assert.Equal(t, expected, result, "validationErrorMessage did not return expected string")
	})
}
