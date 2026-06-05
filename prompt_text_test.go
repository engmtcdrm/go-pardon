package pardon

import (
	"io"
	"testing"

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
}

// Tests for [Text.AnswerFunc] method.
func Test_Text_AnswerFunc(t *testing.T) {
	t.Run("using default answer function", func(t *testing.T) {
		var result string
		confirm := NewQuestion(&result).
			Title("Continue?")
		confirm.terminal.Out = io.Discard
		assert.Nil(t, confirm.answerFn, "Default answer function should be nil")
	})

	t.Run("using custom answer function", func(t *testing.T) {
		var result string
		customFn := func(s string) string {
			return "Custom: " + s
		}
		confirm := NewQuestion(&result).
			Title("Continue?").
			AnswerFunc(customFn)
		confirm.terminal.Out = io.Discard
		assert.Equal(t, customFn("Test"), confirm.answerFn("Test"), "Custom answer function did not return expected result")
	})
}

// TODO: Tests for [Text.Ask] function.
func Test_Text_Ask(t *testing.T) {
	t.Skip("Need to implement")
}

// Tests for [Text.Hide] function.
func Test_Text_Hide(t *testing.T) {
	t.Run("should set hide to true", func(t *testing.T) {
		var result string
		text := NewQuestion(&result).
			Title("Enter password:").
			Hide(true)
		assert.True(t, text.terminal.Hide, "Hide should be set to true")
	})
}

func TestQuestionWithTitle(t *testing.T) {
	var result string
	question := NewQuestion(&result).
		Title("What is your name?")

	if question.title.val != "What is your name?" {
		t.Errorf("Title() = %q; want %q", question.title.val, "What is your name?")
	}
}

func TestQuestionWithValue(t *testing.T) {
	var result string
	question := NewQuestion(&result)

	if question.value != &result {
		t.Error("Question value pointer not properly set")
	}
}

func TestQuestionValidation(t *testing.T) {
	t.Run("no title", func(t *testing.T) {
		var result string
		prompt := NewQuestion(&result)

		// Test that title is empty, which should cause validation to fail
		if prompt.title.val != "" {
			t.Error("Expected title to be empty")
		}

		// We can't easily test Ask() without user interaction,
		// but we can verify the validation conditions
		if prompt.value == nil {
			t.Error("Value should be set")
		}
	})

	t.Run("no value", func(t *testing.T) {
		var result string
		prompt := NewQuestion(&result).
			Title("Test question")

		// Verify title is set correctly
		if prompt.title.val != "Test question" {
			t.Error("Title should be set correctly")
		}
	})
}

func TestQuestionWithValidate(t *testing.T) {
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

func TestQuestionWithIcon(t *testing.T) {
	var result string
	question := NewQuestion(&result).
		Icon("➤ ")

	if question.icon.val != "➤ " {
		t.Errorf("Icon() = %q; want %q", question.icon.val, "➤ ")
	}
}

func TestQuestionWithIconFunc(t *testing.T) {
	var result string
	question := NewQuestion(&result).
		IconFunc(func(input string) string {
			return "🔍 "
		})

	if question.icon.fn == nil {
		t.Error("IconFunc should set the icon function")
	}
}

func TestQuestionWithTitleFunc(t *testing.T) {
	var result string
	question := NewQuestion(&result).
		TitleFunc(func(input string) string {
			return "Dynamic: " + input
		})

	if question.title.fn == nil {
		t.Error("TitleFunc should set the title function")
	}
}
