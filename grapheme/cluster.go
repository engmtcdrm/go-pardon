package grapheme

// Cluster represents a grapheme cluster in a sequence of runes that form a
// single user-perceived character, such as an emoji or a combined character. It
// is used to handle complex input scenarios where multiple runes may be
// combined into a single visual unit.
type Cluster []rune

func New(r ...rune) Cluster {
	return Cluster(r)
}

// Bytes returns the byte representation of the Cluster.
func (c Cluster) Bytes() []byte {
	return []byte(c.String())
}

// String returns the string representation of the Cluster.
func (c Cluster) String() string {
	return string(c)
}
