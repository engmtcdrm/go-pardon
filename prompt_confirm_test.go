package pardon

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/engmtcdrm/go-ansi"
	"github.com/engmtcdrm/go-pardon/internal/keys"
	"github.com/engmtcdrm/go-pardon/internal/runekeys"
	"github.com/engmtcdrm/go-pardon/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for [NewConfirm] function.
func Test_NewConfirm(t *testing.T) {
	t.Run("with default settings", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		require.NotNil(t, confirm, "NewConfirm returned nil")
	})
}

// Tests for [Confirm.AnswerFunc] function.
func Test_Confirm_AnswerFunc(t *testing.T) {
	t.Run("using default answer function", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			Title("Continue?")
		confirm.terminal.Out = io.Discard
		assert.Nil(t, confirm.answerFn, "Default answer function should be nil")
	})

	t.Run("using custom answer function", func(t *testing.T) {
		var result bool
		customFn := func(s string) string {
			return "Custom: " + s
		}
		confirm := NewConfirm(&result).
			Title("Continue?").
			AnswerFunc(customFn)
		confirm.terminal.Out = io.Discard
		assert.Equal(t, customFn("Test"), confirm.answerFn("Test"), "Custom answer function did not return expected result")
	})
}

// Tests for [Confirm.Ask] function.
func Test_Confirm_Ask(t *testing.T) {
	t.Run("should return error with no title", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)

		err := confirm.Ask()
		require.Error(t, err, "Expected error when asking without title")
		require.ErrorAsf(t, err, &ErrNoTitle, "Expected ErrNoTitle but got: %v", err)
	})

	t.Run("should return error with nil value", func(t *testing.T) {
		confirm := NewConfirm(nil).
			Title("Continue?")
		confirm.terminal.Out = io.Discard

		err := confirm.Ask()
		require.Error(t, err, "Expected error when asking with nil value")
		require.ErrorAsf(t, err, &ErrNoValue, "Expected ErrNoValue but got: %v", err)
	})

	t.Run("should set value to true when confirmed", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			Title("Continue?")
		confirm.terminal.Out = io.Discard
		confirm.terminal.In = testutils.CreateValidTestFile(t, string(keys.LowerY))

		err := confirm.Ask()
		require.NoError(t, err, "Expected no error when asking with valid input")
		require.True(t, result, "Expected result to be true when confirmed")
	})

	t.Run("should error when user presses Ctrl+C", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			Title("Continue?")
		confirm.terminal.Out = io.Discard
		confirm.terminal.In = testutils.CreateValidTestFile(t, string(keys.CtrlC))

		err := confirm.Ask()
		require.Error(t, err, "Expected error when user presses Ctrl+C")
		require.ErrorAsf(t, err, &ErrUserAborted, "Expected ErrUserAborted but got: %v", err)
	})
}

// Tests for [Confirm.ConfirmKey] function.
func Test_Confirm_ConfirmKey(t *testing.T) {
	t.Run("default confirm key", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		require.Equal(t, runekeys.UpperY, confirm.confirmKey, "Default confirm key should be 'Y'")
	})

	t.Run("custom confirm key", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			ConfirmKey(runekeys.UpperO)
		require.Equal(t, runekeys.UpperO, confirm.confirmKey, "Custom confirm key should be 'O'")
	})
}

// Tests for [Confirm.DenyKey] function.
func Test_Confirm_DenyKey(t *testing.T) {
	t.Run("default deny key", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		require.Equal(t, runekeys.UpperN, confirm.denyKey, "Default deny key should be 'N'")
	})

	t.Run("custom deny key", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			DenyKey(runekeys.UpperO)
		require.Equal(t, runekeys.UpperO, confirm.denyKey, "Custom deny key should be 'O'")
	})
}

// Tests for [Confirm.Icon] function.
func Test_Confirm_Icon(t *testing.T) {
	t.Run("with icon", func(t *testing.T) {
		expectedIcon := "[?]"

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, iconFn, nil)

		var result bool
		confirm := NewConfirm(&result).
			Icon(expectedIcon)
		confirm.icon.defaultFn = nil
		require.Equal(t, expectedIcon, confirm.icon.val, "Icon() did not set the icon correctly")
		require.Equal(t, expectedIcon, confirm.icon.Get(), "Icon.Get() did not return the expected icon")
	})

	t.Run("without icon", func(t *testing.T) {
		expectedIcon := ""

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, iconFn, nil)

		var result bool
		confirm := NewConfirm(&result).
			Icon("")
		require.Equal(t, expectedIcon, confirm.icon.val, "Icon() did not set the icon correctly")
		require.Equal(t, expectedIcon, confirm.icon.Get(), "Icon.Get() did not return the expected icon")
	})
}

