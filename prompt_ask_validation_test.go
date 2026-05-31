package pardon

import (
	"testing"
)

// TestAskValidation tests the Ask() method validation logic
// without requiring user interaction
func TestAskValidation(t *testing.T) {
	t.Run("question with no title should return error", func(t *testing.T) {
		var result string
		question := NewOldQuestion(&result)

		err := question.Ask()
		if err == nil {
			t.Error("Expected error when asking question with no title")
		}
	})

	t.Run("password with no title should return error", func(t *testing.T) {
		var result []byte
		password := NewOldPassword(&result)

		err := password.Ask()
		if err == nil {
			t.Error("Expected error when asking password with no title")
		}
	})

	t.Run("select with no options should return error", func(t *testing.T) {
		var result string
		selectPrompt := NewSelect(&result).
			Title("Choose")

		err := selectPrompt.Ask()
		if err != ErrNoSelectOptions {
			t.Errorf("Expected ErrNoSelectOptions but got: %v", err)
		}
	})

	t.Run("select with no title should return error", func(t *testing.T) {
		options := []Option[string]{
			{Key: "Option 1", Value: "value1"},
		}

		var result string
		selectPrompt := NewSelect(&result).
			Options(options...)

		err := selectPrompt.Ask()
		if err == nil {
			t.Error("Expected error when asking select with no title")
		}
	})
}

// TestEvalFunctionality tests the eval type behavior in prompts
func TestEvalFunctionality(t *testing.T) {
	t.Run("question with title function", func(t *testing.T) {
		var result string
		question := NewOldQuestion(&result).
			TitleFunc(func(input string) string {
				return "Dynamic title: " + input
			})

		// Test that Get() uses the function
		title := question.title.Get()
		if title != "Dynamic title: " {
			t.Errorf("Title.Get() = %q; want %q", title, "Dynamic title: ")
		}
	})

	t.Run("question with icon function", func(t *testing.T) {
		var result string
		question := NewOldQuestion(&result).
			IconFunc(func(input string) string {
				return "📝 "
			})

		// Test that Get() uses the function
		icon := question.icon.Get()
		if icon != "📝 " {
			t.Errorf("Icon.Get() = %q; want %q", icon, "📝 ")
		}
	})

	t.Run("select with cursor function", func(t *testing.T) {
		options := []Option[string]{
			{Key: "Option 1", Value: "value1"},
		}

		var result string
		selectPrompt := NewSelect(&result).
			Options(options...).
			CursorFunc(func(input string) string {
				return "▶ "
			})

		// Test that Get() uses the function
		cursor := selectPrompt.cursor.Get()
		if cursor != "▶ " {
			t.Errorf("Cursor.Get() = %q; want %q", cursor, "▶ ")
		}
	})
}

// TestMethodChaining ensures all methods return proper instances for chaining
func TestMethodChaining(t *testing.T) {
	t.Run("question method chaining", func(t *testing.T) {
		var result string
		question := NewOldQuestion(&result).
			Title("Name?").
			Icon("👤 ").
			Value(&result).
			AnswerFunc(func(s string) string { return s })

		if question == nil {
			t.Error("Method chaining should return valid instance")
		}
	})

	t.Run("password method chaining", func(t *testing.T) {
		var result []byte
		password := NewOldPassword(&result).
			Title("Password?").
			Icon("🔒 ").
			Value(&result).
			AnswerFunc(func(s string) string { return "***" })

		if password == nil {
			t.Error("Method chaining should return valid instance")
		}
	})

	t.Run("confirm method chaining", func(t *testing.T) {
		var result bool
		confirm := NewConfirm(&result).
			Title("Continue?").
			Icon("❓ ").
			AnswerFunc(func(s string) string { return s })

		if confirm == nil {
			t.Error("Method chaining should return valid instance")
		}
	})

	t.Run("select method chaining", func(t *testing.T) {
		options := []Option[string]{
			{Key: "Option 1", Value: "value1"},
		}

		var result string
		selectPrompt := NewSelect(&result).
			Options(options...).
			Title("Choose:").
			Icon("🎯 ").
			Cursor("→ ").
			AnswerFunc(func(s string) string { return s })

		if selectPrompt == nil {
			t.Error("Method chaining should return valid instance")
		}
	})
}
