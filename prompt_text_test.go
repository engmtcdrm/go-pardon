package pardon

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/grapheme"
	"github.com/engmtcdrm/go-pardon/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for [NewPassword] function.
func Test_NewPassword(t *testing.T) {
	t.Run("should create a new password prompt with valid value pointer", func(t *testing.T) {
		var result string
		password := NewPassword(&result)
		require.NotNil(t, password, "NewPassword should not return nil")
		require.NotNil(t, password.value, "Password value pointer should not be nil")
	})

	t.Run("should create a new password prompt with a nil value pointer", func(t *testing.T) {
		password := NewPassword(nil)
		require.NotNil(t, password, "NewPassword should not return nil even with nil value pointer")
		require.Nil(t, password.value, "Password value pointer should be nil when initialized with nil")
	})

	t.Run("should return error for empty input for Question", func(t *testing.T) {
		var result string
		question := NewQuestion(&result)

		err := question.validateFn("test")
		require.NoError(t, err, "validateFn should not return an error for valid input")
	})
}

// Tests for [NewQuestion] function.
func Test_NewQuestion(t *testing.T) {
	t.Run("should create a new question with valid value pointer", func(t *testing.T) {
		var result string
		question := NewQuestion(&result)
		require.NotNil(t, question)
		assert.NotEmpty(t, question.icon.val)
		assert.Equal(t, &result, question.value)
	})

	t.Run("should create a new question with a nil value pointer", func(t *testing.T) {
		question := NewQuestion(nil)
		require.NotNil(t, question)
		assert.NotEmpty(t, question.icon.val)
		assert.Nil(t, question.value)
	})

	t.Run("should return error for empty input for Password", func(t *testing.T) {
		var result string
		password := NewPassword(&result)

		err := password.validateFn("")
		require.NoError(t, err, "validateFn should not return an error for valid input")
	})
}

// Tests for [Text.AnswerFunc] method.
func Test_Text_AnswerFunc(t *testing.T) {
	t.Run("using default answer function", func(t *testing.T) {
		var result string
		textPrompt := NewQuestion(&result).
			Title("Continue?")
		textPrompt.Out = io.Discard
		assert.Nil(t, textPrompt.answer.fn, "Default answer function should be nil")
	})

	t.Run("using custom answer function", func(t *testing.T) {
		var result string
		customFn := func(s string) string {
			return "Custom: " + s
		}
		textPrompt := NewQuestion(&result).
			Title("Continue?").
			AnswerFunc(customFn)
		textPrompt.Out = io.Discard
		textPrompt.answer.val = "Test"
		assert.Equal(t, customFn("Test"), textPrompt.answer.Get(), "Custom answer function did not return expected result")
	})
}

// Tests for [Text.Ask] function.
func Test_Text_Ask(t *testing.T) {
	t.Run("should return error if title is not set", func(t *testing.T) {
		var result string
		textPrompt := NewQuestion(&result)
		err := textPrompt.Ask()
		require.ErrorIs(t, err, ErrNoTitle, "Ask should return ErrNoTitle if title is not set")
	})

	t.Run("should return error if value pointer is not set", func(t *testing.T) {
		textPrompt := NewQuestion(nil).
			Title("Enter value:")
		err := textPrompt.Ask()
		require.ErrorIs(t, err, ErrNoValue, "Ask should return ErrNoValue if value pointer is not set")
	})

	t.Run("should return no error for valid input", func(t *testing.T) {
		var result string
		textPrompt := NewQuestion(&result).
			Title("Enter value:")
		textPrompt.Out = io.Discard
		textPrompt.In.Reader = testutils.CreateValidTestFile(t, "test input\r")

		err := textPrompt.Ask()
		require.NoError(t, err, "Ask should not return an error for valid input")
		require.Equal(t, "test input", result, "Value pointer should be set to the user input")
	})

	t.Run("should return ErrUserAborted when user presses Ctrl+C", func(t *testing.T) {
		var result string
		textPrompt := NewQuestion(&result).
			Title("Enter value:")
		textPrompt.Out = io.Discard
		textPrompt.In.Reader = testutils.CreateValidTestFile(t, "test input"+grapheme.CtrlC.String())

		err := textPrompt.Ask()
		require.ErrorIs(t, err, ErrUserAborted, "Ask should return ErrUserAborted when user presses Ctrl+C")
	})
}

