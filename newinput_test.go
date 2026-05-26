package pardon

import (
	"bytes"
	"os"
	"testing"

	"github.com/engmtcdrm/go-pardon/internal/keys"
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
		input.pending = []byte{keys.Escape}

		doBreak := input.handleEscapeSequence()

		require.True(t, doBreak)
	})

	t.Run("should return false if second byte is not LeftBracket or CapitalO", func(t *testing.T) {
		input := NewInput()
		input.pending = []byte{keys.Escape, 'X'}

		doBreak := input.handleEscapeSequence()

		require.False(t, doBreak)
		require.Equal(t, []byte{'X'}, input.pending)
	})

	t.Run("final byte of escape sequence is in the range 0x40 to 0x7E", func(t *testing.T) {
		expected := []byte{}

		input := NewInput()
		input.pending = []byte{keys.Escape, keys.LeftBracket, 'A'}

		doBreak := input.handleEscapeSequence()

		require.False(t, doBreak)
		require.Equal(t, expected, input.pending)
	})

	t.Run("final byte of escape sequence is out of the range 0x40 to 0x7E", func(t *testing.T) {
		input := NewInput()
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
		inputText := "Hello, World!"
		expected := "Hello, World!"
		input.Writer = &bytes.Buffer{}
		input.print(inputText)
		require.Equal(t, expected, input.Writer.(*bytes.Buffer).String())
	})

	t.Run("should not print anything if the message is empty", func(t *testing.T) {
		input := NewHiddenInput()
		inputText := "Hello, World!"
		expected := ""
		input.Writer = &bytes.Buffer{}
		input.print(inputText)
		require.Equal(t, expected, input.Writer.(*bytes.Buffer).String())
	})
}

// Tests for [Input.processPending] function.
func Test_Input_processPending(t *testing.T) {
	t.Run("should process pending input as Carriage Return", func(t *testing.T) {
		expectedResult := "hello"
		input := NewInput()
		input.pending = []byte{keys.NewLine}
		input.result = []rune(expectedResult)

		returnString, done, err := input.processPending()
		require.True(t, done)
		require.NoError(t, err)
		require.Equal(t, expectedResult, returnString)
	})

	t.Run("should process pending input as Enter", func(t *testing.T) {
		expectedResult := "hello"
		input := NewInput()
		input.pending = []byte{keys.Enter}
		input.result = []rune(expectedResult)

		returnString, done, err := input.processPending()
		require.True(t, done)
		require.NoError(t, err)
		require.Equal(t, expectedResult, returnString)
	})

	t.Run("should process pending input as CtrlC", func(t *testing.T) {
		expectedResult := "hello"
		input := NewInput()
		input.pending = []byte{keys.CtrlC}
		input.result = []rune(expectedResult)

		returnString, done, err := input.processPending()
		require.True(t, done)
		require.ErrorIs(t, err, ErrUserAborted)
		require.Equal(t, "", returnString)
		require.Equal(t, expectedResult, string(input.result))
	})

	t.Run("should process pending input as Backspace", func(t *testing.T) {
		expectedResult := "hell"
		input := NewInput()
		input.pending = []byte{keys.Delete}
		input.result = []rune("hello")

		returnString, done, err := input.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Equal(t, "", returnString)
		require.Equal(t, expectedResult, string(input.result))
	})

	t.Run("should process pending input as Delete", func(t *testing.T) {
		expectedResult := "hell"
		input := NewInput()
		input.pending = []byte{keys.Backspace}
		input.result = []rune("hello")

		returnString, done, err := input.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Equal(t, "", returnString)
		require.Equal(t, expectedResult, string(input.result))
	})

	t.Run("should process pending input as Escape sequence", func(t *testing.T) {
		expectedResult := "hello"
		input := NewInput()
		input.pending = []byte{keys.Escape, keys.LeftBracket, 'A'}
		input.result = []rune(expectedResult)

		returnString, done, err := input.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Equal(t, "", returnString)
		require.Empty(t, input.pending)
		require.Equal(t, expectedResult, string(input.result))
	})

	t.Run("should process pending input as Escape sequence, too short", func(t *testing.T) {
		expectedResult := "hello"
		input := NewInput()
		input.pending = []byte{keys.Escape}
		input.result = []rune(expectedResult)

		returnString, done, err := input.processPending()
		require.False(t, done)
		require.NoError(t, err)
		require.Equal(t, "", returnString)
		require.Equal(t, []byte{}, input.pending, "pending should be empty because the escape byte should be consumed")
		require.Equal(t, expectedResult, string(input.result))
	})
}

// Tests for [Input.reset] function.
func Test_Input_reset(t *testing.T) {
	t.Run("should reset specified fields", func(t *testing.T) {
		input := NewInput()
		input.result = []rune("hello")
		input.pending = []byte("world")

		input.reset()

		require.Empty(t, input.result)
		require.Empty(t, input.pending)
	})
}
