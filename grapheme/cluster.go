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

func (c Cluster) IsANSIEscapeSequence() bool {
	if !IsFeEscapeSequence(c) {
		return false
	}

	seqEndIdx := IndexOfSequenceEnd(c[2:])
	if seqEndIdx == -1 {
		return false
	}

	if seqEndIdx != len(c[2:])-1 {
		return false
	}

	return true
}

func (c Cluster) Len() int {
	return len(c)
}

func (c Cluster) VisualLen() int {
	if c.IsANSIEscapeSequence() {
		return 0
	}

	return len(c)
}

// String returns the string representation of the Cluster.
func (c Cluster) String() string {
	return string(c)
}
