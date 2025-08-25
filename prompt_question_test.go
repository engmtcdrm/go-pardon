package pardon

import (
	"testing"
)

func TestQuestionCreation(t *testing.T) {
	var result string
	question := NewQuestion(&result)

	if question == nil {
		t.Error("NewQuestion returned nil")
	}

	if question.icon.val == "" {
		t.Error("Question should have default icon")
	}
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
		Validate(func(input string) error {
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
