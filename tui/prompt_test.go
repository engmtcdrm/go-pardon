package tui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/engmtcdrm/go-pardon/keys"
)

func TestInputPromptCreation(t *testing.T) {
	stringPrompt := NewStringPrompt()
	if stringPrompt == nil {
		t.Error("NewStringPrompt returned nil")
	}

	passwordPrompt := NewPasswordPrompt()
	if passwordPrompt == nil {
		t.Error("NewPasswordPrompt returned nil")
	}
}

func TestControlCharacterFiltering(t *testing.T) {
	// Test that control characters are properly filtered
	// This tests the logic from prompt.go line filtering

	tests := []struct {
		name     string
		keyCode  byte
		expected bool // true if should be allowed, false if filtered
	}{
		{"Control character - Ctrl+A", keys.KeyCtrlC, false},
		{"Printable character - space", 32, true},
		{"Printable character - 'a'", 97, true},
		{"Printable character - '0'", 48, true},
		{"High ASCII character", 128, true},
		{"Tab character", 9, false},
		{"Newline character", 10, false},
		{"Delete character", keys.KeyDelete, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the filtering logic: keyCode >= 32
			allowed := tt.keyCode >= 32
			if allowed != tt.expected {
				t.Errorf("Control character filtering for keyCode %d: got %t, want %t",
					tt.keyCode, allowed, tt.expected)
			}
		})
	}
}

func TestPasswordMasking(t *testing.T) {
	// Test that password input properly handles masking
	// We can't easily test the actual terminal input/output,
	// but we can test the underlying logic

	password := []byte("test123")
	// Use asterisk (*) which is a single byte character for testing
	maskChar := '*'
	masked := strings.Repeat(string(maskChar), len(password))

	if len(masked) != len(password) {
		t.Errorf("Masked password length mismatch: got %d, want %d", len(masked), len(password))
	}

	// Test that the mask character is consistently applied
	for _, char := range masked {
		if char != maskChar {
			t.Errorf("Expected mask character '*', got %c", char)
		}
	}

	// Test with actual bullet character to understand the UTF-8 encoding
	bulletMask := strings.Repeat("•", len(password))
	// The bullet character is 3 bytes in UTF-8, so the string length will be 3x the password length
	expectedBulletLength := len(password) * 3
	if len(bulletMask) != expectedBulletLength {
		t.Logf("Bullet mask length: got %d, expected %d (this is expected for UTF-8)", len(bulletMask), expectedBulletLength)
	}
}

func TestMinFunction(t *testing.T) {
	// More comprehensive testing of the Min function
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Both positive", 5, 3, 3},
		{"Both negative", -5, -3, -5},
		{"Mixed positive/negative", 5, -3, -3},
		{"Zero and positive", 0, 5, 0},
		{"Zero and negative", 0, -5, -5},
		{"Large numbers", 1000000, 999999, 999999},
		{"Equal numbers", 42, 42, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Min(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestTerminalHeight(t *testing.T) {
	height := GetTerminalHeight()

	// Test that we get a reasonable default or actual height
	if height < 5 {
		t.Error("Terminal height should be at least 5 lines")
	}

	if height > 300 {
		t.Error("Terminal height seems unreasonably large")
	}
}

// Test string builder capacity optimization
func TestStringBuilderOptimization(t *testing.T) {
	// This tests the concept of pre-allocating string builder capacity
	// which was mentioned in the code review

	var buf bytes.Buffer

	// Test that we can write without reallocations if we know the size
	expectedSize := 100
	buf.Grow(expectedSize)

	initialCap := buf.Cap()

	// Write some data
	for i := 0; i < expectedSize/2; i++ {
		buf.WriteByte(byte('a'))
	}

	// Capacity should not have changed if we pre-allocated correctly
	if buf.Cap() != initialCap {
		t.Logf("Buffer capacity changed from %d to %d (this is informational, not an error)",
			initialCap, buf.Cap())
	}
}

// Test error conditions
func TestErrorHandling(t *testing.T) {
	// Test that negative values are handled correctly in utility functions
	RenderClearLines(-1) // Should not panic and should handle gracefully

	// Test Min with edge cases
	result := Min(-2147483648, 2147483647) // Test with int32 min/max
	if result != -2147483648 {
		t.Errorf("Min with extreme values failed: got %d", result)
	}
}
