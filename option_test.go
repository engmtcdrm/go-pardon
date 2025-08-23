package pardon

import (
	"testing"
)

func TestOption_NewOption(t *testing.T) {
	t.Run("creates option with string value", func(t *testing.T) {
		option := NewOption("Test Key", "test_value")

		if option.Key != "Test Key" {
			t.Errorf("Expected key 'Test Key', got '%s'", option.Key)
		}
		if option.Value != "test_value" {
			t.Errorf("Expected value 'test_value', got '%s'", option.Value)
		}
	})

	t.Run("creates option with int value", func(t *testing.T) {
		option := NewOption("Number", 42)

		if option.Key != "Number" {
			t.Errorf("Expected key 'Number', got '%s'", option.Key)
		}
		if option.Value != 42 {
			t.Errorf("Expected value 42, got %d", option.Value)
		}
	})

	t.Run("creates option with bool value", func(t *testing.T) {
		option := NewOption("Enabled", true)

		if option.Key != "Enabled" {
			t.Errorf("Expected key 'Enabled', got '%s'", option.Key)
		}
		if option.Value != true {
			t.Errorf("Expected value true, got %v", option.Value)
		}
	})

	t.Run("creates option with empty key", func(t *testing.T) {
		option := NewOption("", "value")

		if option.Key != "" {
			t.Errorf("Expected empty key, got '%s'", option.Key)
		}
		if option.Value != "value" {
			t.Errorf("Expected value 'value', got '%s'", option.Value)
		}
	})

	t.Run("creates option with zero value", func(t *testing.T) {
		option := NewOption("Zero", 0)

		if option.Key != "Zero" {
			t.Errorf("Expected key 'Zero', got '%s'", option.Key)
		}
		if option.Value != 0 {
			t.Errorf("Expected value 0, got %d", option.Value)
		}
	})
}

func TestOption_DirectConstruction(t *testing.T) {
	t.Run("creates option via direct struct initialization", func(t *testing.T) {
		option := Option[string]{
			Key:   "Direct Key",
			Value: "direct_value",
		}

		if option.Key != "Direct Key" {
			t.Errorf("Expected key 'Direct Key', got '%s'", option.Key)
		}
		if option.Value != "direct_value" {
			t.Errorf("Expected value 'direct_value', got '%s'", option.Value)
		}
	})

	t.Run("creates option with custom struct type", func(t *testing.T) {
		type CustomType struct {
			Name string
			ID   int
		}

		customValue := CustomType{Name: "Test", ID: 123}
		option := Option[CustomType]{
			Key:   "Custom",
			Value: customValue,
		}

		if option.Key != "Custom" {
			t.Errorf("Expected key 'Custom', got '%s'", option.Key)
		}
		if option.Value.Name != "Test" || option.Value.ID != 123 {
			t.Errorf("Expected custom value {Name: Test, ID: 123}, got %+v", option.Value)
		}
	})
}

func TestOption_Types(t *testing.T) {
	t.Run("option with string value", func(t *testing.T) {
		option := NewOption("Text", "hello")

		if option.Key != "Text" {
			t.Errorf("Expected key 'Text', got '%s'", option.Key)
		}
		if option.Value != "hello" {
			t.Errorf("Expected value 'hello', got '%s'", option.Value)
		}
	})

	t.Run("option with pointer value", func(t *testing.T) {
		str := "pointer_value"
		option := NewOption("Pointer", &str)

		if option.Key != "Pointer" {
			t.Errorf("Expected key 'Pointer', got '%s'", option.Key)
		}
		if *option.Value != "pointer_value" {
			t.Errorf("Expected value 'pointer_value', got '%s'", *option.Value)
		}
	})
}
