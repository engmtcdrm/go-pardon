package pardon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for [NewSelect] function.
func Test_NewSelect(t *testing.T) {
	t.Run("should return a non-nil select prompt", func(t *testing.T) {
		var result string
		selectPrompt := NewSelect(&result)
		require.NotNil(t, selectPrompt, "NewSelect() should not return nil")
		assert.Same(t, &result, selectPrompt.value, "NewSelect() should set the value pointer correctly")
	})

	t.Run("should return a non-nil select prompt when value is nil", func(t *testing.T) {
		selectPrompt := NewSelect[string](nil)
		require.NotNil(t, selectPrompt, "NewSelect() should not return nil even when value is nil")
		assert.Nil(t, selectPrompt.value, "NewSelect() should set the value pointer to nil when nil is passed")
	})
}

// Tests for [Select.AnswerFunc] function.
func Test_Select_AnswerFunc(t *testing.T) {
	t.Run("should set the answer function", func(t *testing.T) {
		var result string
		answerFunc := func(answer string) string {
			return "[SELECTED: " + answer + "]"
		}

		selectPrompt := NewSelect(&result).
			AnswerFunc(answerFunc)
		require.NotNil(t, selectPrompt.answer.fn, "AnswerFunc() should set the answer function")
		testAnswer := "Test Answer"
		assert.Equal(t, answerFunc(testAnswer), selectPrompt.answer.fn(testAnswer), "AnswerFunc() should return the same result as the provided function")
	})
}

// TODO: Tests for [Select.Ask] function.
func Test_Select_Ask(t *testing.T) {
	t.Skip("Need to implement")
}

// Tests for [Select.Cursor] function.
func Test_Select_Cursor(t *testing.T) {
	t.Run("should set the cursor value", func(t *testing.T) {
		expectedCursor := ">> "
		var result string
		selectPrompt := NewSelect(&result).
			Cursor(expectedCursor)
		require.Equal(t, expectedCursor, selectPrompt.cursor.Get(), "Cursor() should set the cursor value correctly")
	})
}

// Tests for [Select.CursorFunc] function.
func Test_Select_CursorFunc(t *testing.T) {
	t.Run("should set the cursor function", func(t *testing.T) {
		preCursor := "⭐ "
		var result string
		selectPrompt := NewSelect(&result).
			CursorFunc(func(input string) string {
				return preCursor + input
			})
		expectedCursor := preCursor + "> "
		require.NotNil(t, selectPrompt.cursor.fn, "CursorFunc should set the cursor function")
		assert.Equal(t, expectedCursor, selectPrompt.cursor.Get(), "CursorFunc should return the expected cursor")
	})

	t.Run("should return cursor value with function applied", func(t *testing.T) {
		preCursor := "⭐ "
		var result string
		selectPrompt := NewSelect(&result).
			Cursor("test").
			CursorFunc(func(input string) string {
				return preCursor + input
			})
		require.NotNil(t, selectPrompt.cursor.fn, "CursorFunc should set the cursor function")
		assert.Equal(t, preCursor+"test", selectPrompt.cursor.Get(), "CursorFunc should return the expected cursor")
	})
}

// Tests for [Select.Icon] function.
func Test_Select_Icon(t *testing.T) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
	}

	var result string
	selectPrompt := NewSelect(&result).
		Options(options...).
		Icon("🎯 ")

	if selectPrompt.icon.val != "🎯 " {
		t.Errorf("Icon() = %q; want %q", selectPrompt.icon.val, "🎯 ")
	}
}

// Tests for [Select.IconFunc] function.
func Test_Select_IconFunc(t *testing.T) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
	}

	var result string
	selectPrompt := NewSelect(&result).
		Options(options...).
		IconFunc(func(input string) string {
			return "🚀 "
		})

	if selectPrompt.icon.fn == nil {
		t.Error("IconFunc should set the icon function")
	}
}

func Test_Select_Options(t *testing.T) {
	t.Run("string type", func(t *testing.T) {
		options := []Option[string]{
			{Key: "First", Value: "first_value"},
			{Key: "Second", Value: "second_value"},
		}

		var result string
		selectPrompt := NewSelect(&result).
			Options(options...)
		require.Len(t, selectPrompt.options, len(options), "Options() should set the correct number of options")
	})

	t.Run("int type", func(t *testing.T) {
		options := []Option[int]{
			{Key: "One", Value: 1},
			{Key: "Two", Value: 2},
		}

		var result int
		selectPrompt := NewSelect(&result).
			Options(options...)
		require.Len(t, selectPrompt.options, len(options), "Options() should set the correct number of options")
	})

	t.Run("struct type", func(t *testing.T) {
		type CustomStruct struct {
			ID   int
			Name string
		}

		options := []Option[CustomStruct]{
			{Key: "First Item", Value: CustomStruct{ID: 1, Name: "first"}},
			{Key: "Second Item", Value: CustomStruct{ID: 2, Name: "second"}},
		}

		var result CustomStruct
		selectPrompt := NewSelect(&result).
			Options(options...)
		require.Len(t, selectPrompt.options, len(options), "Options() should set the correct number of options")
	})

	t.Run("should not change options when empty slice is provided", func(t *testing.T) {
		var result string
		selectPrompt := NewSelect(&result).
			Options()
		require.Len(t, selectPrompt.options, 0, "Options() should not change options when empty slice is provided")
	})
}

