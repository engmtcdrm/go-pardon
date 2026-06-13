package pardon

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/internal/keys"
	"github.com/engmtcdrm/go-pardon/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for [NewConfirm] function.
func Test_NewConfirm(t *testing.T) {
	t.Run("with default settings", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		require.NotNil(t, confirmPrompt, "NewConfirm returned nil")
	})
}

// Tests for [Confirm.AnswerFunc] function.
func Test_Confirm_AnswerFunc(t *testing.T) {
	t.Run("using default answer function", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("Continue?")
		confirmPrompt.terminal.Out = io.Discard
		assert.Nil(t, confirmPrompt.answer.fn, "Default answer function should be nil")
	})

	t.Run("using custom answer function", func(t *testing.T) {
		var result bool
		customFn := func(s string) string {
			return "Custom: " + s
		}
		confirmPrompt := NewConfirm(&result).
			Title("Continue?").
			AnswerFunc(customFn)
		confirmPrompt.terminal.Out = io.Discard
		confirmPrompt.answer.val = "Test"
		assert.Equal(t, customFn("Test"), confirmPrompt.answer.Get(), "Custom answer function did not return expected result")
	})
}

// Tests for [Confirm.Ask] function.
func Test_Confirm_Ask(t *testing.T) {
	t.Run("should return error with no title", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)

		err := confirmPrompt.Ask()
		require.Error(t, err, "Expected error when asking without title")
		require.ErrorAsf(t, err, &ErrNoTitle, "Expected ErrNoTitle but got: %v", err)
	})

	t.Run("should return error with nil value", func(t *testing.T) {
		confirmPrompt := NewConfirm(nil).
			Title("Continue?")
		confirmPrompt.terminal.Out = io.Discard

		err := confirmPrompt.Ask()
		require.Error(t, err, "Expected error when asking with nil value")
		require.ErrorAsf(t, err, &ErrNoValue, "Expected ErrNoValue but got: %v", err)
	})

	t.Run("should set value to true when confirmed", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("Continue?")
		confirmPrompt.terminal.Out = io.Discard
		confirmPrompt.terminal.In = testutils.CreateValidTestFile(t, string(keys.LowerY))

		err := confirmPrompt.Ask()
		require.NoError(t, err, "Expected no error when asking with valid input")
		require.True(t, result, "Expected result to be true when confirmed")
	})

	t.Run("should error when user presses Ctrl+C", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("Continue?")
		confirmPrompt.terminal.Out = io.Discard
		confirmPrompt.terminal.In = testutils.CreateValidTestFile(t, string(keys.CtrlC))

		err := confirmPrompt.Ask()
		require.Error(t, err, "Expected error when user presses Ctrl+C")
		require.ErrorAsf(t, err, &ErrUserAborted, "Expected ErrUserAborted but got: %v", err)
	})
}

// Tests for [Confirm.ConfirmKey] function.
func Test_Confirm_ConfirmKey(t *testing.T) {
	t.Run("default confirm key", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		require.True(t, bytes.Equal(confirmPrompt.confirmKey, []byte{keys.UpperY}), "Default confirm key should be 'Y'")
	})

	t.Run("custom confirm key", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			ConfirmKey(keys.UpperO)
		require.True(t, bytes.Equal(confirmPrompt.confirmKey, []byte{keys.UpperO}), "Custom confirm key should be 'O'")
	})
}

// Tests for [Confirm.DenyKey] function.
func Test_Confirm_DenyKey(t *testing.T) {
	t.Run("default deny key", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		require.True(t, bytes.Equal(confirmPrompt.denyKey, []byte{keys.UpperN}), "Default deny key should be 'N'")
	})

	t.Run("custom deny key", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			DenyKey(keys.UpperO)
		require.True(t, bytes.Equal(confirmPrompt.denyKey, []byte{keys.UpperO}), "Custom deny key should be 'O'")
	})
}

