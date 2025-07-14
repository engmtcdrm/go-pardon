package gocliselect

const (
	// Key codes for navigation and actions
	KeyCtrlC          = byte(3)  // Ctrl+C generates ASCII 3 in raw mode
	KeyCarriageReturn = byte(10) // Additional for cross-platform compatibility
	KeyEnter          = byte(13)
	KeyEscape         = byte(27)
	KeyUp             = byte(65)
	KeyDown           = byte(66)
	KeyNoUpper        = byte(78)
	KeyYesUpper       = byte(89)
	KeyNo             = byte(110)
	KeyYes            = byte(121)
)
