package testutils

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for [CreateValidTestFile] function.
func Test_CreateValidTestFile(t *testing.T) {
	t.Run("should create a temporary file with the specified content and an EOF character", func(t *testing.T) {
		content := "Hello, World!"
		f := CreateValidTestFile(t, content)

		data, err := io.ReadAll(f)
		require.NoError(t, err, "failed to read test file")
		require.Equal(t, content, string(data))
	})
}

// Tests for [CreateInvalidTestFile] function.
func Test_CreateInvalidTestFile(t *testing.T) {
	t.Run("should create a temporary file with the specified content without an EOF character", func(t *testing.T) {
		content := "Hello, World!"
		f := CreateInvalidTestFile(t, content)

		data, err := io.ReadAll(f)
		require.NoError(t, err, "failed to read test file")
		require.NotEqual(t, content, string(data))
	})
}
