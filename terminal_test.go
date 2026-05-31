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

// Tests for [NewConfirmTerminal] function.
func Test_NewConfirmTerminal(t *testing.T) {
	t.Run("should create a new Terminal instance with Confirm set to true", func(t *testing.T) {
		terminal := NewConfirmTerminal()
		require.NotNil(t, terminal)
		require.Equal(t, terminal.Hide, false)
		require.Equal(t, terminal.Confirm, true)
		require.Equal(t, terminal.Out, os.Stdout)
		require.Equal(t, terminal.In, os.Stdin)
	})
}

// Tests for [Terminal.RawRead] function.
func Test_Terminal_RawRead(t *testing.T) {
	t.Run("should read input from the In and return it as a slice of runes", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("pty tests skipped on Windows. Pty is not supported.")
		}

		f := testutils.CreatePTY(t, "hello\n")
		defer f.Close()

		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.In = f

		results, err := terminal.RawRead()
		require.NoError(t, err)
		require.Equal(t, []rune("hello"), results)
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

		results, err := terminal.RawRead()
		require.NoError(t, err)
		require.Equal(t, []rune("hello"), results)
	})

	t.Run("should error when term.MakeRaw fails", func(t *testing.T) {
		t.Skip("Cannot reliably test term.MakeRaw failure without mocking, and mocking is not currently implemented.")
	})
}

// Tests for [Terminal.Reset] function.
func Test_Terminal_Reset(t *testing.T) {
	t.Run("should reset specified fields", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.result = []rune("hello")
		terminal.pending = []byte("world")

		terminal.Reset()

		require.Empty(t, terminal.result)
		require.Empty(t, terminal.pending)
	})
}

// Tests for [Terminal.handleErase] function.
func Test_Terminal_handleErase(t *testing.T) {
	t.Run("should remove the last character from result and pending", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.Out = &bytes.Buffer{}
		terminal.result = []rune("hello")
		terminal.pending = []byte{keys.Delete, 'o', 'r', 'l', 'd'}

		terminal.handleErase()

		require.Equal(t, []rune("hell"), terminal.result)
		require.Equal(t, []byte("orld"), terminal.pending)
		require.Equal(t, "\b \b", terminal.Out.(*bytes.Buffer).String())
	})
}

// Tests for [Terminal.handleEscapeSequence] function.
func Test_Terminal_handleEscapeSequence(t *testing.T) {
	t.Run("should return true if pending is less than 2 bytes", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Escape}

		doBreak := terminal.handleEscapeSequence()

		require.True(t, doBreak)
	})

	t.Run("should return false if second byte is not LeftBracket or CapitalO", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Escape, 'X'}

		doBreak := terminal.handleEscapeSequence()

		require.False(t, doBreak)
		require.Equal(t, []byte{'X'}, terminal.pending)
	})

	t.Run("final byte of escape sequence is in the range 0x40 to 0x7E", func(t *testing.T) {
		expected := []byte{}

		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Escape, keys.LeftBracket, 'A'}

		doBreak := terminal.handleEscapeSequence()

		require.False(t, doBreak)
		require.Equal(t, expected, terminal.pending)
	})

	t.Run("final byte of escape sequence is out of the range 0x40 to 0x7E", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Escape, keys.LeftBracket, byte(32)}

		doBreak := terminal.handleEscapeSequence()

		require.True(t, doBreak)
		require.Equal(t, []byte{keys.Escape, keys.LeftBracket, byte(32)}, terminal.pending)
	})
}

// Tests for [Terminal.print] function.
func Test_Terminal_print(t *testing.T) {
	t.Run("should print the message to the Writer", func(t *testing.T) {
		terminal := NewTerminal()
		inputText := "hello"
		expected := "hello"
		terminal.Out = &bytes.Buffer{}
		terminal.print(inputText)
		require.Equal(t, expected, terminal.Out.(*bytes.Buffer).String())
	})

	t.Run("should not print anything if the message is empty", func(t *testing.T) {
		terminal := NewHiddenTerminal()
		inputText := "hello"
		expected := ""
		terminal.Out = &bytes.Buffer{}
		terminal.print(inputText)
		require.Equal(t, expected, terminal.Out.(*bytes.Buffer).String())
	})
}