// Tests for [Confirm.IconFunc] function.
func Test_Confirm_IconFunc(t *testing.T) {
	t.Run("with icon function", func(t *testing.T) {
		expectedIcon := "dynamic icon: test"
		var result bool

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, iconFn, nil)

		confirm := NewConfirm(&result).
			Icon("test").
			IconFunc(func(s string) string {
				return "dynamic icon: " + s
			})
		require.NotNil(t, confirm.icon.fn, "Icon function should not be nil")
		assert.Equal(t, expectedIcon, confirm.icon.Get(), "Icon function did not return expected result")
	})

	t.Run("with icon function being nil", func(t *testing.T) {
		expectedIcon := "[?]"
		var result bool

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, iconFn, nil)

		confirm := NewConfirm(&result).
			Icon(expectedIcon).
			IconFunc(nil)
		require.Nil(t, confirm.icon.fn, "Icon function should be nil")
		assert.Equal(t, expectedIcon, confirm.icon.Get(), "Icon function should return input when nil")
	})
}

// Tests for [Confirm.Title] function.
func Test_Confirm_Title(t *testing.T) {
	t.Run("with title", func(t *testing.T) {
		expectedTitle := "Are you sure?"

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, titleFn, nil)

		var result bool
		confirm := NewConfirm(&result).
			Title(expectedTitle)
		require.Equal(t, expectedTitle, confirm.title.val, "Title() did not set the title correctly")
		require.Equal(t, expectedTitle, confirm.title.Get(), "Title.Get() did not return the expected title")
	})

	t.Run("without title", func(t *testing.T) {
		expectedTitle := ""

		// Need to set this to nil so we can compare the exact value without transformation
		changeDefaultFunc(t, titleFn, nil)

		var result bool
		confirm := NewConfirm(&result)
		require.Equal(t, expectedTitle, confirm.title.val, "Title() did not set the title correctly")
		require.Equal(t, expectedTitle, confirm.title.Get(), "Title.Get() did not return the expected title")
	})
}

// Tests for [Confirm.TitleFunc] function.
func Test_Confirm_TitleFunc(t *testing.T) {
	t.Run("with title function", func(t *testing.T) {
		expectedTitle := "dynamic title: test"
		var result bool
		confirm := NewConfirm(&result).
			Title("test").
			TitleFunc(func(s string) string {
				return "dynamic title: " + s
			})
		require.NotNil(t, confirm.title.fn, "Title function should not be nil")
		assert.Equal(t, expectedTitle, confirm.title.Get(), "Title function did not return expected result")
	})

	t.Run("with title function being nil", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			Title("test").
			TitleFunc(nil)
		require.Nil(t, confirm.title.fn, "Title function should be nil")
		assert.Equal(t, "test", confirm.title.Get(), "Title function should return input when nil")
	})
}

// Tests for [Confirm.Value] function.
func Test_Confirm_Value(t *testing.T) {
	t.Run("setting value to true", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		require.Equal(t, false, *confirm.value, "Initial value should be false")

		result2 := true
		confirm = confirm.Value(&result2)
		require.Equal(t, true, *confirm.value, "Value() did not set the value to true")
	})
}

// Tests for [Confirm.ask] function.
func Test_Confirm_ask(t *testing.T) {
	t.Run("should set value to true when confirmed", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			Title("Continue?")
		confirm.terminal.Out = io.Discard
		confirm.terminal.In = testutils.CreateValidTestFile(t, string(keys.LowerY))

		err := confirm.ask()
		require.NoError(t, err, "Expected no error when asking with valid input")
		require.True(t, result, "Expected result to be true when confirmed")
	})

	t.Run("should error when user presses Ctrl+C", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			Title("Continue?")
		confirm.terminal.Out = io.Discard
		confirm.terminal.In = testutils.CreateValidTestFile(t, string(keys.CtrlC))

		err := confirm.ask()
		require.Error(t, err, "Expected error when user presses Ctrl+C")
		require.ErrorAsf(t, err, &ErrUserAborted, "Expected ErrUserAborted but got: %v", err)
	})

	t.Run("should error when In is not os.File", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			Title("Continue?")
		confirm.terminal.Out = io.Discard
		confirm.terminal.In = bytes.NewBufferString(string(keys.LowerY))

		err := confirm.ask()
		require.Error(t, err, "Expected error when In is not os.File")
	})
}

