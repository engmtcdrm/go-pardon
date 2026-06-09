package testutils

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
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
	t.Helper()

	testFile := filepath.Join(t.TempDir(), "test_input.txt")

	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err, "failed to write test file")

	f, err := os.Open(testFile)
	require.NoError(t, err, "failed to open test file")
	t.Cleanup(func() {
		f.Close()
	})

	return f
}

// CreateInvalidTestFile is a helper function to create a temporary file with
// specific content for emulating user input for tests. It does not append an
// EOF character to the content, which can be used to test how the code handles
// unexpected end of input.
func CreateInvalidTestFile(t *testing.T, content string) *os.File {
	t.Helper()

	testFile := filepath.Join(t.TempDir(), "test_input.txt")

	f, err := os.Create(testFile)
	require.NoError(t, err, "failed to create test file")

	_, err = f.WriteString(content)
	require.NoError(t, err, "failed to write to test file")

	t.Cleanup(func() {
		f.Close()
	})

	return f
}

// CreatePTY creates a pseudo-terminal pair for testing terminal interactions.
// The returned master and slave *os.File can be used to simulate terminal input
// and output in tests. The master end can be used to write input as if typed by
// a user, while the slave end can be used to read output from the terminal.
// Both files are automatically closed after the test completes.
func CreatePTY(t *testing.T) (master *os.File, slave *os.File) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("pty is not supported on Windows")
	}

	m, s, err := pty.Open()
	require.NoError(t, err, "failed to open pty")
	t.Cleanup(func() {
		m.Close()
		s.Close()
	})

	return m, s
}

// CreatePTYWithSize creates a pseudo-terminal pair with the specified size for
// testing terminal interactions. The returned master and slave *os.File can be
// used to simulate terminal input and output in tests that require specific
// terminal dimensions. Both files are automatically closed after the test
// completes.
func CreatePTYWithSize(t *testing.T, columns, rows int) (master *os.File, slave *os.File) {
	t.Helper()

	master, slave = CreatePTY(t)

	err := pty.Setsize(slave, &pty.Winsize{Cols: uint16(columns), Rows: uint16(rows)})
	require.NoError(t, err, "failed to set pty size")

	return master, slave
}

// CreateWritePTY creates a pseudo-terminal pair and writes the provided content
// to the master end. The returned slave *os.File can be used as a terminal
// reader in tests. The master is closed after writing so the slave will observe
// EOF when appropriate.
func CreateWritePTY(t *testing.T, content string) (master *os.File, slave *os.File) {
	t.Helper()

	m, s := CreatePTY(t)

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

	return m, s
}

// CreateWritePTYWithSize creates a pseudo-terminal pair with the specified size
// and writes the provided content to the master end. The returned slave
// *os.File can be used as a terminal reader in tests that require specific
// terminal dimensions. The master is closed after writing so the slave will
// observe EOF when appropriate.
func CreateWritePTYWithSize(t *testing.T, content string, columns, rows int) (master *os.File, slave *os.File) {
	t.Helper()

	master, slave = CreateWritePTY(t, content)

	err := pty.Setsize(slave, &pty.Winsize{Cols: uint16(columns), Rows: uint16(rows)})
	require.NoError(t, err, "failed to set pty size")

	return master, slave
}

// ReadPTYOutput reads all available output from the provided pseudo-terminal
// file until EOF is reached.
func ReadPTYOutput(t *testing.T, ptyFile *os.File, bufferSize int) string {
	t.Helper()

	if bufferSize <= 0 {
		bufferSize = 1024
	}

	var output bytes.Buffer
	buffer := make([]byte, bufferSize)
	for {
		n, readErr := ptyFile.Read(buffer)
		if n > 0 {
			_, _ = output.Write(buffer[:n])
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) || errors.Is(readErr, syscall.EIO) {
				break
			}
			require.NoError(t, readErr, "reading output should not return an error")
		}
		if n == 0 {
			break
		}
	}

	return output.String()
}
