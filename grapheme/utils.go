package grapheme

import (
	"strings"

	"github.com/clipperhouse/uax29/v2/graphemes"
)

func ClusterSetFromString(s string) ClusterSet {
	var pendingSet ClusterSet
	gr := graphemes.FromString(s)

	for gr.Next() {
		cluster := New([]rune(gr.Value())...)
		pendingSet = append(pendingSet, cluster)
	}

	var cs ClusterSet

	for len(pendingSet) > 0 {
		c := pendingSet[0]

		if !Equal(c, Escape) {
			cs = append(cs, c)
			pendingSet = pendingSet[1:]
			continue
		}

		pendingEscSeqSet := ClusterSet{}
		pendingEscSeqSet = append(pendingEscSeqSet, c)
		if len(pendingSet) == 1 {
			// We have an escape character but no more input, so we should treat
			// it as a normal cluster.
			cs = append(cs, c)
			pendingSet = pendingSet[1:]
			continue
		}

		pendingEscSeqSet = append(pendingEscSeqSet, pendingSet[1])
		pendingEscSeqCluster := New(pendingEscSeqSet.Runes()...)
		switch {
		case !IsFeEscapeSequence(pendingEscSeqCluster):
			// We have an escape character and one more input, but it's not a
			// full escape sequence, so we should treat the escape character as
			// a normal cluster and continue processing the next cluster.
			cs = append(cs, c)
			pendingSet = pendingSet[1:]
			continue
		}

		nextCluster := New(pendingSet[2:].Runes()...)
		seqEndIdx := IndexOfSequenceEnd(nextCluster)
		if seqEndIdx == -1 {
			// We have the start of an escape sequence but we don't have the
			// full sequence yet, so we should treat the escape character as a
			// normal cluster and continue processing the next cluster.
			cs = append(cs, c)
			pendingSet = pendingSet[1:]
			continue
		}

		pendingEscSeqCluster = append(pendingEscSeqCluster, pendingSet[2:3+seqEndIdx].Runes()...)
		// cs = append(cs, pendingEscSeqCluster)
		pendingSet = pendingSet[3+seqEndIdx:]
	}

	return cs
}

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