// Tests for [Select.SelectFunc] function.
func Test_Select_SelectFunc(t *testing.T) {
	t.Run("should set the select function", func(t *testing.T) {
		var result string
		selectFunc := func(input string) string {
			return "✓ " + input
		}
		selectPrompt := NewSelect(&result).
			SelectFunc(selectFunc)
		require.NotNil(t, selectPrompt.selectFn, "SelectFunc() should set the select function")
		testSelect := "Test Option"
		assert.Equal(t, selectFunc(testSelect), selectPrompt.selectFn(testSelect), "SelectFunc() should return the same result as the provided function")
	})
}

// Tests for [Select.Title] function.
func Test_Select_Title(t *testing.T) {
	t.Run("should set the title value", func(t *testing.T) {
		expectedTitle := "Test Title"
		var result string
		selectPrompt := NewSelect(&result).
			Title(expectedTitle)
		assert.Equal(t, expectedTitle, selectPrompt.title.Get(), "Title() should set the title value correctly")
	})

	t.Run("should return an empty title when not set", func(t *testing.T) {
		var result string
		selectPrompt := NewSelect(&result)
		assert.Equal(t, "", selectPrompt.title.Get(), "Title() should return an empty title when not set")
	})
}

// Tests for [Select.TitleFunc] function.
func Test_Select_TitleFunc(t *testing.T) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
	}

	var result string
	selectPrompt := NewSelect(&result).
		Options(options...).
		TitleFunc(func(input string) string {
			return "Choose: " + input
		})

	if selectPrompt.title.fn == nil {
		t.Error("TitleFunc should set the title function")
	}
}

// Tests for [Select.Value] function.
func Test_Select_Value(t *testing.T) {
	t.Run("should set the value pointer", func(t *testing.T) {
		var result string
		result2 := "test"
		selectPrompt := NewSelect(&result).
			Value(&result2)
		require.NotNil(t, selectPrompt.value, "Value() should set the value pointer")
		assert.NotSame(t, &result, selectPrompt.value, "Value() should set the value pointer correctly")
		assert.Same(t, &result2, selectPrompt.value, "Value() should set the value pointer to the new address")
	})

	t.Run("should allow setting a nil value pointer", func(t *testing.T) {
		result := "test"
		selectPrompt := NewSelect[string](nil).
			Value(&result)
		require.NotNil(t, selectPrompt.value, "Value() should allow setting a nil value pointer")
		assert.Same(t, &result, selectPrompt.value, "Value() should set the value pointer to the new address even when initially nil")
	})
}

// Tests for [Select.getSelectFunc] function.
func Test_Select_getSelectFunc(t *testing.T) {
	var result string
	selectPrompt := NewSelect(&result)
	t.Run("should return the string value when no select function or default section function is set", func(t *testing.T) {
		originalDefaultSelectFunc := defaultFuncs.selectFn
		t.Cleanup(func() { defaultFuncs.selectFn = originalDefaultSelectFunc })
		defaultFuncs.selectFn = nil

		expectedOutput := "Test Option"
		assert.Equal(t, expectedOutput, selectPrompt.callSelectFunc(expectedOutput), "getSelectFunc() should return the input string when no functions are set")
	})

	t.Run("should return the string transformed by default selection function if set and no prompt-specific function is set", func(t *testing.T) {
		expectedOutput := "Test Option"
		assert.Equal(t, expectedOutput, selectPrompt.callSelectFunc(expectedOutput), "getSelectFunc() should return the string transformed by the default selection function when no prompt-specific function is set")
	})

	t.Run("should return the string transformed by custom default selection function if set and no prompt-specific function is set", func(t *testing.T) {
		originalDefaultSelectFunc := defaultFuncs.selectFn
		t.Cleanup(func() { defaultFuncs.selectFn = originalDefaultSelectFunc })
		defaultFuncs.selectFn = func(input string) string {
			return "✓ " + input
		}

		expectedOutput := "✓ Test Option"
		assert.Equal(t, expectedOutput, selectPrompt.callSelectFunc("Test Option"), "getSelectFunc() should return the string transformed by the custom default selection function when no prompt-specific function is set")
	})

	t.Run("should return the string transformed by the prompt-specific selection function if set", func(t *testing.T) {
		selectFunc := func(input string) string {
			return "✓ " + input
		}
		selectPrompt.SelectFunc(selectFunc)

		expectedOutput := "✓ Test Option"
		assert.Equal(t, expectedOutput, selectPrompt.callSelectFunc("Test Option"), "getSelectFunc() should return the string transformed by the prompt-specific selection function when it is set")
	})
}

// TODO: Tests for [Select.renderOptions] function.
func Test_Select_renderOptions(t *testing.T) {
	t.Skip("Need to implement")
}
