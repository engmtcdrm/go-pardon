package pardon

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/engmtcdrm/go-pardon/internal/testutils"
	"github.com/stretchr/testify/require"
)

// Tests for [NewTerminal] function.
func Test_NewTerminal(t *testing.T) {
	t.Run("should create a new Terminal instance with default values", func(t *testing.T) {
		terminal := NewTerminal()
		require.NotNil(t, terminal)
		require.Equal(t, terminal.Hide, false)
		require.Equal(t, terminal.Out, os.Stdout)
		require.Equal(t, terminal.In, os.Stdin)
	})
}

// Tests for [NewHiddenTerminal] function.
func Test_NewHiddenTerminal(t *testing.T) {
	t.Run("should create a new Terminal instance with Hide set to true", func(t *testing.T) {
		terminal := NewHiddenTerminal()
		require.NotNil(t, terminal)
		require.Equal(t, terminal.Hide, true)
		require.Equal(t, terminal.Out, os.Stdout)
		require.Equal(t, terminal.In, os.Stdin)
	})
}

// Tests for [Terminal.RawRead] function.
func Test_Terminal_RawRead(t *testing.T) {
	t.Run("should read input from the In and return it as a slice of runes", func(t *testing.T) {
		_, f := testutils.CreateWritePTY(t, "hello\n")

		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.In = f

		var results []byte

		for {
			input, err := terminal.RawRead()
			if err != nil {
				require.NoError(t, err)
			}

			results = append(results, input...)

			if len(results) > 5 {
				break
			}
		}

		require.Equal(t, []byte("hello\n"), results)
	})

	t.Run("should return an error if In is not a file", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.In = nil

		results, err := terminal.RawRead()
		require.Error(t, err)
		require.Nil(t, results)
	})

	t.Run("should read input from the In if not a terminal but valid file", func(t *testing.T) {
		f := testutils.CreateValidTestFile(t, "hello\n")
		defer f.Close()

		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.In = f

		var results []byte

		for {
			input, err := terminal.RawRead()
			if err != nil {
				require.NoError(t, err)
			}

			results = append(results, input...)

			if len(results) > 5 {
				break
			}
		}

		require.Equal(t, []byte("hello\n"), results)
	})

	t.Run("should error when term.MakeRaw fails", func(t *testing.T) {
		t.Skip("Cannot reliably test term.MakeRaw failure without mocking, and mocking is not currently implemented.")
	})
}

// Tests for [Terminal.print] function.
func Test_Terminal_print(t *testing.T) {
	t.Run("should print the message to the Writer", func(t *testing.T) {
		terminal := NewTerminal()
		inputText := "hello"
		expected := "hello"
		terminal.Out = &bytes.Buffer{}
		terminal.PrintInput(inputText)
		require.Equal(t, expected, terminal.Out.(*bytes.Buffer).String())
	})

	t.Run("should not print anything if the message is empty", func(t *testing.T) {
		terminal := NewHiddenTerminal()
		inputText := "hello"
		expected := ""
		terminal.Out = &bytes.Buffer{}
		terminal.PrintInput(inputText)
		require.Equal(t, expected, terminal.Out.(*bytes.Buffer).String())
	})
}
