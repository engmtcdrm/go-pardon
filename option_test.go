package pardon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for [NewOption] function.
func Test_NewOption(t *testing.T) {
	t.Run("creates option with string value", func(t *testing.T) {
		option := NewOption("Test Key", "test_value")
		assert.Equal(t, "Test Key", option.Key)
		assert.Equal(t, "test_value", option.Value)
	})

	t.Run("creates option with int value", func(t *testing.T) {
		option := NewOption("Number", 42)
		assert.Equal(t, "Number", option.Key)
		assert.Equal(t, 42, option.Value)
	})

	t.Run("creates option with bool value", func(t *testing.T) {
		option := NewOption("Enabled", true)
		assert.Equal(t, "Enabled", option.Key)
		assert.Equal(t, true, option.Value)
	})

	t.Run("creates option with empty key", func(t *testing.T) {
		option := NewOption("", "value")
		assert.Equal(t, "", option.Key)
		assert.Equal(t, "value", option.Value)
	})

	t.Run("creates option with zero value", func(t *testing.T) {
		option := NewOption("Zero", 0)
		assert.Equal(t, "Zero", option.Key)
		assert.Equal(t, 0, option.Value)
	})

	t.Run("creates option via direct struct initialization", func(t *testing.T) {
		option := Option[string]{
			Key:   "Direct Key",
			Value: "direct_value",
		}

		assert.Equal(t, "Direct Key", option.Key)
		assert.Equal(t, "direct_value", option.Value)
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
		assert.Equal(t, "Custom", option.Key)
		assert.Equal(t, customValue, option.Value)
	})

	t.Run("option with string value", func(t *testing.T) {
		option := NewOption("Text", "hello")

		assert.Equal(t, "Text", option.Key)
		assert.Equal(t, "hello", option.Value)
	})

	t.Run("option with pointer value", func(t *testing.T) {
		str := "pointer_value"
		option := NewOption("Pointer", &str)

		assert.Equal(t, "Pointer", option.Key)
		assert.Equal(t, &str, option.Value)
	})
}
