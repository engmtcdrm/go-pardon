package pardon

import (
	"testing"
)

func TestSelectCreation(t *testing.T) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
		{Key: "Option 2", Value: "value2"},
	}

	var result string
	selectPrompt := NewSelect[string]().Options(options...).Value(&result)

	if selectPrompt == nil {
		t.Error("NewSelect returned nil")
	}

	// Test that the result pointer is properly set
	if selectPrompt.value != &result {
		t.Error("Select value pointer not properly set")
	}
}

func TestSelectWithTitle(t *testing.T) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
	}

	var result string
	selectPrompt := NewSelect[string]().Options(options...).Value(&result).Title("Test Title")

	if selectPrompt.title.val != "Test Title" {
		t.Errorf("Title() = %q; want %q", selectPrompt.title.val, "Test Title")
	}
}

func TestSelectWithCursor(t *testing.T) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
	}

	var result string
	selectPrompt := NewSelect[string]().Options(options...).Value(&result).Cursor("→ ")

	if selectPrompt.cursor.val != "→ " {
		t.Errorf("Cursor() = %q; want %q", selectPrompt.cursor.val, "→ ")
	}
}

func TestSelectValidation(t *testing.T) {
	t.Run("no title", func(t *testing.T) {
		options := []Option[string]{
			{Key: "Option 1", Value: "value1"},
		}

		var result string
		prompt := NewSelect[string]().Options(options...).Value(&result)

		// Test that title is empty, which should cause validation to fail
		if prompt.title.val != "" {
			t.Error("Expected title to be empty")
		}

		// We can't easily test Ask() without user interaction,
		// but we can verify the validation conditions
		if len(prompt.options) == 0 {
			t.Error("Options should be set")
		}

		if prompt.value == nil {
			t.Error("Value should be set")
		}
	})

	t.Run("no options", func(t *testing.T) {
		var result string
		prompt := NewSelect[string]().Value(&result).Title("Test")

		// Test that options are empty, which should cause validation to fail
		if len(prompt.options) != 0 {
			t.Error("Expected options to be empty")
		}

		// Verify other conditions are met
		if prompt.title.val != "Test" {
			t.Error("Title should be set correctly")
		}

		if prompt.value == nil {
			t.Error("Value should be set")
		}
	})

	t.Run("no value", func(t *testing.T) {
		options := []Option[string]{
			{Key: "Option 1", Value: "value1"},
		}

		prompt := NewSelect[string]().Options(options...).Title("Test")

		// Test that value is nil, which should cause validation to fail
		if prompt.value != nil {
			t.Error("Expected value to be nil")
		}

		// Verify other conditions are met
		if len(prompt.options) == 0 {
			t.Error("Options should be set")
		}

		if prompt.title.val != "Test" {
			t.Error("Title should be set correctly")
		}
	})
}

func TestSelectGenericTypes(t *testing.T) {
	t.Run("string type", func(t *testing.T) {
		options := []Option[string]{
			{Key: "First", Value: "first_value"},
			{Key: "Second", Value: "second_value"},
		}

		var result string
		selectPrompt := NewSelect[string]().Options(options...).Value(&result)

		if selectPrompt == nil {
			t.Error("NewSelect for string type returned nil")
		}
	})

	t.Run("int type", func(t *testing.T) {
		options := []Option[int]{
			{Key: "One", Value: 1},
			{Key: "Two", Value: 2},
		}

		var result int
		selectPrompt := NewSelect[int]().Options(options...).Value(&result)

		if selectPrompt == nil {
			t.Error("NewSelect for int type returned nil")
		}
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
		selectPrompt := NewSelect[CustomStruct]().Options(options...).Value(&result)

		if selectPrompt == nil {
			t.Error("NewSelect for struct type returned nil")
		}
	})
}

func TestSelectCursorPositioning(t *testing.T) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
		{Key: "Option 2", Value: "value2"},
		{Key: "Option 3", Value: "value3"},
	}

	var result string
	selectPrompt := NewSelect[string]().Options(options...).Value(&result)

	// Test initial cursor position
	if selectPrompt.cursorPos != 0 {
		t.Errorf("Initial cursor position = %d; want 0", selectPrompt.cursorPos)
	}
}

func TestSelectScrollOffset(t *testing.T) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
		{Key: "Option 2", Value: "value2"},
		{Key: "Option 3", Value: "value3"},
	}

	var result string
	selectPrompt := NewSelect[string]().Options(options...).Value(&result)

	// Test initial scroll offset
	if selectPrompt.scrollOffset != 0 {
		t.Errorf("Initial scroll offset = %d; want 0", selectPrompt.scrollOffset)
	}
}

func TestNewOption(t *testing.T) {
	option := NewOption("test key", "test value")

	if option.Key != "test key" {
		t.Errorf("NewOption Key = %q; want %q", option.Key, "test key")
	}

	if option.Value != "test value" {
		t.Errorf("NewOption Value = %q; want %q", option.Value, "test value")
	}
}
