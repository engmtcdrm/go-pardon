package grapheme

import "strings"

// Equal reports whether a and b are the same length and contain the same runes.
// A nil argument is equivalent to an empty slice.
func Equal(a Cluster, b Cluster) bool {
	return string(a) == string(b)
}

// EqualFold reports whether a and b, interpreted as UTF-8 strings,
// are equal under simple Unicode case-folding, which is a more general
// form of case-insensitivity.
func EqualFold(a Cluster, b Cluster) bool {
	return strings.EqualFold(string(a), string(b))
}

// HasSequenceEnd checks if the given Cluster contains a ANSI escape sequence
// end character.
func HasSequenceEnd(gc Cluster) bool {
	for _, b := range gc {
		if b >= 0x40 && b <= 0x7E {
			return true
		}
	}

	return false
}

// IndexOfSequenceEnd returns the index of the ANSI escape sequence end
// character in the given Cluster, or -1 if there is no such character.
func IndexOfSequenceEnd(gc Cluster) int {
	for i, b := range gc {
		if b >= 0x40 && b <= 0x7E {
			return i
		}
	}

	return -1
}

// IsFeEscapeSequence checks if the given Cluster is a Fe Escape Sequence.
func IsFeEscapeSequence(gc Cluster) bool {
	if len(gc) != 2 {
		return false
	}

	if gc[0] != escape {
		return false
	}

	if gc[1] < 0x40 || gc[1] > 0x5F {
		return false
	}

	return true
}

// IsSequenceEnd checks if the given Cluster is a ANSI escape sequence end
// character.
func IsSequenceEnd(gc Cluster) bool {
	if len(gc) != 1 {
		return false
	}

	if gc[0] < 0x40 || gc[0] > 0x7E {
		return false
	}

	return true
}
