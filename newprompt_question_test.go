package pardon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for [NewQuestion] function.
func Test_NewQuestion(t *testing.T) {
	t.Run("should create a new question", func(t *testing.T) {
		var result string
		question := NewQuestion(&result)
		require.NotNil(t, question)
		assert.NotEmpty(t, question.icon.val)
		assert.Equal(t, &result, question.value)
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

func TestQuestionWithAnswerFunc(t *testing.T) {
	var result string
	question := NewQuestion(&result).
		Title("Test").
		AnswerFunc(func(answer string) string {
			return "[" + answer + "]"
		})

	if question.answerFn == nil {
		t.Error("AnswerFunc should set the answer function")
	}
}
