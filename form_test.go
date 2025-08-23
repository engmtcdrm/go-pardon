package pardon

import (
	"errors"
	"testing"
)

// MockPrompt is a test double for the Prompt interface
type MockPrompt struct {
	shouldError bool
	errorMsg    string
	askCalled   bool
}

func (m *MockPrompt) Ask() error {
	m.askCalled = true
	if m.shouldError {
		return errors.New(m.errorMsg)
	}
	return nil
}

func TestNewForm(t *testing.T) {
	t.Run("creates form with no prompts", func(t *testing.T) {
		form := NewForm()
		if form == nil {
			t.Fatal("NewForm() returned nil")
		}
		if len(form.prompts) != 0 {
			t.Errorf("Expected 0 prompts, got %d", len(form.prompts))
		}
	})

	t.Run("creates form with single prompt", func(t *testing.T) {
		mockPrompt := &MockPrompt{}
		form := NewForm(mockPrompt)

		if form == nil {
			t.Fatal("NewForm() returned nil")
		}
		if len(form.prompts) != 1 {
			t.Errorf("Expected 1 prompt, got %d", len(form.prompts))
		}
	})

	t.Run("creates form with multiple prompts", func(t *testing.T) {
		mockPrompt1 := &MockPrompt{}
		mockPrompt2 := &MockPrompt{}
		mockPrompt3 := &MockPrompt{}

		form := NewForm(mockPrompt1, mockPrompt2, mockPrompt3)

		if form == nil {
			t.Fatal("NewForm() returned nil")
		}
		if len(form.prompts) != 3 {
			t.Errorf("Expected 3 prompts, got %d", len(form.prompts))
		}
	})
}

func TestForm_Ask(t *testing.T) {
	t.Run("executes all prompts successfully", func(t *testing.T) {
		mockPrompt1 := &MockPrompt{}
		mockPrompt2 := &MockPrompt{}
		mockPrompt3 := &MockPrompt{}

		form := NewForm(mockPrompt1, mockPrompt2, mockPrompt3)

		err := form.Ask()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify all prompts were called
		if !mockPrompt1.askCalled {
			t.Error("First prompt was not called")
		}
		if !mockPrompt2.askCalled {
			t.Error("Second prompt was not called")
		}
		if !mockPrompt3.askCalled {
			t.Error("Third prompt was not called")
		}
	})

	t.Run("stops on first error", func(t *testing.T) {
		mockPrompt1 := &MockPrompt{}
		mockPrompt2 := &MockPrompt{shouldError: true, errorMsg: "test error"}
		mockPrompt3 := &MockPrompt{}

		form := NewForm(mockPrompt1, mockPrompt2, mockPrompt3)

		err := form.Ask()

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "test error" {
			t.Errorf("Expected 'test error', got %v", err)
		}

		// Verify execution stopped after the error
		if !mockPrompt1.askCalled {
			t.Error("First prompt was not called")
		}
		if !mockPrompt2.askCalled {
			t.Error("Second prompt was not called")
		}
		if mockPrompt3.askCalled {
			t.Error("Third prompt should not have been called after error")
		}
	})

	t.Run("handles error on first prompt", func(t *testing.T) {
		mockPrompt1 := &MockPrompt{shouldError: true, errorMsg: "first prompt error"}
		mockPrompt2 := &MockPrompt{}

		form := NewForm(mockPrompt1, mockPrompt2)

		err := form.Ask()

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "first prompt error" {
			t.Errorf("Expected 'first prompt error', got %v", err)
		}

		// Verify only first prompt was called
		if !mockPrompt1.askCalled {
			t.Error("First prompt was not called")
		}
		if mockPrompt2.askCalled {
			t.Error("Second prompt should not have been called after error")
		}
	})

	t.Run("handles empty form", func(t *testing.T) {
		form := NewForm()

		err := form.Ask()

		if err != nil {
			t.Errorf("Expected no error for empty form, got %v", err)
		}
	})

	t.Run("handles single successful prompt", func(t *testing.T) {
		mockPrompt := &MockPrompt{}
		form := NewForm(mockPrompt)

		err := form.Ask()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !mockPrompt.askCalled {
			t.Error("Prompt was not called")
		}
	})

	t.Run("handles single failing prompt", func(t *testing.T) {
		mockPrompt := &MockPrompt{shouldError: true, errorMsg: "single prompt error"}
		form := NewForm(mockPrompt)

		err := form.Ask()

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "single prompt error" {
			t.Errorf("Expected 'single prompt error', got %v", err)
		}
		if !mockPrompt.askCalled {
			t.Error("Prompt was not called")
		}
	})
}

func TestForm_Integration(t *testing.T) {
	t.Run("complex scenario with mixed success and failure", func(t *testing.T) {
		// Create a scenario with multiple prompts where the 3rd one fails
		prompts := []*MockPrompt{
			{}, // Success
			{}, // Success
			{shouldError: true, errorMsg: "validation failed"}, // Failure
			{}, // Should not be reached
			{}, // Should not be reached
		}

		form := NewForm(
			prompts[0],
			prompts[1],
			prompts[2],
			prompts[3],
			prompts[4],
		)

		err := form.Ask()

		// Should return error from 3rd prompt
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "validation failed" {
			t.Errorf("Expected 'validation failed', got %v", err)
		}

		// Check execution sequence
		if !prompts[0].askCalled {
			t.Error("First prompt should have been called")
		}
		if !prompts[1].askCalled {
			t.Error("Second prompt should have been called")
		}
		if !prompts[2].askCalled {
			t.Error("Third prompt should have been called")
		}
		if prompts[3].askCalled {
			t.Error("Fourth prompt should not have been called")
		}
		if prompts[4].askCalled {
			t.Error("Fifth prompt should not have been called")
		}
	})
}
