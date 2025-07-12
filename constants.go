package gocliselect

const (

	// Key codes for navigation and actions
	KeyUp             = byte(65)
	KeyDown           = byte(66)
	KeyEscape         = byte(27)
	KeyEnter          = byte(13)
	KeyCarriageReturn = byte(10) // Additional for cross-platform compatibility
	KeyCtrlC          = byte(3)  // Ctrl+C generates ASCII 3 in raw mode
)

// NavigationKeys defines a map of specific byte keycodes related to
// navigation functionality, such as up or down actions.
var NavigationKeys = map[byte]bool{
	KeyUp:   true,
	KeyDown: true,
}
