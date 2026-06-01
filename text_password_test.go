package pardon

import (
	"testing"
)

func TestPasswordCreation(t *testing.T) {
	var result string
	password := NewPassword(&result)

	if password == nil {
		t.Error("NewPassword returned nil")
	}

	if password.icon.val == "" {
		t.Error("Password should have default icon")
	}
}

func TestPasswordWithTitle(t *testing.T) {
	var result string
	password := NewPassword(&result).
		Title("Enter password:")

	if password.title.val != "Enter password:" {
		t.Errorf("Title() = %q; want %q", password.title.val, "Enter password:")
	}
}

func TestPasswordWithValue(t *testing.T) {
	var result string
	password := NewPassword(&result)

	if password.value != &result {
		t.Error("Password value pointer not properly set")
	}
}

func TestPasswordValidation(t *testing.T) {
	t.Run("no title", func(t *testing.T) {
		var result string
		prompt := NewPassword(&result)

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
		prompt := NewPassword(&result).
			Title("Enter password:")

		// Verify title is set correctly
		if prompt.title.val != "Enter password:" {
			t.Error("Title should be set correctly")
		}
	})
}

func TestPasswordWithValidate(t *testing.T) {
	var result string
	password := NewPassword(&result).
		Title("Enter password:").
		ValidateFunc(func(input string) error {
			if len(input) < 6 {
				return ErrNoValue // Using existing error for test simplicity
			}
			return nil
		})

	// Validate functionality test - we can't easily test the actual input
	// but we can verify the password was configured properly
	if password == nil {
		t.Error("Password with validation returned nil")
	}
}

func TestPasswordWithIcon(t *testing.T) {
	var result string
	password := NewPassword(&result).
		Icon("🔐 ")

	if password.icon.val != "🔐 " {
		t.Errorf("Icon() = %q; want %q", password.icon.val, "🔐 ")
	}
}

func TestPasswordWithIconFunc(t *testing.T) {
	var result string
	password := NewPassword(&result).
		IconFunc(func(input string) string {
			return "🛡️ "
		})

	if password.icon.fn == nil {
		t.Error("IconFunc should set the icon function")
	}
}

func TestPasswordWithTitleFunc(t *testing.T) {
	var result string
	password := NewPassword(&result).
		TitleFunc(func(input string) string {
			return "Secure: " + input
		})

	if password.title.fn == nil {
		t.Error("TitleFunc should set the title function")
	}
}

func TestPasswordWithAnswerFunc(t *testing.T) {
	var result string
	password := NewPassword(&result).
		Title("Password").
		AnswerFunc(func(answer string) string {
			return "***hidden***"
		})

	if password.answerFn == nil {
		t.Error("AnswerFunc should set the answer function")
	}
}
