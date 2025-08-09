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
		{"Ctrl+C", KeyCtrlC, 3},
		{"Delete", KeyDelete, 8},
		{"Carriage Return", KeyCarriageReturn, 10},
		{"Enter", KeyEnter, 13},
		{"Escape", KeyEscape, 27},
		{"Up Arrow", KeyUp, 65},
		{"Down Arrow", KeyDown, 66},
		{"Right Arrow", KeyRight, 67},
		{"Left Arrow", KeyLeft, 68},
		{"No Upper", KeyNoUpper, 78},
		{"Yes Upper", KeyYesUpper, 89},
		{"Left Bracket", KeyLeftBracket, 91},
		{"No", KeyNo, 110},
		{"Yes", KeyYes, 121},
		{"Backspace", KeyBackspace, 127},
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
		"Up":    KeyUp,
		"Down":  KeyDown,
		"Right": KeyRight,
		"Left":  KeyLeft,
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
		{"Yes lowercase", KeyYes, 121},
		{"Yes uppercase", KeyYesUpper, 89},
		{"No lowercase", KeyNo, 110},
		{"No uppercase", KeyNoUpper, 78},
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
		{"Ctrl+C", KeyCtrlC, 3},
		{"Delete", KeyDelete, 8},
		{"Enter", KeyEnter, 13},
		{"Carriage Return", KeyCarriageReturn, 10},
		{"Escape", KeyEscape, 27},
		{"Backspace", KeyBackspace, 127},
		{"Left Bracket", KeyLeftBracket, 91},
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
	if KeyCtrlC < 1 || KeyCtrlC > 31 {
		t.Error("KeyCtrlC should be in control character range (1-31)")
	}

	// Arrow keys should be in the expected ANSI range
	arrowKeys := []byte{KeyUp, KeyDown, KeyRight, KeyLeft}
	for i, key := range arrowKeys {
		if key < 65 || key > 68 {
			t.Errorf("Arrow key %d should be in range 65-68, got %d", i, key)
		}
	}

	// Yes/No keys should be printable ASCII
	yesNoKeys := []byte{KeyYes, KeyYesUpper, KeyNo, KeyNoUpper}
	for _, key := range yesNoKeys {
		if key < 32 || key > 126 {
			t.Errorf("Yes/No key %d should be in printable ASCII range (32-126)", key)
		}
	}
}
