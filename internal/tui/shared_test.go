package tui

import (
	"testing"
)

func TestGetTerminalHeight(t *testing.T) {
	height := GetTerminalHeight()

	// Should return a reasonable default or actual terminal height
	if height < 10 || height > 200 {
		t.Errorf("GetTerminalHeight() = %d; expected a reasonable value between 10 and 200", height)
	}
}

// Helper function to check if a string contains another string
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