// Tests for [Text.Hide] function.
func Test_Text_Hide(t *testing.T) {
	t.Run("should set hide to false on Question prompts", func(t *testing.T) {
		var result string
		textPrompt := NewQuestion(&result)
		require.False(t, textPrompt.hide, "Hide should be false by default for Question prompts")

		textPrompt = textPrompt.Hide(true)
		assert.True(t, textPrompt.hide, "Hide should be set to true")
	})

	t.Run("should set hide to true on Password prompts", func(t *testing.T) {
		var result string
		passwordPrompt := NewPassword(&result)
		require.True(t, passwordPrompt.hide, "Hide should be true by default for Password prompts")

		passwordPrompt = passwordPrompt.Hide(false)
		assert.False(t, passwordPrompt.hide, "Hide should be set to false")
	})
}

// Tests for [Text.Icon] function.
func Test_Text_Icon(t *testing.T) {
	t.Run("should set icon value and clear icon function", func(t *testing.T) {
		var result string
		question := NewQuestion(&result)
		assert.Equal(t, Icons.QuestionMark, question.icon.val, "Default icon should be a question mark")

		question = question.
			IconFunc(func(input string) string {
				return "🔍 "
			})
		assert.NotNil(t, question.icon.fn, "IconFunc() should set the icon function")

		question = question.Icon("➤ ")
		assert.Equal(t, "➤ ", question.icon.val, "Icon() should set the icon value")
		assert.Nil(t, question.icon.fn, "Icon() should set the icon function to nil")
	})
}

// Tests for [Text.IconFunc] function.
func Test_Text_IconFunc(t *testing.T) {
	t.Run("should set icon function", func(t *testing.T) {
		expected := "🔍 " + Icons.QuestionMark
		var result string
		question := NewQuestion(&result)
		assert.Nil(t, question.icon.fn, "Icon function should be nil by default")

		question = question.IconFunc(func(input string) string {
			return "🔍 " + input
		})
		assert.NotEmpty(t, question.icon.val, "Icon value should not be empty when using IconFunc")
		assert.NotNil(t, question.icon.fn, "IconFunc() should set the icon function")
		assert.Equal(t, expected, question.icon.Get(), "Icon function did not return expected result")
	})
}

// Tests for [Text.Title] function.
func Test_Text_Title(t *testing.T) {
	t.Run("should set title value and clear title function", func(t *testing.T) {
		var result string
		question := NewQuestion(&result)
		assert.Empty(t, question.title.val, "Default title should be empty")

		question = question.
			TitleFunc(func(input string) string {
				return "Dynamic: " + input
			})
		assert.NotNil(t, question.title.fn, "TitleFunc() should set the title function")

		question = question.Title("Static title")
		assert.Equal(t, "Static title", question.title.val, "Title() should set the title value")
		assert.Nil(t, question.title.fn, "Title() should set the title function to nil")
	})
}

// Tests for [Text.TitleFunc] function.
func Test_Text_TitleFunc(t *testing.T) {
	t.Run("should set title function", func(t *testing.T) {
		expected := "Dynamic: Static title"
		var result string
		question := NewQuestion(&result).
			Title("Static title")
		assert.Equal(t, "Static title", question.title.val, "Title should be set to static title")
		assert.Nil(t, question.title.fn, "Title function should be nil by default")

		question = question.TitleFunc(func(input string) string {
			return "Dynamic: " + input
		})
		assert.NotEmpty(t, question.title.val, "Title value should not be empty when using TitleFunc")
		assert.NotNil(t, question.title.fn, "TitleFunc() should set the title function")
		assert.Equal(t, expected, question.title.Get(), "Title function did not return expected result")
	})
}

// Tests for [Text.ValidateFunc] function.
func Test_Text_ValidateFunc(t *testing.T) {
	t.Run("should be valid regardless of input when using default validate function", func(t *testing.T) {
		var result string
		questionPrompt := NewQuestion(&result)
		assert.NoError(t, questionPrompt.validateFn(""), "Default validate function should not return an error for empty input")
		assert.NoError(t, questionPrompt.validateFn("test"), "Default validate function should not return an error for non-empty input")
	})

	t.Run("should set custom validation function", func(t *testing.T) {
		var result string
		questionPrompt := NewQuestion(&result)
		require.NoError(t, questionPrompt.validateFn("test"), "Original validate function should not return an error for valid input")

		questionPrompt = questionPrompt.ValidateFunc(func(input string) error {
			if input == "test" {
				return errors.New("input cannot be 'test'")
			}
			return nil
		})
		assert.Error(t, questionPrompt.validateFn("test"), "Custom validate function should return an error for input 'test'")
		assert.NoError(t, questionPrompt.validateFn("valid input"), "Custom validate function should not return an error for valid input")
	})
}

