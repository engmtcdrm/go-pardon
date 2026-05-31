package pardon

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/engmtcdrm/go-pardon/internal/keys"
	"github.com/engmtcdrm/go-pardon/internal/runekeys"
	"github.com/engmtcdrm/go-pardon/internal/testutils"
	"github.com/stretchr/testify/require"
)

// Tests for [NewInput] function.
func Test_NewInput(t *testing.T) {
	t.Run("should create a new Input instance with default values", func(t *testing.T) {
		input := NewInput()
		require.NotNil(t, input)
		require.Equal(t, input.Hide, false)
		require.Equal(t, input.Writer, os.Stdout)
		require.Equal(t, input.Reader, os.Stdin)
	})
}

// Tests for [NewHiddenInput] function.
func Test_NewHiddenInput(t *testing.T) {
	t.Run("should create a new Input instance with Hide set to true", func(t *testing.T) {
		input := NewHiddenInput()
		require.NotNil(t, input)
		require.Equal(t, input.Hide, true)
		require.Equal(t, input.Writer, os.Stdout)
		require.Equal(t, input.Reader, os.Stdin)
	})
}

// Tests for [NewConfirmInput] function.
func Test_NewConfirmInput(t *testing.T) {
	t.Run("should create a new Input instance with Hide set to true", func(t *testing.T) {
		input := NewConfirmInput()
		require.NotNil(t, input)
		require.Equal(t, input.Hide, false)
		require.Equal(t, input.Confirm, true)
		require.Equal(t, input.Writer, os.Stdout)
		require.Equal(t, input.Reader, os.Stdin)
	})
}

// Tests for [Input.RawRead] function.
func Test_Input_RawRead(t *testing.T) {
	t.Run("should read input from the Reader and return it as a slice of runes", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("pty tests skipped on Windows. Pty is not supported.")
		}

		f := testutils.CreatePTY(t, "hello\n")
		defer f.Close()

		input := NewInput()
		input.Writer = io.Discard
		input.Reader = f

		results, err := input.RawRead()
		require.NoError(t, err)
		require.Equal(t, []rune("hello"), results)
	})

	t.Run("should return an error if Reader is not a file", func(t *testing.T) {
		input := NewInput()
		input.Writer = io.Discard
		input.Reader = nil

		results, err := input.RawRead()
		require.Error(t, err)
		require.Nil(t, results)
	})

	t.Run("should read input from the Reader if not a terminal but valid file", func(t *testing.T) {
		f := testutils.CreateValidTestFile(t, "hello\n")
		defer f.Close()

		input := NewInput()
		input.Writer = io.Discard
		input.Reader = f

		results, err := input.RawRead()
		require.NoError(t, err)
		require.Equal(t, []rune("hello"), results)
	})

	t.Run("should error when term.MakeRaw fails", func(t *testing.T) {
		t.Skip("Cannot reliably test term.MakeRaw failure without mocking, and mocking is not currently implemented.")
	})
}

// Tests for [Input.Reset] function.
func Test_Input_Reset(t *testing.T) {
	t.Run("should reset specified fields", func(t *testing.T) {
		input := NewInput()
		input.result = []rune("hello")
		input.pending = []byte("world")

		input.Reset()

		require.Empty(t, input.result)
		require.Empty(t, input.pending)
	})
}

// Tests for [Input.handleErase] function.
func Test_Input_handleErase(t *testing.T) {
	t.Run("should remove the last character from result and pending", func(t *testing.T) {
		input := NewInput()
		input.Writer = &bytes.Buffer{}
		input.result = []rune("hello")
		input.pending = []byte{keys.Delete, 'o', 'r', 'l', 'd'}

		input.handleErase()

		require.Equal(t, []rune("hell"), input.result)
		require.Equal(t, []byte("orld"), input.pending)
		require.Equal(t, "\b \b", input.Writer.(*bytes.Buffer).String())
	})
}

// Tests for [Input.handleEscapeSequence] function.
func Test_Input_handleEscapeSequence(t *testing.T) {
	t.Run("should return true if pending is less than 2 bytes", func(t *testing.T) {
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Escape}

		doBreak := input.handleEscapeSequence()

		require.True(t, doBreak)
	})

	t.Run("should return false if second byte is not LeftBracket or CapitalO", func(t *testing.T) {
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Escape, 'X'}

		doBreak := input.handleEscapeSequence()

		require.False(t, doBreak)
		require.Equal(t, []byte{'X'}, input.pending)
	})

	t.Run("final byte of escape sequence is in the range 0x40 to 0x7E", func(t *testing.T) {
		expected := []byte{}

		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Escape, keys.LeftBracket, 'A'}

		doBreak := input.handleEscapeSequence()

		require.False(t, doBreak)
		require.Equal(t, expected, input.pending)
	})

	t.Run("final byte of escape sequence is out of the range 0x40 to 0x7E", func(t *testing.T) {
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Escape, keys.LeftBracket, byte(32)}

		doBreak := input.handleEscapeSequence()

		require.True(t, doBreak)
		require.Equal(t, []byte{keys.Escape, keys.LeftBracket, byte(32)}, input.pending)
	})
}

