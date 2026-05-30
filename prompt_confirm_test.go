package pardon

import (
	"testing"

	"github.com/engmtcdrm/go-pardon/internal/runekeys"
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
		assert.Equal(t, customFn("Test"), confirm.answerFn("Test"), "Custom answer function did not return expected result")
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
	t.Skip("not implemented")
}

// Tests for [Confirm.Ask] function.
func Test_Confirm_Ask(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.ask] function.
func Test_Confirm_ask(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.callAnswerFunc] function.
func Test_Confirm_callAnswerFunc(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.processLine] function.
func Test_Confirm_processLine(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.equal] function.
func Test_Confirm_equal(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.getPromptOptions] function.
func Test_Confirm_getPromptOptions(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.getValueAsRunes] function.
func Test_Confirm_getValueAsRunes(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.getValueAsString] function.
func Test_Confirm_getValueAsString(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.printFinalPromptLine] function.
func Test_Confirm_printFinalPromptLine(t *testing.T) {
	t.Skip("not implemented")
}

// Tests for [Confirm.trimSpace] function.
func Test_Confirm_trimSpace(t *testing.T) {
	t.Skip("not implemented")
}

type funcType int

const (
	answerFn funcType = iota
	cursorFn
	iconFn
	selectFn
	titleFn
)

func changeDefaultFunc(t *testing.T, fnType funcType, newFn func(string) string) {
	var originalFn func(string) string

	switch fnType {
	case answerFn:
		originalFn = defaultFuncs.answerFn
		defaultFuncs.answerFn = newFn
	case cursorFn:
		originalFn = defaultFuncs.cursorFn
		defaultFuncs.cursorFn = newFn
	case iconFn:
		originalFn = defaultFuncs.iconFn
		defaultFuncs.iconFn = newFn
	case selectFn:
		originalFn = defaultFuncs.selectFn
		defaultFuncs.selectFn = newFn
	case titleFn:
		originalFn = defaultFuncs.titleFn
		defaultFuncs.titleFn = newFn
	default:
		t.Fatalf("Unknown function type: %v", fnType)
	}

	t.Cleanup(func() {
		switch fnType {
		case answerFn:
			defaultFuncs.answerFn = originalFn
		case cursorFn:
			defaultFuncs.cursorFn = originalFn
		case iconFn:
			defaultFuncs.iconFn = originalFn
		case selectFn:
			defaultFuncs.selectFn = originalFn
		case titleFn:
			defaultFuncs.titleFn = originalFn
		}
	})
}