// Tests for [Confirm.Icon] function.
func Test_Confirm_Icon(t *testing.T) {
	t.Run("with icon", func(t *testing.T) {
		expectedIcon := "[?]"

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, iconFn, nil)

		var result bool
		confirmPrompt := NewConfirm(&result).
			Icon(expectedIcon)
		confirmPrompt.icon.defaultFn = nil
		require.Equal(t, expectedIcon, confirmPrompt.icon.val, "Icon() did not set the icon correctly")
		require.Equal(t, expectedIcon, confirmPrompt.icon.Get(), "Icon.Get() did not return the expected icon")
	})

	t.Run("without icon", func(t *testing.T) {
		expectedIcon := ""

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, iconFn, nil)

		var result bool
		confirmPrompt := NewConfirm(&result).
			Icon("")
		require.Equal(t, expectedIcon, confirmPrompt.icon.val, "Icon() did not set the icon correctly")
		require.Equal(t, expectedIcon, confirmPrompt.icon.Get(), "Icon.Get() did not return the expected icon")
	})
}

// Tests for [Confirm.IconFunc] function.
func Test_Confirm_IconFunc(t *testing.T) {
	t.Run("with icon function", func(t *testing.T) {
		expectedIcon := "dynamic icon: test"
		var result bool

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, iconFn, nil)

		confirmPrompt := NewConfirm(&result).
			Icon("test").
			IconFunc(func(s string) string {
				return "dynamic icon: " + s
			})
		require.NotNil(t, confirmPrompt.icon.fn, "Icon function should not be nil")
		assert.Equal(t, expectedIcon, confirmPrompt.icon.Get(), "Icon function did not return expected result")
	})

	t.Run("with icon function being nil", func(t *testing.T) {
		expectedIcon := "[?]"
		var result bool

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, iconFn, nil)

		confirmPrompt := NewConfirm(&result).
			Icon(expectedIcon).
			IconFunc(nil)
		require.Nil(t, confirmPrompt.icon.fn, "Icon function should be nil")
		assert.Equal(t, expectedIcon, confirmPrompt.icon.Get(), "Icon function should return input when nil")
	})
}

// Tests for [Confirm.Title] function.
func Test_Confirm_Title(t *testing.T) {
	t.Run("with title", func(t *testing.T) {
		expectedTitle := "Are you sure?"

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, titleFn, nil)

		var result bool
		confirmPrompt := NewConfirm(&result).
			Title(expectedTitle)
		require.Equal(t, expectedTitle, confirmPrompt.title.val, "Title() did not set the title correctly")
		require.Equal(t, expectedTitle, confirmPrompt.title.Get(), "Title.Get() did not return the expected title")
	})

	t.Run("without title", func(t *testing.T) {
		expectedTitle := ""

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, titleFn, nil)

		var result bool
		confirmPrompt := NewConfirm(&result)
		require.Equal(t, expectedTitle, confirmPrompt.title.val, "Title() did not set the title correctly")
		require.Equal(t, expectedTitle, confirmPrompt.title.Get(), "Title.Get() did not return the expected title")
	})
}

// Tests for [Confirm.TitleFunc] function.
func Test_Confirm_TitleFunc(t *testing.T) {
	t.Run("with title function", func(t *testing.T) {
		expectedTitle := "dynamic title: test"
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("test").
			TitleFunc(func(s string) string {
				return "dynamic title: " + s
			})
		require.NotNil(t, confirmPrompt.title.fn, "Title function should not be nil")
		assert.Equal(t, expectedTitle, confirmPrompt.title.Get(), "Title function did not return expected result")
	})

	t.Run("with title function being nil", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("test").
			TitleFunc(nil)
		require.Nil(t, confirmPrompt.title.fn, "Title function should be nil")
		assert.Equal(t, "test", confirmPrompt.title.Get(), "Title function should return input when nil")
	})
}

// Tests for [Confirm.Value] function.
func Test_Confirm_Value(t *testing.T) {
	t.Run("setting value to true", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		require.Equal(t, false, *confirmPrompt.value, "Initial value should be false")

		result2 := true
		confirmPrompt = confirmPrompt.Value(&result2)
		require.Equal(t, true, *confirmPrompt.value, "Value() did not set the value to true")
	})
}

