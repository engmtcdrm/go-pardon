package keys

import (
	"testing"
)

func TestKeyConstants(t *testing.T) {
	tests := []struct {
		name     string
		key      byte
		expected byte
	}{
		{"Ctrl+C", CtrlC, 3},
		{"Delete", Backspace, 8},
		{"Carriage Return", NewLine, 10},
		{"Enter", Enter, 13},
		{"Escape", Escape, 27},
		{"Up Arrow", Up, 65},
		{"Down Arrow", Down, 66},
		{"Right Arrow", Right, 67},
		{"Left Arrow", Left, 68},
		{"No Upper", UpperN, 78},
		{"Yes Upper", UpperY, 89},
		{"Left Bracket", LeftBracket, 91},
		{"No", LowerN, 110},
		{"Yes", LowerY, 121},
		{"Backspace", Delete, 127},
		{"Capital O", UpperO, byte('O')},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != tt.expected {
				t.Errorf("%s key constant = %d; want %d", tt.name, tt.key, tt.expected)
			}
		})
	}
}

func TestArrowKeys(t *testing.T) {
	// Test arrow key constants
	arrowKeys := map[string]byte{
		"Up":    Up,
		"Down":  Down,
		"Right": Right,
		"Left":  Left,
	}

	expectedValues := map[string]byte{
		"Up":    65,
		"Down":  66,
		"Right": 67,
		"Left":  68,
	}

	for name, key := range arrowKeys {
		t.Run(name+" arrow key", func(t *testing.T) {
			expected := expectedValues[name]
			if key != expected {
				t.Errorf("Key%s = %d; want %d", name, key, expected)
			}
		})
	}
}

func TestConfirmationKeys(t *testing.T) {
	// Test yes/no keys
	confirmKeys := []struct {
		name     string
		key      byte
		expected byte
	}{
		{"Yes lowercase", LowerY, 121},
		{"Yes uppercase", UpperY, 89},
		{"No lowercase", LowerN, 110},
		{"No uppercase", UpperN, 78},
	}

	for _, tt := range confirmKeys {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != tt.expected {
				t.Errorf("%s key constant = %d; want %d", tt.name, tt.key, tt.expected)
			}
		})
	}
}

func TestSpecialKeys(t *testing.T) {
	specialKeys := []struct {
		name     string
		key      byte
		expected byte
	}{
		{"Ctrl+C", CtrlC, 3},
		{"Delete", Backspace, 8},
		{"Enter", Enter, 13},
		{"Carriage Return", NewLine, 10},
		{"Escape", Escape, 27},
		{"Backspace", Delete, 127},
		{"Left Bracket", LeftBracket, 91},
	}

	for _, tt := range specialKeys {
		t.Run(tt.name, func(t *testing.T) {
			if tt.key != tt.expected {
				t.Errorf("%s key constant = %d; want %d", tt.name, tt.key, tt.expected)
			}
		})
	}
}

// Test that key constants are reasonable
func TestKeyRanges(t *testing.T) {
	// Control characters should be in range 1-31
	if CtrlC < 1 || CtrlC > 31 {
		t.Error("KeyCtrlC should be in control character range (1-31)")
	}

	// Arrow keys should be in the expected ANSI range
	arrowKeys := []byte{Up, Down, Right, Left}
	for i, key := range arrowKeys {
		if key < 65 || key > 68 {
			t.Errorf("Arrow key %d should be in range 65-68, got %d", i, key)
		}
	}

	// Yes/No keys should be printable ASCII
	yesNoKeys := []byte{LowerY, UpperY, LowerN, UpperN}
	for _, key := range yesNoKeys {
		if key < 32 || key > 126 {
			t.Errorf("Yes/No key %d should be in printable ASCII range (32-126)", key)
		}
	}
}
