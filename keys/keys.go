package keys

const (
	escape = 27
)

var (
	CtrlC       = New(3)
	Backspace   = New(8)
	Escape      = New(escape)
	Delete      = New(127)
	Enter       = New('\r')
	Newline     = New('\n')
	LeftBracket = New('[')

	LowerN = New('n')
	LowerY = New('y')

	UpperA = New('A')
	UpperN = New('N')
	UpperO = New('O')
	UpperY = New('Y')

	UpArrow    = New(escape, '[', 'A')
	DownArrow  = New(escape, '[', 'B')
	RightArrow = New(escape, '[', 'C')
	LeftArrow  = New(escape, '[', 'D')
)
