package grapheme

const (
	escape = '\x1b'
)

var (
	CtrlC     = New(3)
	Backspace = New(8)
	Escape    = New(escape)
	CSI       = NewFromRunes(Escape, New('['))
	Delete    = New(127)
	Enter     = New('\r')
	Newline   = New('\n')

	UpArrow    = NewFromRunes(CSI, New('A'))
	DownArrow  = NewFromRunes(CSI, New('B'))
	RightArrow = NewFromRunes(CSI, New('C'))
	LeftArrow  = NewFromRunes(CSI, New('D'))
)
