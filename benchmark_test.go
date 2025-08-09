package pardon

import (
	"fmt"
	"testing"
)

// Benchmark tests for performance measurement

func BenchmarkSelectCreation(b *testing.B) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
		{Key: "Option 2", Value: "value2"},
		{Key: "Option 3", Value: "value3"},
		{Key: "Option 4", Value: "value4"},
		{Key: "Option 5", Value: "value5"},
	}

	var result string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewSelect[string]().Options(options...).Value(&result).Title("Test")
	}
}

func BenchmarkSelectCreationLargeOptions(b *testing.B) {
	options := make([]Option[string], 100)
	for i := 0; i < 100; i++ {
		options[i] = Option[string]{
			Key:   fmt.Sprintf("Option %d", i+1),
			Value: fmt.Sprintf("value%d", i+1),
		}
	}

	var result string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewSelect[string]().Options(options...).Value(&result).Title("Test")
	}
}

func BenchmarkQuestionCreation(b *testing.B) {
	var result string

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewQuestion().Value(&result).Title("Test Question")
	}
}

func BenchmarkPasswordCreation(b *testing.B) {
	var result []byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewPassword().Value(&result).Title("Enter Password")
	}
}

func BenchmarkConfirmCreation(b *testing.B) {
	var result bool

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewConfirm().Value(&result).Title("Are you sure?")
	}
}

func BenchmarkNewOption(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewOption("Test Key", "test_value")
	}
}

// Benchmark memory allocations
func BenchmarkSelectAllocation(b *testing.B) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
		{Key: "Option 2", Value: "value2"},
		{Key: "Option 3", Value: "value3"},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result string
		selectPrompt := NewSelect[string]().Options(options...).Value(&result).Title("Test")
		_ = selectPrompt
	}
}

func BenchmarkFluentAPI(b *testing.B) {
	options := []Option[string]{
		{Key: "Option 1", Value: "value1"},
		{Key: "Option 2", Value: "value2"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result string
		_ = NewSelect[string]().
			Title("Select an option").
			Options(options...).
			Value(&result).
			Cursor("> ")
	}
}
