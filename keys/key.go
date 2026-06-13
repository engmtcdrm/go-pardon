package keys

import "bytes"

type Key []byte

func New(b ...byte) Key {
	return Key(b)
}

func (k Key) String() string {
	return string(k)
}

func (k Key) Runes() []rune {
	return bytes.Runes(k)
}
