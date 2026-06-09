package pardon

import (
	"errors"
	"fmt"
	"testing"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/internal/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for [min] function.
func Test_min(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"a smaller than b", 3, 7, 3},
		{"b smaller than a", 10, 5, 5},
		{"equal values", 4, 4, 4},
		{"negative values", -3, -7, -7},
		{"zero values", 0, 0, 0},
		{"negative and positive", -5, 3, -5},
		// -2147483648, 2147483647
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := min(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Min(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

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

// Tests for [validateEscapeSequence] function.
func Test_validateEscapeSequence(t *testing.T) {
	type cases struct {
		name     string
		input    byte
		expected bool
	}

	tests := []cases{
		{
			name:     "valid escape sequence byte: LeftBracket",
			input:    keys.LeftBracket,
			expected: true,
		},
		{
			name:     "valid escape sequence byte: CapitalO",
			input:    keys.UpperO,
			expected: true,
		},
		{
			name:     "invalid escape sequence byte: \\n",
			input:    '\n',
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := validateEscapeSequence(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
