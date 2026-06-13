package keys

type Key []byte

func New(b ...byte) *Key {
	bs := Key(b)
	return &bs
}

func (b Key) Equal(input []byte) bool {
	if len(input) != len(b) {
		return false
	}

	for i := range b {
		if input[i] != b[i] {
			return false
		}
	}

	return true
}

func (b Key) String() string {
	return string(b)
}