// Tests for [Terminal.processPending] function.
func Test_Terminal_processPending(t *testing.T) {
	t.Run("should process pending input as Carriage Return", func(t *testing.T) {
		expectedResult := []rune("hello")
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.NewLine}
		terminal.result = []rune(expectedResult)

		returnRunes, done, err := terminal.processPending()
		require.True(t, done)
		require.NoError(t, err)
		require.Equal(t, expectedResult, returnRunes)
	})

	t.Run("should process pending input as Enter", func(t *testing.T) {
		expectedResult := []rune("hello")
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Enter}
		terminal.result = []rune(expectedResult)

		returnRunes, done, err := terminal.processPending()
		require.True(t, done)
		require.NoError(t, err)
		require.Equal(t, expectedResult, returnRunes)
	})

	t.Run("should process pending input as CtrlC", func(t *testing.T) {
		expectedResult := []rune("hello")
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.CtrlC}
		terminal.result = []rune(expectedResult)

		returnRunes, done, err := terminal.processPending()
		require.True(t, done)
		require.ErrorIs(t, err, ErrUserAborted)
		require.Nil(t, returnRunes)
		require.Equal(t, expectedResult, terminal.result)
	})

	t.Run("should process pending input as Backspace", func(t *testing.T) {
		expectedResult := []rune("hell")
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Delete}
		terminal.result = []rune("hello")

		returnRunes, done, err := terminal.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Nil(t, returnRunes)
		require.Equal(t, expectedResult, terminal.result)
	})

	t.Run("should process pending input as Delete", func(t *testing.T) {
		expectedResult := []rune("hell")
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Backspace}
		terminal.result = []rune("hello")

		returnRunes, done, err := terminal.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Nil(t, returnRunes)
		require.Equal(t, expectedResult, terminal.result)
	})

	t.Run("should process pending input as Escape sequence", func(t *testing.T) {
		expectedResult := []rune("hello")
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Escape, keys.LeftBracket, 'A'}
		terminal.result = []rune(expectedResult)

		returnRunes, done, err := terminal.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Nil(t, returnRunes)
		require.Empty(t, terminal.pending)
		require.Equal(t, expectedResult, terminal.result)
	})

	t.Run("should process pending input as Escape sequence, too short", func(t *testing.T) {
		expectedResult := []rune("hello")
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.Escape}
		terminal.result = []rune(expectedResult)

		returnRunes, done, err := terminal.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Nil(t, returnRunes)
		require.Equal(t, []byte{}, terminal.pending, "pending should be empty because the escape byte should be consumed")
		require.Equal(t, expectedResult, terminal.result)
	})

	t.Run("invalid rune in pending input", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{'H', 0xFF, 'I', keys.Enter} // Invalid UTF-8 byte

		returnRunes, done, err := terminal.processPending()
		require.NoError(t, err)
		require.True(t, done)
		require.NotNil(t, returnRunes)
		require.Equal(t, []rune{'H', 'I'}, returnRunes, "pending should remain unchanged on invalid rune")
	})

	t.Run("invalid full rune in pending input", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.Out = io.Discard
		// 0xE2 starts a 3-byte UTF-8 sequence (but buffer is incomplete)
		terminal.pending = []byte{0xE2}

		returnRunes, done, err := terminal.processPending()
		require.NoError(t, err)
		require.False(t, done)
		require.Nil(t, returnRunes, "runes should be nil on invalid full rune")
		require.Equal(t, []byte{0xE2}, terminal.pending, "expected pending unchanged")
	})

	t.Run("valid confirm input", func(t *testing.T) {
		terminal := NewConfirmTerminal()
		terminal.Out = io.Discard
		terminal.pending = []byte{keys.UpperN}

		returnRunes, done, err := terminal.processPending()
		require.NoError(t, err)
		require.True(t, done)
		require.Equal(t, []rune{runekeys.UpperN}, returnRunes)
	})
}

// Tests for [Terminal.rawReadline] function.
func Test_Terminal_rawReadline(t *testing.T) {
	t.Run("should read runes from the file until a newline is encountered", func(t *testing.T) {
		f := testutils.CreateValidTestFile(t, "hello\n")
		defer f.Close()

		terminal := NewTerminal()
		terminal.Out = io.Discard

		results, err := terminal.rawReadline(f)
		require.NoError(t, err)
		require.Equal(t, []rune("hello"), results)
	})

	t.Run("should return an error if there is an issue reading from the file", func(t *testing.T) {
		f := testutils.CreateInvalidTestFile(t, "hello\n")
		defer f.Close()

		terminal := NewTerminal()
		terminal.Out = io.Discard
		results, err := terminal.rawReadline(f)
		require.NoError(t, err)
		require.Nil(t, results)
	})

	t.Run("should return an error if the file is nil", func(t *testing.T) {
		terminal := NewTerminal()
		terminal.Out = io.Discard

		results, err := terminal.rawReadline(nil)
		require.Error(t, err)
		require.Nil(t, results)
	})
}
