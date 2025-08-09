package pardon

import (
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "user aborted error",
			err:      ErrUserAborted,
			expected: "user aborted",
		},
		{
			name:     "no title error",
			err:      ErrNoTitle,
			expected: "prompt requires a title",
		},
		{
			name:     "no select options error",
			err:      ErrNoSelectOptions,
			expected: "select prompt requires at least one option",
		},
		{
			name:     "no value error",
			err:      ErrNoValue,
			expected: "value must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Error message = %q; want %q", tt.err.Error(), tt.expected)
			}
		})
	}
}

func TestErrorsAreNotNil(t *testing.T) {
	errors := []error{
		ErrUserAborted,
		ErrNoTitle,
		ErrNoSelectOptions,
		ErrNoValue,
	}

	for i, err := range errors {
		if err == nil {
			t.Errorf("Error at index %d is nil", i)
		}
	}
}
