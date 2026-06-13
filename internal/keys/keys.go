package keys

const (
	escape = byte(27)

	leftBracket = byte('[')

	upperA = byte('A')
	upperB = byte('B')
	upperC = byte('C')
	upperD = byte('D')
	upperN = byte('N')
	upperO = byte('O')
	upperY = byte('Y')
)

var (
	CtrlC       = New(3)
	Backspace   = New(8)
	Escape      = New(27)
	Delete      = New(127)
	Enter       = New('\r')
	Newline     = New('\n')
	LeftBracket = New(leftBracket)

	LowerN = New('n')
	LowerY = New('y')

	UpperA = New(upperA)
	UpperN = New(upperN)
	UpperO = New(upperO)
	UpperY = New(upperY)

	UpArrow    = New(escape, leftBracket, upperA)
	DownArrow  = New(escape, leftBracket, upperB)
	RightArrow = New(escape, leftBracket, upperC)
	LeftArrow  = New(escape, leftBracket, upperD)
)