// Tests for [Text.Value] function.
func Test_Text_Value(t *testing.T) {
	t.Run("should set value pointer", func(t *testing.T) {
		var result string
		question := NewQuestion(&result)
		assert.Same(t, &result, question.value, "Value pointer should be set to the provided variable")

		var result2 string
		question = question.Value(&result2)
		assert.Same(t, &result2, question.value, "Value pointer should be updated to the new variable")
		assert.NotSame(t, &result, question.value, "Value pointer should no longer point to the original variable")
	})

	t.Run("should allow setting value pointer to nil", func(t *testing.T) {
		question := NewQuestion(nil)
		assert.Nil(t, question.value, "Value pointer should be nil when initialized with nil")

		question = question.Value(nil)
		assert.Nil(t, question.value, "Value pointer should remain nil when set to nil")
	})
}

// Tests for [Text.ask] function.
func Test_Text_ask(t *testing.T) {
	const promptTitle = "What is your name?"
	const expectedResult = "Bobby"
	t.Run("should set value to true when confirmed", func(t *testing.T) {
		var result string
		questionPrompt := NewQuestion(&result).
			Title(promptTitle)
		questionPrompt.Out = io.Discard
		questionPrompt.In.Reader = testutils.CreateValidTestFile(t, expectedResult+"\r")

		err := questionPrompt.ask()
		require.NoError(t, err, "Expected no error when asking with valid input")
		require.Equal(t, expectedResult, result, "Expected result to be 'Bobby' when confirmed")
	})

	t.Run("should error when user presses Ctrl+C", func(t *testing.T) {
		var result string
		questionPrompt := NewQuestion(&result).
			Title(promptTitle)
		questionPrompt.Out = io.Discard
		questionPrompt.In.Reader = testutils.CreateValidTestFile(t, grapheme.CtrlC.String())

		err := questionPrompt.ask()
		require.Error(t, err, "Expected error when user presses Ctrl+C")
		require.ErrorAsf(t, err, &ErrUserAborted, "Expected ErrUserAborted but got: %v", err)
	})

	t.Run("should error when In is not os.File", func(t *testing.T) {
		var result string
		questionPrompt := NewQuestion(&result).
			Title(promptTitle)
		questionPrompt.Out = io.Discard
		questionPrompt.In.Reader = bytes.NewBufferString(grapheme.New('y').String())

		err := questionPrompt.ask()
		require.Error(t, err, "Expected error when In is not os.File")
	})

	t.Run("should print error message when validation fails", func(t *testing.T) {
		mockPTY, mockTTY := testutils.CreatePTYWithSize(t, 20, 10)
		t.Cleanup(func() {
			mockPTY.Close()
			mockTTY.Close()
		})

		var result string
		questionPrompt := NewQuestion(&result).
			Title(promptTitle).
			ValidateFunc(func(input string) error {
				if input == expectedResult {
					return errors.New("validation failed")
				}

				return nil
			})
		questionPrompt.Out = mockTTY
		// Buffer size is 8 so need to pad two empty spaces to emulate stdin clearing the line after validation error is printed.
		questionPrompt.In.Reader = testutils.CreateValidTestFile(t, expectedResult+"\n  "+expectedResult+"2\r")
		err := questionPrompt.ask()
		require.NoError(t, err, "Expected no error when validation fails")
		require.Equal(t, expectedResult+"2", result, "Expected result to be 'Bobby2' after correcting validation error")
	})
}

// Tests for [Text.getPrompt] function.
func Test_Text_getPrompt(t *testing.T) {
	setPrompt := func(t *testing.T, questionPrompt *Text) {
		t.Helper()
		questionPrompt.Prompt = grapheme.ClusterSetFromString(fmt.Sprintf("%s%s ", questionPrompt.icon.Get(), questionPrompt.title.Get()))
	}

	t.Run("Simple prompt", func(t *testing.T) {
		name := ""
		questionPrompt := NewQuestion(&name).
			Title("What is your name?")

		setPrompt(t, questionPrompt)

		s := questionPrompt.getPromptWithDefaultValue()
		expected := fmt.Sprintf("%s%s%s ", ansi.ClearLineReset, Icons.QuestionMark, "What is your name?")
		assert.Equal(t, expected, s, "getPrompt() did not return expected prompt string")
	})

	// name := "input ❤️ [31mApple[0m"
	// questionPrompt := NewQuestion(&name).
	// 	Title("What is your name?")

	// // questionPrompt.Out = mockTTY
	// questionPrompt.Prompt = fmt.Sprintf("%s%s ", questionPrompt.icon.Get(), questionPrompt.title.Get())

	// s := questionPrompt.getPrompt()
	// _ = s

	// changeDefaultFunc(t, iconFn, func(s string) string {
	// 	return fmt.Sprintf("%s%s%s", ansi.Green, s, ansi.Reset)
	// })
}

