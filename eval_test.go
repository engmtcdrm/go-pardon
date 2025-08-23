package pardon

import (
	"strings"
	"testing"
)

func TestEval_Get(t *testing.T) {
	t.Run("returns static value when no functions are set", func(t *testing.T) {
		e := eval[string]{val: "hello"}

		result := e.Get()
		if result != "hello" {
			t.Errorf("Expected 'hello', got '%s'", result)
		}
	})

	t.Run("applies fn when set", func(t *testing.T) {
		e := eval[string]{
			val: "hello",
			fn:  func(s string) string { return strings.ToUpper(s) },
		}

		result := e.Get()
		if result != "HELLO" {
			t.Errorf("Expected 'HELLO', got '%s'", result)
		}
	})

	t.Run("applies defaultFn when fn is nil", func(t *testing.T) {
		e := eval[string]{
			val:       "hello",
			fn:        nil,
			defaultFn: func(s string) string { return strings.ToUpper(s) },
		}

		result := e.Get()
		if result != "HELLO" {
			t.Errorf("Expected 'HELLO', got '%s'", result)
		}
	})

	t.Run("fn takes precedence over defaultFn", func(t *testing.T) {
		e := eval[string]{
			val:       "hello",
			fn:        func(s string) string { return strings.ToUpper(s) },
			defaultFn: func(s string) string { return strings.ToLower(s) },
		}

		result := e.Get()
		if result != "HELLO" {
			t.Errorf("Expected 'HELLO' (fn should take precedence), got '%s'", result)
		}
	})

	t.Run("works with int type", func(t *testing.T) {
		e := eval[int]{
			val: 5,
			fn:  func(i int) int { return i * 2 },
		}

		result := e.Get()
		if result != 10 {
			t.Errorf("Expected 10, got %d", result)
		}
	})

	t.Run("works with bool type", func(t *testing.T) {
		e := eval[bool]{
			val: true,
			fn:  func(b bool) bool { return !b },
		}

		result := e.Get()
		if result != false {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("works with zero values", func(t *testing.T) {
		e := eval[string]{val: ""}

		result := e.Get()
		if result != "" {
			t.Errorf("Expected empty string, got '%s'", result)
		}
	})

	t.Run("works with nil defaultFn", func(t *testing.T) {
		e := eval[string]{
			val:       "test",
			fn:        nil,
			defaultFn: nil,
		}

		result := e.Get()
		if result != "test" {
			t.Errorf("Expected 'test', got '%s'", result)
		}
	})
}

func TestEval_ComplexTransformations(t *testing.T) {
	t.Run("string concatenation with fn", func(t *testing.T) {
		e := eval[string]{
			val: "world",
			fn:  func(s string) string { return "hello " + s },
		}

		result := e.Get()
		if result != "hello world" {
			t.Errorf("Expected 'hello world', got '%s'", result)
		}
	})

	t.Run("mathematical operations with defaultFn", func(t *testing.T) {
		e := eval[int]{
			val:       10,
			defaultFn: func(i int) int { return i*i + 5 },
		}

		result := e.Get()
		if result != 105 { // 10*10 + 5
			t.Errorf("Expected 105, got %d", result)
		}
	})

	t.Run("chaining behavior simulation", func(t *testing.T) {
		// Simulate what might happen in a prompt chain
		baseValue := "user input"

		// First transformation (like title formatting)
		e1 := eval[string]{
			val: baseValue,
			fn:  func(s string) string { return "Question: " + s },
		}

		intermediateResult := e1.Get()

		// Second transformation (like icon prepending)
		e2 := eval[string]{
			val: intermediateResult,
			fn:  func(s string) string { return "[?] " + s },
		}

		finalResult := e2.Get()
		expected := "[?] Question: user input"

		if finalResult != expected {
			t.Errorf("Expected '%s', got '%s'", expected, finalResult)
		}
	})
}

func TestEval_CustomTypes(t *testing.T) {
	type CustomStruct struct {
		Name  string
		Count int
	}

	t.Run("works with custom struct type", func(t *testing.T) {
		original := CustomStruct{Name: "test", Count: 1}
		e := eval[CustomStruct]{
			val: original,
			fn: func(cs CustomStruct) CustomStruct {
				cs.Count = cs.Count * 2
				return cs
			},
		}

		result := e.Get()
		if result.Name != "test" || result.Count != 2 {
			t.Errorf("Expected {Name: test, Count: 2}, got %+v", result)
		}
	})

	t.Run("works with pointer type", func(t *testing.T) {
		original := "original"
		e := eval[*string]{
			val: &original,
			fn: func(s *string) *string {
				modified := "modified " + *s
				return &modified
			},
		}

		result := e.Get()
		if *result != "modified original" {
			t.Errorf("Expected 'modified original', got '%s'", *result)
		}
	})
}