// Tests for [Confirm.callAnswerFunc] function.
func Test_Confirm_callAnswerFunc(t *testing.T) {
	var result bool
	confirmPrompt := NewConfirm(&result)
	t.Run("should return the string value when no answer function or default answer function is set", func(t *testing.T) {
		originalDefaultAnswerFunc := defaultFuncs.answerFn
		t.Cleanup(func() { defaultFuncs.answerFn = originalDefaultAnswerFunc })
		defaultFuncs.answerFn = nil

		expectedOutput := "Test Answer"
		assert.Equal(t, expectedOutput, confirmPrompt.callAnswerFunc(expectedOutput), "getAnswerFunc() should return the input string when no functions are set")
	})

	t.Run("should return the string transformed by default answer function if set and no prompt-specific function is set", func(t *testing.T) {
		expectedOutput := "Test Answer"
		assert.Equal(t, expectedOutput, confirmPrompt.callAnswerFunc(expectedOutput), "getAnswerFunc() should return the string transformed by the default answer function when no prompt-specific function is set")
	})

	t.Run("should return the string transformed by custom default answer function if set and no prompt-specific function is set", func(t *testing.T) {
		originalDefaultAnswerFunc := defaultFuncs.answerFn
		t.Cleanup(func() { defaultFuncs.answerFn = originalDefaultAnswerFunc })
		defaultFuncs.answerFn = func(input string) string {
			return "[ANSWER: " + input + "]"
		}

		expectedOutput := "[ANSWER: Test Answer]"
		assert.Equal(t, expectedOutput, confirmPrompt.callAnswerFunc("Test Answer"), "getAnswerFunc() should return the string transformed by the custom default answer function when no prompt-specific function is set")
	})

	t.Run("should return the string transformed by the prompt-specific answer function if set", func(t *testing.T) {
		answerFunc := func(input string) string {
			return "[ANSWER: " + input + "]"
		}
		confirmPrompt.AnswerFunc(answerFunc)

		expectedOutput := "[ANSWER: Test Answer]"
		assert.Equal(t, expectedOutput, confirmPrompt.callAnswerFunc("Test Answer"), "getAnswerFunc() should return the string transformed by the prompt-specific answer function when it is set")
	})
}

// Tests for [Confirm.equal] function.
func Test_Confirm_equal(t *testing.T) {
	t.Run("should return true when values match", func(t *testing.T) {
		confirm := NewConfirm(nil)
		require.True(t, confirm.equal(runekeys.UpperY, runekeys.UpperY), "Expected equal to return true")
		require.True(t, confirm.equal(runekeys.UpperN, runekeys.UpperN), "Expected equal to return true")
	})

	t.Run("should return false when values do not match", func(t *testing.T) {
		confirm := NewConfirm(nil)

		require.False(t, confirm.equal(runekeys.UpperO, runekeys.UpperY), "Expected equal to return false")
		require.False(t, confirm.equal(runekeys.UpperO, runekeys.UpperN), "Expected equal to return false")
	})
}

// Tests for [Confirm.getPromptOptions] function.
func Test_Confirm_getPromptOptions(t *testing.T) {
	t.Run("with default confirm and deny keys", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		expectedOptions := "[y/N]"
		require.Equal(t, expectedOptions, confirm.getPromptOptions(), "getPromptOptions() did not return expected options")

		result = true
		expectedOptions = "[Y/n]"
		require.Equal(t, expectedOptions, confirm.getPromptOptions(), "getPromptOptions() did not return expected options when value is true")
	})

	t.Run("with custom confirm and deny keys", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			ConfirmKey(runekeys.UpperO).
			DenyKey(runekeys.UpperA)
		expectedOptions := "[o/A]"
		require.Equal(t, expectedOptions, confirm.getPromptOptions(), "getPromptOptions() did not return expected options with custom keys")

		result = true
		expectedOptions = "[O/a]"
		require.Equal(t, expectedOptions, confirm.getPromptOptions(), "getPromptOptions() did not return expected options with custom keys when value is true")
	})
}

