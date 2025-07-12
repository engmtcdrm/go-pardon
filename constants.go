package gocliselect

const (
	// Key codes for navigation and actions
	KeyCtrlC          = byte(3)  // Ctrl+C generates ASCII 3 in raw mode
	KeyCarriageReturn = byte(10) // Additional for cross-platform compatibility
	KeyEnter          = byte(13)
	KeyEscape         = byte(27)
	KeyUp             = byte(65)
	KeyDown           = byte(66)
)

// NavigationKeys defines a map of specific byte keycodes related to
// navigation functionality, such as up or down actions.
var NavigationKeys = map[byte]bool{
	KeyUp:   true,
	KeyDown: true,
}
