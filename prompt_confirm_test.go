package pardon

import (
	"testing"
)

func TestConfirmCreation(t *testing.T) {
	confirm := NewConfirm()

	if confirm == nil {
		t.Error("NewConfirm returned nil")
	}

	// Test initial state
	if confirm.value != nil {
		t.Error("Confirm value should be nil initially")
	}

	if confirm.icon.val == "" {
		t.Error("Confirm should have default icon")
	}
}

func TestConfirmWithTitle(t *testing.T) {
	var result bool
	confirm := NewConfirm().Value(&result).Title("Are you sure?")

	if confirm.title.val != "Are you sure?" {
		t.Errorf("Title() = %q; want %q", confirm.title.val, "Are you sure?")
	}
}

func TestConfirmWithValue(t *testing.T) {
	var result bool
	confirm := NewConfirm().Value(&result)

	if confirm.value != &result {
		t.Error("Confirm value pointer not properly set")
	}
}

func TestConfirmValidation(t *testing.T) {
	t.Run("no title", func(t *testing.T) {
		var result bool
		prompt := NewConfirm().Value(&result)

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
		prompt := NewConfirm().Title("Proceed?")

		// Test that value is nil, which should cause validation to fail
		if prompt.value != nil {
			t.Error("Expected value to be nil")
		}

		// Verify title is set correctly
		if prompt.title.val != "Proceed?" {
			t.Error("Title should be set correctly")
		}
	})
}

func TestConfirmWithAnswerFunc(t *testing.T) {
	var result bool
	confirm := NewConfirm().Value(&result).Title("Continue?").AnswerFunc(func(answer string) string {
		return "[" + answer + "]"
	})

	// We can't test the actual interaction, but we can verify configuration
	if confirm == nil {
		t.Error("Confirm with answer function returned nil")
	}
}