// Tests for [Confirm.getValueAsRunes] function.
func Test_Confirm_getValueAsRunes(t *testing.T) {
	t.Run("with default confirm and deny keys", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		require.Equal(t, []rune{confirm.denyKey}, confirm.getValueAsRunes(), "getValueAsRunes() did not return expected runes when value is false")

		result = true
		require.Equal(t, []rune{confirm.confirmKey}, confirm.getValueAsRunes(), "getValueAsRunes() did not return expected runes when value is true")
	})

	t.Run("with custom confirm and deny keys", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			ConfirmKey(runekeys.UpperO).
			DenyKey(runekeys.UpperA)
		require.Equal(t, []rune{confirm.denyKey}, confirm.getValueAsRunes(), "getValueAsRunes() did not return expected runes with custom keys when value is false")

		result = true
		require.Equal(t, []rune{confirm.confirmKey}, confirm.getValueAsRunes(), "getValueAsRunes() did not return expected runes with custom keys when value is true")
	})
}

// Tests for [Confirm.getValueAsString] function.
func Test_Confirm_getValueAsString(t *testing.T) {
	t.Run("with default confirm and deny keys", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		require.Equal(t, "N", confirm.getValueAsString(), "getValueAsString() did not return expected string when value is false")

		result = true
		require.Equal(t, "Y", confirm.getValueAsString(), "getValueAsString() did not return expected string when value is true")
	})

	t.Run("with custom confirm and deny keys", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			ConfirmKey(runekeys.UpperO).
			DenyKey(runekeys.UpperA)
		require.Equal(t, "A", confirm.getValueAsString(), "getValueAsString() did not return expected string with custom keys when value is false")

		result = true
		require.Equal(t, "O", confirm.getValueAsString(), "getValueAsString() did not return expected string with custom keys when value is true")
	})
}

// Tests for [Confirm.printFinalPromptLine] function.
func Test_Confirm_printFinalPromptLine(t *testing.T) {
	t.Run("should print final prompt line with prompt and answer", func(t *testing.T) {
		expectedOutput := ansi.ClearLineReset + "[?] Continue? N\n" + ansi.ClearLineReset
		var result bool
		confirm := NewConfirm(&result).
			Title("Continue?")
		confirm.terminal.Out = &bytes.Buffer{}
		confirm.prompt = fmt.Sprintf("%s%s ", confirm.icon.Get(), confirm.title.Get())
		confirm.printFinalPromptLine()
		require.Equal(t, expectedOutput, confirm.terminal.Out.(*bytes.Buffer).String(), "printFinalPromptLine() did not print expected output when value is false")
	})
}

// Tests for [Confirm.processLine] function.
func Test_Confirm_processLine(t *testing.T) {
	t.Run("should return false for empty input", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		require.False(t, confirm.processLine([]rune{}), "processLine() should return false for empty input")
	})

	t.Run("should return set value if user enters enter character", func(t *testing.T) {
		result := true
		confirm := NewConfirm(&result)
		require.True(t, *confirm.value, "Initial value should be true")

		require.True(t, confirm.processLine([]rune{runekeys.Enter}), "processLine() should return true when user hits enter key")
		require.True(t, *confirm.value, "Value should be set to true when user hits enter key")
	})

	t.Run("should return set value if user enters new line character", func(t *testing.T) {
		result := true
		confirm := NewConfirm(&result)
		require.True(t, *confirm.value, "Initial value should be true")

		require.True(t, confirm.processLine([]rune{runekeys.NewLine}), "processLine() should return true when user hits new line key")
		require.True(t, *confirm.value, "Value should be set to true when user hits new line key")
	})

	t.Run("should set value to true when user confirms", func(t *testing.T) {
		result := false
		confirm := NewConfirm(&result)
		require.False(t, *confirm.value, "Initial value should be false")

		require.True(t, confirm.processLine([]rune{confirm.confirmKey}), "processLine() should return true when user confirms")
		require.True(t, *confirm.value, "Value should be set to true when user confirms")
	})

	t.Run("should set value to false when user denies", func(t *testing.T) {
		result := true
		confirm := NewConfirm(&result)
		require.True(t, *confirm.value, "Initial value should be true")

		require.True(t, confirm.processLine([]rune{confirm.denyKey}), "processLine() should return true when user denies")
		require.False(t, *confirm.value, "Value should be set to false when user denies")
	})

	t.Run("should return false for unrecognized input", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result)
		require.False(t, confirm.processLine([]rune{runekeys.UpperO}), "processLine() should return false for unrecognized input")
	})
}
