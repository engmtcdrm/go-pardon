package pardon

import (
	"testing"
)

func TestPasswordCreation(t *testing.T) {
	password := NewPassword()

	if password == nil {
		t.Error("NewPassword returned nil")
	}

	// Test initial state
	if password.value != nil {
		t.Error("Password value should be nil initially")
	}

	if password.icon.val == "" {
		t.Error("Password should have default icon")
	}
}

func TestPasswordWithTitle(t *testing.T) {
	var result []byte
	password := NewPassword().Value(&result).Title("Enter password:")

	if password.title.val != "Enter password:" {
		t.Errorf("Title() = %q; want %q", password.title.val, "Enter password:")
	}
}

func TestPasswordWithValue(t *testing.T) {
	var result []byte
	password := NewPassword().Value(&result)

	if password.value != &result {
		t.Error("Password value pointer not properly set")
	}
}

func TestPasswordValidation(t *testing.T) {
	t.Run("no title", func(t *testing.T) {
		var result []byte
		prompt := NewPassword().Value(&result)

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
		prompt := NewPassword().Title("Enter password:")

		// Test that value is nil, which should cause validation to fail
		if prompt.value != nil {
			t.Error("Expected value to be nil")
		}

		// Verify title is set correctly
		if prompt.title.val != "Enter password:" {
			t.Error("Title should be set correctly")
		}
	})
}

func TestPasswordWithValidate(t *testing.T) {
	var result []byte
	password := NewPassword().Value(&result).Title("Enter password:").Validate(func(input []byte) error {
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
