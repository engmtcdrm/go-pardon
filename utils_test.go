package pardon

import (
	"testing"

	"github.com/engmtcdrm/go-pardon/internal/keys"
	"github.com/stretchr/testify/require"
)

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
