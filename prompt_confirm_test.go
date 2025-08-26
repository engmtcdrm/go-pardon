package pardon

import (
	"testing"
)

func TestConfirmCreation(t *testing.T) {
	var result bool
	confirm := NewConfirm(&result)

	if confirm == nil {
		t.Error("NewConfirm returned nil")
	}

	if confirm.icon.val == "" {
		t.Error("Confirm should have default icon")
	}
}

func TestConfirmWithTitle(t *testing.T) {
	var result bool
	confirm := NewConfirm(&result).
		Title("Are you sure?")

	if confirm.title.val != "Are you sure?" {
		t.Errorf("Title() = %q; want %q", confirm.title.val, "Are you sure?")
	}
}

func TestConfirmWithValue(t *testing.T) {
	var result bool
	confirm := NewConfirm(&result)

	if confirm.value != &result {
		t.Error("Confirm value pointer not properly set")
	}
}

func TestConfirmValidation(t *testing.T) {
	t.Run("no title", func(t *testing.T) {
		var result bool
		prompt := NewConfirm(&result)

		// We can't easily test Ask() without user interaction,
		// but we can verify the validation conditions
		if prompt.value == nil {
			t.Error("Value should be set")
		}
	})

	t.Run("no value", func(t *testing.T) {
		var result bool
		prompt := NewConfirm(&result).
			Title("Proceed?")

		// Verify title is set correctly
		if prompt.title.val != "Proceed?" {
			t.Error("Title should be set correctly")
		}
	})
}

func TestConfirmWithAnswerFunc(t *testing.T) {
	var result bool
	confirm := NewConfirm(&result).
		Title("Continue?").
		AnswerFunc(func(answer string) string {
			return "[" + answer + "]"
		})

	// We can't test the actual interaction, but we can verify configuration
	if confirm == nil {
		t.Error("Confirm with answer function returned nil")
	}
}

func TestConfirmWithIcon(t *testing.T) {
	var result bool
	confirm := NewConfirm(&result).
		Icon("❓ ")

	if confirm.icon.val != "❓ " {
		t.Errorf("Icon() = %q; want %q", confirm.icon.val, "❓ ")
	}
}

func TestConfirmWithIconFunc(t *testing.T) {
	var result bool
	confirm := NewConfirm(&result).
		IconFunc(func(input string) string {
			return "⚡ "
		})

	if confirm.icon.fn == nil {
		t.Error("IconFunc should set the icon function")
	}
}

func TestConfirmWithTitleFunc(t *testing.T) {
	var result bool
	confirm := NewConfirm(&result).
		TitleFunc(func(input string) string {
			return "Confirm: " + input
		})

	if confirm.title.fn == nil {
		t.Error("TitleFunc should set the title function")
	}
}