// Tests for [Confirm.ask] function.
func Test_Confirm_ask(t *testing.T) {
	t.Run("should set value to true when confirmed", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("Continue?")
		confirmPrompt.terminal.Out = io.Discard
		confirmPrompt.terminal.In = testutils.CreateValidTestFile(t, string(keys.LowerY))

		err := confirmPrompt.ask()
		require.NoError(t, err, "Expected no error when asking with valid input")
		require.True(t, result, "Expected result to be true when confirmed")
	})

	t.Run("should error when user presses Ctrl+C", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("Continue?")
		confirmPrompt.terminal.Out = io.Discard
		confirmPrompt.terminal.In = testutils.CreateValidTestFile(t, string(keys.CtrlC))

		err := confirmPrompt.ask()
		require.Error(t, err, "Expected error when user presses Ctrl+C")
		require.ErrorAsf(t, err, &ErrUserAborted, "Expected ErrUserAborted but got: %v", err)
	})

	t.Run("should error when In is not os.File", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("Continue?")
		confirmPrompt.terminal.Out = io.Discard
		confirmPrompt.terminal.In = bytes.NewBufferString(string(keys.LowerY))

		err := confirmPrompt.ask()
		require.Error(t, err, "Expected error when In is not os.File")
	})
}

// Tests for [Confirm.getPromptOptions] function.
func Test_Confirm_getPromptOptions(t *testing.T) {
	t.Run("with default confirm and deny keys", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		expectedOptions := "[y/N]"
		require.Equal(t, expectedOptions, confirmPrompt.getPromptOptions(), "getPromptOptions() did not return expected options")

		result = true
		expectedOptions = "[Y/n]"
		require.Equal(t, expectedOptions, confirmPrompt.getPromptOptions(), "getPromptOptions() did not return expected options when value is true")
	})

	t.Run("with custom confirm and deny keys", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			ConfirmKey(keys.UpperO).
			DenyKey(keys.UpperA)
		expectedOptions := "[o/A]"
		require.Equal(t, expectedOptions, confirmPrompt.getPromptOptions(), "getPromptOptions() did not return expected options with custom keys")

		result = true
		expectedOptions = "[O/a]"
		require.Equal(t, expectedOptions, confirmPrompt.getPromptOptions(), "getPromptOptions() did not return expected options with custom keys when value is true")
	})
}

// Tests for [Confirm.getValueAsBytes] function.
func Test_Confirm_getValueAsBytes(t *testing.T) {
	t.Run("with default confirm and deny keys", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		require.True(t, bytes.Equal(confirmPrompt.denyKey, confirmPrompt.getValueAsBytes()), "getValueAsBytes() did not return expected bytes when value is false")

		result = true
		require.True(t, bytes.Equal(confirmPrompt.confirmKey, confirmPrompt.getValueAsBytes()), "getValueAsBytes() did not return expected bytes when value is true")
	})

	t.Run("with custom confirm and deny keys", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			ConfirmKey(keys.UpperO).
			DenyKey(keys.UpperA)
		require.True(t, bytes.Equal(confirmPrompt.denyKey, confirmPrompt.getValueAsBytes()), "getValueAsBytes() did not return expected bytes with custom keys when value is false")

		result = true
		require.True(t, bytes.Equal(confirmPrompt.confirmKey, confirmPrompt.getValueAsBytes()), "getValueAsBytes() did not return expected bytes with custom keys when value is true")
	})
}

// Tests for [Confirm.getValueAsString] function.
func Test_Confirm_getValueAsString(t *testing.T) {
	t.Run("with default confirm and deny keys", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		require.Equal(t, "N", confirmPrompt.getValueAsString(), "getValueAsString() did not return expected string when value is false")

		result = true
		require.Equal(t, "Y", confirmPrompt.getValueAsString(), "getValueAsString() did not return expected string when value is true")
	})

	t.Run("with custom confirm and deny keys", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result).
			ConfirmKey(keys.UpperO).
			DenyKey(keys.UpperA)
		require.Equal(t, "A", confirmPrompt.getValueAsString(), "getValueAsString() did not return expected string with custom keys when value is false")

		result = true
		require.Equal(t, "O", confirmPrompt.getValueAsString(), "getValueAsString() did not return expected string with custom keys when value is true")
	})
}

