package pardon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for [SetDefaultAnswerFunc] function.
func Test_SetDefaultAnswerFunc(t *testing.T) {
	t.Run("Set to nil", func(t *testing.T) {
		originalFn := defaultFuncs.answerFn
		t.Cleanup(func() {
			SetDefaultAnswerFunc(originalFn)
		})

		require.NotNil(t, defaultFuncs.answerFn)
		SetDefaultAnswerFunc(nil)
		require.Nil(t, defaultFuncs.answerFn)
	})

	t.Run("Set to a custom function", func(t *testing.T) {
		originalFn := defaultFuncs.answerFn
		t.Cleanup(func() {
			SetDefaultAnswerFunc(originalFn)
		})

		expectedBefore := "test"
		require.NotNil(t, defaultFuncs.answerFn)
		require.Equal(t, expectedBefore, defaultFuncs.answerFn("test"))

		expected := "custom: test"
		customFn := func(s string) string { return "custom: " + s }
		SetDefaultAnswerFunc(customFn)
		require.Equal(t, expected, defaultFuncs.answerFn("test"))
	})
}

// Tests for [SetDefaultCursorFunc] function.
func Test_SetDefaultCursorFunc(t *testing.T) {
	t.Run("Set to nil", func(t *testing.T) {
		originalFn := defaultFuncs.cursorFn
		t.Cleanup(func() {
			SetDefaultCursorFunc(originalFn)
		})

		require.NotNil(t, defaultFuncs.cursorFn)
		SetDefaultCursorFunc(nil)
		require.Nil(t, defaultFuncs.cursorFn)
	})

	t.Run("Set to a custom function", func(t *testing.T) {
		originalFn := defaultFuncs.cursorFn
		t.Cleanup(func() {
			SetDefaultCursorFunc(originalFn)
		})

		expectedBefore := "test"
		require.NotNil(t, defaultFuncs.cursorFn)
		require.Equal(t, expectedBefore, defaultFuncs.cursorFn("test"))

		expected := "custom: test"
		customFn := func(s string) string { return "custom: " + s }
		SetDefaultCursorFunc(customFn)
		require.Equal(t, expected, defaultFuncs.cursorFn("test"))
	})
}

// Tests for [SetDefaultIconFunc] function.
func Test_SetDefaultIconFunc(t *testing.T) {
	t.Run("Set to nil", func(t *testing.T) {
		originalFn := defaultFuncs.iconFn
		t.Cleanup(func() {
			SetDefaultIconFunc(originalFn)
		})

		require.NotNil(t, defaultFuncs.iconFn)
		SetDefaultIconFunc(nil)
		require.Nil(t, defaultFuncs.iconFn)
	})

	t.Run("Set to a custom function", func(t *testing.T) {
		originalFn := defaultFuncs.iconFn
		t.Cleanup(func() {
			SetDefaultIconFunc(originalFn)
		})

		expectedBefore := "test"
		require.NotNil(t, defaultFuncs.iconFn)
		require.Equal(t, expectedBefore, defaultFuncs.iconFn("test"))

		expected := "custom: test"
		customFn := func(s string) string { return "custom: " + s }
		SetDefaultIconFunc(customFn)
		require.Equal(t, expected, defaultFuncs.iconFn("test"))
	})
}

// Tests for [SetDefaultSelectFunc] function.
func Test_SetDefaultSelectFunc(t *testing.T) {
	t.Run("Set to nil", func(t *testing.T) {
		originalFn := defaultFuncs.selectFn
		t.Cleanup(func() {
			SetDefaultSelectFunc(originalFn)
		})

		require.NotNil(t, defaultFuncs.selectFn)
		SetDefaultSelectFunc(nil)
		require.Nil(t, defaultFuncs.selectFn)
	})

	t.Run("Set to a custom function", func(t *testing.T) {
		originalFn := defaultFuncs.selectFn
		t.Cleanup(func() {
			SetDefaultSelectFunc(originalFn)
		})

		expectedBefore := "test"
		require.NotNil(t, defaultFuncs.selectFn)
		require.Equal(t, expectedBefore, defaultFuncs.selectFn("test"))

		expected := "custom: test"
		customFn := func(s string) string { return "custom: " + s }
		SetDefaultSelectFunc(customFn)
		require.Equal(t, expected, defaultFuncs.selectFn("test"))
	})
}

// Tests for [SetDefaultTitleFunc] function.
func Test_SetDefaultTitleFunc(t *testing.T) {
	t.Run("Set to nil", func(t *testing.T) {
		originalFn := defaultFuncs.titleFn
		t.Cleanup(func() {
			SetDefaultTitleFunc(originalFn)
		})

		require.NotNil(t, defaultFuncs.titleFn)
		SetDefaultTitleFunc(nil)
		require.Nil(t, defaultFuncs.titleFn)
	})

	t.Run("Set to a custom function", func(t *testing.T) {
		originalFn := defaultFuncs.titleFn
		t.Cleanup(func() {
			SetDefaultTitleFunc(originalFn)
		})

		expectedBefore := "test"
		require.NotNil(t, defaultFuncs.titleFn)
		require.Equal(t, expectedBefore, defaultFuncs.titleFn("test"))

		expected := "custom: test"
		customFn := func(s string) string { return "custom: " + s }
		SetDefaultTitleFunc(customFn)
		require.Equal(t, expected, defaultFuncs.titleFn("test"))
	})
}
