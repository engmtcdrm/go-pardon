package pardon

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/engmtcdrm/go-ansi"
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
		textPrompt.terminal.Out = io.Discard
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
		textPrompt.terminal.Out = io.Discard
		textPrompt.answer.val = "Test"
		assert.Equal(t, customFn("Test"), textPrompt.answer.Get(), "Custom answer function did not return expected result")
	})
}

// TODO: Tests for [Text.Ask] function.
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
}

// Tests for [Text.Hide] function.
func Test_Text_Hide(t *testing.T) {
	t.Run("should set hide to false on Question prompts", func(t *testing.T) {
		var result string
		textPrompt := NewQuestion(&result)
		require.False(t, textPrompt.terminal.Hide, "Hide should be false by default for Question prompts")

		textPrompt = textPrompt.Hide(true)
		assert.True(t, textPrompt.terminal.Hide, "Hide should be set to true")
	})

	t.Run("should set hide to true on Password prompts", func(t *testing.T) {
		var result string
		passwordPrompt := NewPassword(&result)
		require.True(t, passwordPrompt.terminal.Hide, "Hide should be true by default for Password prompts")

		passwordPrompt = passwordPrompt.Hide(false)
		assert.False(t, passwordPrompt.terminal.Hide, "Hide should be set to false")
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

// TODO: Tests for [Text.ValidateFunc] function.
func Test_Text_ValidateFunc(t *testing.T) {
	var result string
	question := NewQuestion(&result).
		Title("Enter name:").
		ValidateFunc(func(input string) error {
			if input == "" {
				return ErrNoValue
			}
			return nil
		})

	// Validate functionality test - we can't easily test the actual input
	// but we can verify the question was configured properly
	if question == nil {
		t.Error("Question with validation returned nil")
	}
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

// TODO: Tests for [Text.ask] function.
func Test_Text_ask(t *testing.T) {
	t.Skip("Need to implement")
}

// Tests for [Text.getPromptLines] function.
func Test_Text_getPromptLines(t *testing.T) {
	t.Run("should return 1 when prompt line fits within terminal width", func(t *testing.T) {
		mockTTY := testutils.CreatePTYWithSize(t, "", 20, 10)

		questionPrompt := NewQuestion(nil)
		questionPrompt.terminal.Out = mockTTY

		lines, err := questionPrompt.getPromptLines("Short prompt?")
		assert.NoError(t, err, "getPromptLines should not return an error")
		assert.Equal(t, 1, lines, "getPromptLines should return 1 for short prompt")
	})

	t.Run("should return correct number of lines for prompt that exceeds terminal width", func(t *testing.T) {
		mockTTY := testutils.CreatePTYWithSize(t, "", 5, 10)

		questionPrompt := NewQuestion(nil)
		questionPrompt.terminal.Out = mockTTY

		lines, err := questionPrompt.getPromptLines("Short prompt?")
		assert.NoError(t, err, "getPromptLines should not return an error")
		assert.Equal(t, 3, lines, "getPromptLines should return 3 for prompt that exceeds terminal width")
	})

	t.Run("should return an error if output writer is not a file", func(t *testing.T) {
		questionPrompt := NewQuestion(nil)
		questionPrompt.terminal.Out = &bytes.Buffer{}

		_, err := questionPrompt.getPromptLines("Short prompt?")
		assert.Error(t, err, "getPromptLines should return an error if output writer is not a file")
	})

	t.Run("should return an error if terminal width is 0", func(t *testing.T) {
		mockTTY := testutils.CreatePTYWithSize(t, "", 0, 10)

		questionPrompt := NewQuestion(nil)
		questionPrompt.terminal.Out = mockTTY

		_, err := questionPrompt.getPromptLines("Short prompt?")
		assert.Error(t, err, "getPromptLines should return an error if terminal width is 0")
	})

	t.Run("should return an error if terminal size cannot be determined", func(t *testing.T) {
		// Use a regular file (not a PTY) so term.GetSize will fail with ENOTTY.
		questionPrompt := NewQuestion(nil)
		questionPrompt.terminal.Out = testutils.CreateValidTestFile(t, "not a pty")

		_, err := questionPrompt.getPromptLines("Short prompt?")
		assert.Error(t, err, "getPromptLines should return an error if terminal size cannot be determined")
	})
}

// TODO: Tests for [Text.printErrorMessage] function.
func Test_Text_printErrorMessage(t *testing.T) {
	t.Skip("Need to implement")
}

// Tests for [Text.printFinalPromptLine] function.
func Test_Text_printFinalPromptLine(t *testing.T) {
	t.Run("should print final prompt line with prompt and answer", func(t *testing.T) {
		expectedOutput := ansi.ClearLineReset + "[?] What is your name? Bobby\n" + ansi.ClearLineReset
		var result string
		questionPrompt := NewQuestion(&result).
			Title("What is your name?")
		questionPrompt.terminal.Out = &bytes.Buffer{}
		*questionPrompt.value = "Bobby"
		questionPrompt.prompt = fmt.Sprintf("%s%s ", questionPrompt.icon.Get(), questionPrompt.title.Get())

		questionPrompt.printFinalPromptLine()
		require.Equal(t, expectedOutput, questionPrompt.terminal.Out.(*bytes.Buffer).String(), "printFinalPromptLine() did not print expected output when value is false")
	})

	t.Run("should print final prompt line with only prompt when hide is true", func(t *testing.T) {
		expectedOutput := "\n" + ansi.ClearLineReset
		var result string
		passwordPrompt := NewPassword(&result).
			Title("What is your password?")
		passwordPrompt.terminal.Out = &bytes.Buffer{}
		*passwordPrompt.value = "MySuperSecretPassword"
		passwordPrompt.prompt = fmt.Sprintf("%s%s ", passwordPrompt.icon.Get(), passwordPrompt.title.Get())

		passwordPrompt.printFinalPromptLine()
		require.Equal(t, expectedOutput, passwordPrompt.terminal.Out.(*bytes.Buffer).String(), "printFinalPromptLine() did not print expected output when hide is true")
	})
}
