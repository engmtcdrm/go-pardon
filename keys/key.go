package keys

type Key []rune

func New(r ...rune) Key {
	return Key(r)
}

func (rk Key) String() string {
	return string(rk)
}

func (rk Key) Bytes() []byte {
	return []byte(rk.String())
}

func RuneKeysToRunes(rk []Key) []rune {
	var result []rune
	for _, r := range rk {
		result = append(result, r...)
	}
	return result
}

// Keys represents a slice of [Key], providing a way to handle multiple
// Key instances as a single entity.
type Keys []Key

func (rks Keys) String() string {
	return string(rks.Bytes())
}

func (rks Keys) Bytes() []byte {
	var result []byte
	for _, rk := range rks {
		result = append(result, rk.Bytes()...)
	}

	return result
}
