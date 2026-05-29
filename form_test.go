package pardon

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

// Tests for [NewForm] function.
func Test_NewForm(t *testing.T) {
	t.Run("creates form with no prompts", func(t *testing.T) {
		form := NewForm()
		require.NotNil(t, form, "NewForm() returned nil")
		require.Equalf(t, len(form.prompts), 0, "Expected 0 prompts, got %d", len(form.prompts))
	})

	t.Run("creates form with single prompt", func(t *testing.T) {
		mockPrompt := &MockPrompt{}
		form := NewForm(mockPrompt)
		require.NotNil(t, form, "NewForm() returned nil")
		require.Equalf(t, len(form.prompts), 1, "Expected 1 prompt, got %d", len(form.prompts))
	})

	t.Run("creates form with multiple prompts", func(t *testing.T) {
		mockPrompt1 := &MockPrompt{}
		mockPrompt2 := &MockPrompt{}
		mockPrompt3 := &MockPrompt{}

		form := NewForm(mockPrompt1, mockPrompt2, mockPrompt3)
		require.NotNil(t, form, "NewForm() returned nil")
		require.Equalf(t, len(form.prompts), 3, "Expected 3 prompts, got %d", len(form.prompts))
	})
}

// Tests for [Form.Ask] function.
func Test_Form_Ask(t *testing.T) {
	t.Run("executes all prompts successfully", func(t *testing.T) {
		mockPrompt1 := &MockPrompt{}
		mockPrompt2 := &MockPrompt{}
		mockPrompt3 := &MockPrompt{}

		form := NewForm(
			mockPrompt1,
			mockPrompt2,
			mockPrompt3,
		)

		err := form.Ask()
		require.NoErrorf(t, err, "Expected no error, got %v", err)
		assert.True(t, mockPrompt1.askCalled, "First prompt was not called")
		assert.True(t, mockPrompt2.askCalled, "Second prompt was not called")
		assert.True(t, mockPrompt3.askCalled, "Third prompt was not called")
	})

	t.Run("stops on first error", func(t *testing.T) {
		mockPrompt1 := &MockPrompt{}
		mockPrompt2 := &MockPrompt{shouldError: true, errorMsg: "test error"}
		mockPrompt3 := &MockPrompt{}

		form := NewForm(
			mockPrompt1,
			mockPrompt2,
			mockPrompt3,
		)

		err := form.Ask()
		require.Error(t, err, "Expected error, got nil")
		require.Equal(t, "test error", err.Error(), "Error message did not match")

		// Verify execution stopped after the error
		assert.True(t, mockPrompt1.askCalled, "First prompt was not called")
		assert.True(t, mockPrompt2.askCalled, "Second prompt was not called")
		assert.False(t, mockPrompt3.askCalled, "Third prompt should not have been called after error")
	})

	t.Run("handles error on first prompt", func(t *testing.T) {
		mockPrompt1 := &MockPrompt{shouldError: true, errorMsg: "first prompt error"}
		mockPrompt2 := &MockPrompt{}

		form := NewForm(mockPrompt1, mockPrompt2)

		err := form.Ask()
		require.Error(t, err, "Expected error, got nil")
		require.Equal(t, "first prompt error", err.Error(), "Error message did not match")

		// Verify only first prompt was called
		assert.True(t, mockPrompt1.askCalled, "First prompt was not called")
		assert.False(t, mockPrompt2.askCalled, "Second prompt should not have been called after error")
	})

	t.Run("handles empty form", func(t *testing.T) {
		form := NewForm()

		err := form.Ask()
		require.Error(t, err, "Expected error, got nil")
		assert.Equal(t, ErrNoPrompts, err, "Expected ErrNoPrompts for empty form")
	})

	t.Run("handles single successful prompt", func(t *testing.T) {
		mockPrompt := &MockPrompt{}
		form := NewForm(mockPrompt)

		err := form.Ask()
		require.NoErrorf(t, err, "Expected no error, got %v", err)
		assert.True(t, mockPrompt.askCalled, "Prompt was not called")
	})

	t.Run("handles single failing prompt", func(t *testing.T) {
		mockPrompt := &MockPrompt{shouldError: true, errorMsg: "single prompt error"}
		form := NewForm(mockPrompt)

		err := form.Ask()
		require.Error(t, err, "Expected error, got nil")
		require.Equal(t, "single prompt error", err.Error(), "Error message did not match")
		assert.True(t, mockPrompt.askCalled, "Prompt was not called")
	})
}