// Tests for [Input.print] function.
func Test_Input_print(t *testing.T) {
	t.Run("should print the message to the Writer", func(t *testing.T) {
		input := NewInput()
		inputText := "hello"
		expected := "hello"
		input.Writer = &bytes.Buffer{}
		input.print(inputText)
		require.Equal(t, expected, input.Writer.(*bytes.Buffer).String())
	})

	t.Run("should not print anything if the message is empty", func(t *testing.T) {
		input := NewHiddenInput()
		inputText := "hello"
		expected := ""
		input.Writer = &bytes.Buffer{}
		input.print(inputText)
		require.Equal(t, expected, input.Writer.(*bytes.Buffer).String())
	})
}

// Tests for [Input.processPending] function.
func Test_Input_processPending(t *testing.T) {
	t.Run("should process pending input as Carriage Return", func(t *testing.T) {
		expectedResult := []rune("hello")
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.NewLine}
		input.result = []rune(expectedResult)

		returnRunes, done, err := input.processPending()
		require.True(t, done)
		require.NoError(t, err)
		require.Equal(t, expectedResult, returnRunes)
	})

	t.Run("should process pending input as Enter", func(t *testing.T) {
		expectedResult := []rune("hello")
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Enter}
		input.result = []rune(expectedResult)

		returnRunes, done, err := input.processPending()
		require.True(t, done)
		require.NoError(t, err)
		require.Equal(t, expectedResult, returnRunes)
	})

	t.Run("should process pending input as CtrlC", func(t *testing.T) {
		expectedResult := []rune("hello")
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.CtrlC}
		input.result = []rune(expectedResult)

		returnRunes, done, err := input.processPending()
		require.True(t, done)
		require.ErrorIs(t, err, ErrUserAborted)
		require.Nil(t, returnRunes)
		require.Equal(t, expectedResult, input.result)
	})

	t.Run("should process pending input as Backspace", func(t *testing.T) {
		expectedResult := []rune("hell")
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Delete}
		input.result = []rune("hello")

		returnRunes, done, err := input.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Nil(t, returnRunes)
		require.Equal(t, expectedResult, input.result)
	})

	t.Run("should process pending input as Delete", func(t *testing.T) {
		expectedResult := []rune("hell")
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Backspace}
		input.result = []rune("hello")

		returnRunes, done, err := input.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Nil(t, returnRunes)
		require.Equal(t, expectedResult, input.result)
	})

	t.Run("should process pending input as Escape sequence", func(t *testing.T) {
		expectedResult := []rune("hello")
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Escape, keys.LeftBracket, 'A'}
		input.result = []rune(expectedResult)

		returnRunes, done, err := input.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Nil(t, returnRunes)
		require.Empty(t, input.pending)
		require.Equal(t, expectedResult, input.result)
	})

	t.Run("should process pending input as Escape sequence, too short", func(t *testing.T) {
		expectedResult := []rune("hello")
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.Escape}
		input.result = []rune(expectedResult)

		returnRunes, done, err := input.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Nil(t, returnRunes)
		require.Equal(t, []byte{}, input.pending, "pending should be empty because the escape byte should be consumed")
		require.Equal(t, expectedResult, input.result)
	})

	t.Run("invalid rune in pending input", func(t *testing.T) {
		input := NewInput()
		input.Writer = io.Discard
		input.pending = []byte{'H', 0xFF, 'I', keys.Enter} // Invalid UTF-8 byte

		returnRunes, done, err := input.processPending()
		require.NoError(t, err)
		require.True(t, done)
		require.NotNil(t, returnRunes)
		require.Equal(t, []rune{'H', 'I'}, returnRunes, "pending should remain unchanged on invalid rune")
	})

	t.Run("invalid full rune in pending input", func(t *testing.T) {
		input := NewInput()
		input.Writer = io.Discard
		// 0xE2 starts a 3-byte UTF-8 sequence (but buffer is incomplete)
		input.pending = []byte{0xE2}

		returnRunes, done, err := input.processPending()
		require.NoError(t, err)
		require.False(t, done)
		require.Nil(t, returnRunes, "runes should be nil on invalid full rune")
		require.Equal(t, []byte{0xE2}, input.pending, "expected pending unchanged")
	})

	t.Run("valid confirm input", func(t *testing.T) {
		input := NewConfirmInput()
		input.Writer = io.Discard
		input.pending = []byte{keys.UpperN}

		returnRunes, done, err := input.processPending()
		require.NoError(t, err)
		require.True(t, done)
		require.Equal(t, []rune{runekeys.UpperN}, returnRunes)
	})
}

// Tests for [Input.rawReadline] function.
func Test_Input_rawReadline(t *testing.T) {
	t.Run("should read runes from the file until a newline is encountered", func(t *testing.T) {
		f := testutils.CreateValidTestFile(t, "hello\n")
		defer f.Close()

		input := NewInput()
		input.Writer = io.Discard

		results, err := input.rawReadline(f)
		require.NoError(t, err)
		require.Equal(t, []rune("hello"), results)
	})

	t.Run("should return an error if there is an issue reading from the file", func(t *testing.T) {
		f := testutils.CreateInvalidTestFile(t, "hello\n")
		defer f.Close()

		input := NewInput()
		input.Writer = io.Discard
		results, err := input.rawReadline(f)
		require.NoError(t, err)
		require.Nil(t, results)
	})

	t.Run("should return an error if the file is nil", func(t *testing.T) {
		input := NewInput()
		input.Writer = io.Discard

		results, err := input.rawReadline(nil)
		require.Error(t, err)
		require.Nil(t, results)
	})
}
