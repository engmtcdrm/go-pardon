package grapheme

import "github.com/mattn/go-runewidth"

// Cluster represents a grapheme cluster in a sequence of runes that form a
// single user-perceived character, such as an emoji or a combined character. It
// is used to handle complex input scenarios where multiple runes may be
// combined into a single visual unit.
type Cluster []rune

// New creates a new [Cluster] from the provided runes.
func New(r ...rune) Cluster {
	return Cluster(r)
}

// NewFromRunes creates a new [Cluster] by combining multiple slices of runes
// into a single Cluster.
func NewFromRunes(runes ...[]rune) Cluster {
	var combinedRunes []rune
	for _, r := range runes {
		combinedRunes = append(combinedRunes, r...)
	}

	return Cluster(combinedRunes)
}

// NewFromString creates a new [Cluster] from the provided string.
func NewFromString(s string) Cluster {
	return Cluster([]rune(s))
}

// Bytes returns the byte representation of the Cluster.
func (c Cluster) Bytes() []byte {
	return []byte(c.String())
}

// IsANSIEscapeSequence checks if the Cluster represents an ANSI escape
// sequence. This may not be a valid ANSI escape sequence, but it will validate
// the start being an FE escape sequence and the end being a valid ANSI sequence
// terminator.
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

// Len returns the visual length of the cluster, including ANSI escape
// sequences.
func (c Cluster) Len() int {
	width := 0
	for _, r := range c {
		width += runewidth.RuneWidth(r)
	}

	return width
}

// String returns the string representation of the Cluster.
func (c Cluster) String() string {
	return string(c)
}

// VisualLen returns the visual length of the cluster, ignorning ANSI escape
// sequences.
func (c Cluster) VisualLen() int {
	if c.IsANSIEscapeSequence() {
		return 0
	}

	return c.Len()
}