// Tests for [Confirm.printFinalPromptLine] function.
func Test_Confirm_printFinalPromptLine(t *testing.T) {
	t.Run("should print final prompt line with prompt and answer", func(t *testing.T) {
		expectedOutput := ansi.ClearLineReset + "[?] Continue? N\n" + ansi.ClearLineReset
		var result bool
		confirmPrompt := NewConfirm(&result).
			Title("Continue?")
		confirmPrompt.terminal.Out = &bytes.Buffer{}
		confirmPrompt.prompt = fmt.Sprintf("%s%s ", confirmPrompt.icon.Get(), confirmPrompt.title.Get())
		confirmPrompt.printFinalPromptLine()
		require.Equal(t, expectedOutput, confirmPrompt.terminal.Out.(*bytes.Buffer).String(), "printFinalPromptLine() did not print expected output when value is false")
	})
}

// Tests for [Confirm.processLine] function.
func Test_Confirm_processLine(t *testing.T) {
	t.Run("should return false for empty input", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		done, err := confirmPrompt.processInput([]byte{})
		require.False(t, done, "processLine() should return false for empty input")
		require.NoError(t, err, "processLine() should not return an error for empty input")
	})

	t.Run("should return set value if user enters enter character", func(t *testing.T) {
		result := true
		confirmPrompt := NewConfirm(&result)
		require.True(t, *confirmPrompt.value, "Initial value should be true")

		done, err := confirmPrompt.processInput([]byte{keys.Enter})
		require.True(t, done, "processLine() should return true when user hits enter key")
		require.NoError(t, err, "processLine() should not return an error when user hits enter key")
		require.True(t, *confirmPrompt.value, "Value should be set to true when user hits enter key")
	})

	t.Run("should return set value if user enters new line character", func(t *testing.T) {
		result := true
		confirmPrompt := NewConfirm(&result)
		require.True(t, *confirmPrompt.value, "Initial value should be true")

		done, err := confirmPrompt.processInput([]byte{keys.NewLine})
		require.True(t, done, "processLine() should return true when user hits new line key")
		require.NoError(t, err, "processLine() should not return an error when user hits new line key")
		require.True(t, *confirmPrompt.value, "Value should be set to true when user hits new line key")
	})

	t.Run("should set value to true when user confirms", func(t *testing.T) {
		result := false
		confirmPrompt := NewConfirm(&result)
		require.False(t, *confirmPrompt.value, "Initial value should be false")

		done, err := confirmPrompt.processInput(confirmPrompt.confirmKey)
		require.True(t, done, "processLine() should return true when user confirms")
		require.NoError(t, err, "processLine() should not return an error when user confirms")
		require.True(t, *confirmPrompt.value, "Value should be set to true when user confirms")
	})

	t.Run("should set value to false when user denies", func(t *testing.T) {
		result := true
		confirmPrompt := NewConfirm(&result)
		require.True(t, *confirmPrompt.value, "Initial value should be true")

		done, err := confirmPrompt.processInput(confirmPrompt.denyKey)
		require.True(t, done, "processLine() should return true when user denies")
		require.NoError(t, err, "processLine() should not return an error when user denies")
		require.False(t, *confirmPrompt.value, "Value should be set to false when user denies")
	})

	t.Run("should return false for unrecognized input", func(t *testing.T) {
		var result bool
		confirmPrompt := NewConfirm(&result)
		done, err := confirmPrompt.processInput([]byte{keys.UpperO})
		require.False(t, done, "processLine() should return false for unrecognized input")
		require.NoError(t, err, "processLine() should not return an error for unrecognized input")
	})
}
