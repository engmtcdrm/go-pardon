package keys

const (
	CtrlC     = byte(3)
	Backspace = byte(8)
	Escape    = byte(27)
	Delete    = byte(127)

	Up    = UpperA
	Down  = UpperB
	Right = UpperC
	Left  = UpperD

	Enter       = byte('\r')
	LeftBracket = byte('[')
	LowerN      = byte('n')
	LowerY      = byte('y')
	NewLine     = byte('\n') // Additional for cross-platform compatibility

	UpperA = byte('A')
	UpperB = byte('B')
	UpperC = byte('C')
	UpperD = byte('D')
	UpperN = byte('N')
	UpperO = byte('O')
	UpperY = byte('Y')
)

var (
	UpArrow    = []byte{Escape, LeftBracket, UpperA}
	DownArrow  = []byte{Escape, LeftBracket, UpperB}
	RightArrow = []byte{Escape, LeftBracket, UpperC}
	LeftArrow  = []byte{Escape, LeftBracket, UpperD}
)
