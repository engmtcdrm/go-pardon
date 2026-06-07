package testutils

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/stretchr/testify/require"
)

// CreateValidTestFile is a helper function to create a temporary file with
// specific content for emulating user input for tests. It appends an EOF
// character to the content to ensure that the file is properly terminated for
// reading.
func CreateValidTestFile(t *testing.T, content string) *os.File {
	// content += "\x04"
	testFile := filepath.Join(t.TempDir(), "test_input.txt")

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err, "failed to write test file")

	f, err := os.Open(testFile)
	require.NoError(t, err, "failed to open test file")
	return f
}

// CreateInvalidTestFile is a helper function to create a temporary file with
// specific content for emulating user input for tests. It does not append an
// EOF character to the content, which can be used to test how the code handles
// unexpected end of input.
func CreateInvalidTestFile(t *testing.T, content string) *os.File {
	testFile := filepath.Join(t.TempDir(), "test_input.txt")

	f, err := os.Create(testFile)
	require.NoError(t, err, "failed to create test file")

	_, err = f.WriteString(content)
	require.NoError(t, err, "failed to write to test file")

	return f
}

// CreatePTY creates a pseudo-terminal pair and writes the provided content to
// the master end so the returned slave *os.File can be used as a terminal
// reader in tests. The master is closed after writing so the slave will
// observe EOF when appropriate.
func CreatePTY(t *testing.T, content string) *os.File {
	m, s, err := pty.Open()
	require.NoError(t, err, "failed to open pty")
	t.Cleanup(func() { s.Close() })

	// Write content as if typed by a user to ensure ordering of bytes so
	// that the reader processes printable runes before the newline is
	// encountered. Close the master after writing so the slave observes EOF.
	go func() {
		defer m.Close()
		for i := 0; i < len(content); i++ {
			_, _ = m.Write([]byte{content[i]})
			time.Sleep(5 * time.Millisecond)
		}
	}()

	return s
}