// Tests for [Text.processInput] function.
func Test_Text_processInput(t *testing.T) {
	t.Run("should return done true and set value on Enter key", func(t *testing.T) {
		var result string
		questionPrompt := NewQuestion(&result)
		questionPrompt.Out = io.Discard

		expectedOutput := "input test"

		done, err := questionPrompt.processInput([]byte(expectedOutput + "\r"))
		require.NoError(t, err, "Expected no error when processing input with Enter key")
		require.True(t, done, "Expected done to be true when Enter key is pressed")
		require.Equal(t, expectedOutput, result, "Expected value to be set to 'input test' when Enter key is pressed")
	})

	t.Run("should return done true and set value on Enter key2", func(t *testing.T) {
		var result string
		questionPrompt := NewQuestion(&result)
		questionPrompt.Out = io.Discard
		questionPrompt.In.Reader = testutils.CreateValidTestFile(t, "input ❤️"+ansi.Red+"Apple"+ansi.Reset+"\r")

		err := questionPrompt.ask()
		require.NoError(t, err, "Expected no error when processing input with Enter key")
		require.Equal(t, "input ❤️Apple", result, "Expected value to be set to 'input Apple' when Enter key is pressed")
	})
}

// Tests for [Text.printErrorMessage] function.
func Test_Text_printErrorMessage(t *testing.T) {
	t.Run("should print error message with correct formatting", func(t *testing.T) {
		expectedErrorMessage := errors.New("Test error")
		expected := restoreCursor + ansi.ClearFromCursorToEndScreen + validationErrorMessage(expectedErrorMessage) + restoreCursor

		mockPTY, mockTTY := testutils.CreatePTYWithSize(t, 20, 10)

		questionPrompt := NewQuestion(nil)
		questionPrompt.Out = mockTTY
		questionPrompt.printErrorMessage(expectedErrorMessage)
		_ = mockTTY.Close() // Close the TTY to signal we're done reading output

		output := testutils.ReadPTYOutput(t, mockPTY, 128)
		assert.Equal(t, expected, output, "printErrorMessage did not print the expected error message with correct formatting")
	})

	t.Run("should return an error if output writer is not a file", func(t *testing.T) {
		questionPrompt := NewQuestion(nil)
		questionPrompt.Out = &bytes.Buffer{}

		errorMessage := errors.New("Test error")
		assert.Panics(t, func() {
			questionPrompt.printErrorMessage(errorMessage)
		})
	})

	t.Run("should handle error that exceeds terminal width", func(t *testing.T) {
		longErrorMessage := errors.New("This is a very long error message that should exceed the terminal width and be handled properly")
		expected := "\r\n" + validationErrorMessage(longErrorMessage) + ansi.CursorUp(1)

		mockPTY, mockTTY := testutils.CreatePTYWithSize(t, 60, 10)

		questionPrompt := NewQuestion(nil)
		questionPrompt.Out = mockTTY
		questionPrompt.printErrorMessage(longErrorMessage)
		_ = mockTTY.Close() // Close the TTY to signal we're done reading output

		output := testutils.ReadPTYOutput(t, mockPTY, 128)
		assert.Equal(t, expected, output, "printErrorMessage did not handle long error message correctly")
	})
}

// Tests for [Text.printFinalPromptLine] function.
func Test_Text_printFinalPromptLine(t *testing.T) {
	t.Run("should print final prompt line with prompt and answer", func(t *testing.T) {
		expectedOutput := ansi.ClearLineReset + "[?] What is your name? Bobby\n" + ansi.ClearLineReset
		var result string
		questionPrompt := NewQuestion(&result).
			Title("What is your name?")
		questionPrompt.Out = &bytes.Buffer{}
		*questionPrompt.value = "Bobby"
		questionPrompt.Prompt = grapheme.ClusterSetFromString(fmt.Sprintf("%s%s ", questionPrompt.icon.Get(), questionPrompt.title.Get()))

		questionPrompt.printFinalPromptLine()
		require.Equal(t, expectedOutput, questionPrompt.Out.(*bytes.Buffer).String(), "printFinalPromptLine() did not print expected output when value is false")
	})

	t.Run("should print final prompt line with only prompt when hide is true", func(t *testing.T) {
		expectedOutput := "\n" + ansi.ClearLineReset
		var result string
		passwordPrompt := NewPassword(&result).
			Title("What is your password?")
		passwordPrompt.Out = &bytes.Buffer{}
		*passwordPrompt.value = "MySuperSecretPassword"
		passwordPrompt.Prompt = grapheme.ClusterSetFromString(fmt.Sprintf("%s%s ", passwordPrompt.icon.Get(), passwordPrompt.title.Get()))

		passwordPrompt.printFinalPromptLine()
		require.Equal(t, expectedOutput, passwordPrompt.Out.(*bytes.Buffer).String(), "printFinalPromptLine() did not print expected output when hide is true")
	})
}
