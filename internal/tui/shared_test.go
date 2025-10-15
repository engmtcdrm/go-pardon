package tui

import (
	"testing"
)

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"a smaller than b", 3, 7, 3},
		{"b smaller than a", 10, 5, 5},
		{"equal values", 4, 4, 4},
		{"negative values", -3, -7, -7},
		{"zero values", 0, 0, 0},
		{"negative and positive", -5, 3, -5},
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

func TestGetTerminalHeight(t *testing.T) {
	height := GetTerminalHeight()

	// Should return a reasonable default or actual terminal height
	if height < 10 || height > 200 {
		t.Errorf("GetTerminalHeight() = %d; expected a reasonable value between 10 and 200", height)
	}
}

func TestRenderFormattedOutput(t *testing.T) {
	tests := []struct {
		name     string
		question string
		result   string
		contains []string // strings that should be in the output
	}{
		{
			name:     "basic formatting",
			question: "What is your name?",
			result:   "John",
			contains: []string{"What is your name?", "John"},
		},
		{
			name:     "empty result",
			question: "Enter value:",
			result:   "",
			contains: []string{"Enter value:"},
		},
		{
			name:     "special characters",
			question: "Select [y/N]:",
			result:   "y",
			contains: []string{"Select [y/N]:", "y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := RenderFormattedOutput(tt.question, tt.result)

			for _, contain := range tt.contains {
				if !containsString(output, contain) {
					t.Errorf("RenderFormattedOutput() output doesn't contain %q\nOutput: %q", contain, output)
				}
			}
		})
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
