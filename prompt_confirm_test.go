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
	// TODO: Fillout
}

// Tests for [Confirm.IconFunc] function.
func Test_Confirm_IconFunc(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.Title] function.
func Test_Confirm_Title(t *testing.T) {
	t.Run("with title", func(t *testing.T) {
		expectedTitle := "Are you sure?"
		var result bool
		confirm := NewConfirm(&result).
			Title(expectedTitle)
		require.Equal(t, expectedTitle, confirm.title.val, "Title() did not set the title correctly")
	})

	t.Run("without title", func(t *testing.T) {
		expectedTitle := ""
		var result bool
		confirm := NewConfirm(&result)
		require.Equal(t, expectedTitle, confirm.title.val, "Title() did not set the title correctly")
	})
}

// Tests for [Confirm.TitleFunc] function.
func Test_Confirm_TitleFunc(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.Value] function.
func Test_Confirm_Value(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.Ask] function.
func Test_Confirm_Ask(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.ask] function.
func Test_Confirm_ask(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.callAnswerFunc] function.
func Tests_Confirm_callAnswerFunc(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.processLine] function.
func Tests_Confirm_processLine(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.equal] function.
func Tests_Confirm_equal(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.getPromptOptions] function.
func Tests_Confirm_getPromptOptions(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.getValueAsRunes] function.
func Tests_Confirm_getValueAsRunes(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.getValueAsString] function.
func Tests_Confirm_getValueAsString(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.printFinalPromptLine] function.
func Tests_Confirm_printFinalPromptLine(t *testing.T) {
	// TODO: Fillout
}

// Tests for [Confirm.trimSpace] function.
func Tests_Confirm_trimSpace(t *testing.T) {
	// TODO: Fillout
}
