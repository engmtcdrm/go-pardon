package grapheme

import (
	"strings"

	"github.com/clipperhouse/uax29/v2/graphemes"
)

// ClusterSetFromString converts a string into a [ClusterSet]. Each [Cluster]
// represents a character, a grapheme, or an ANSI escape sequence each as a
// single cluster.
func ClusterSetFromString(s string) ClusterSet {
	var pendingSet ClusterSet
	gr := graphemes.FromString(s)

	for gr.Next() {
		cluster := New([]rune(gr.Value())...)
		pendingSet = append(pendingSet, cluster)
	}

	return clusterANSIEscape(pendingSet)
}

// Equal reports whether a and b are the same length and contain the same runes.
// A nil argument is equivalent to an empty slice.
func Equal(a Cluster, b Cluster) bool {
	if len(a) != len(b) {
		return false
	}

	return string(a) == string(b)
}

// EqualFold reports whether a and b, interpreted as UTF-8 strings,
// are equal under simple Unicode case-folding, which is a more general
// form of case-insensitivity.
func EqualFold(a Cluster, b Cluster) bool {
	if len(a) != len(b) {
		return false
	}

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
	if len(gc) < 2 {
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

// clusterANSIEscape processes a [ClusterSet] to identify and group ANSI escape
// sequences into single clusters. It returns a new [ClusterSet] where ANSI escape
// sequences are treated as single clusters, while non-escape clusters are left
// unchanged.
func clusterANSIEscape(in ClusterSet) ClusterSet {
	var out ClusterSet

	for len(in) > 0 {
		c := in[0]

		if !Equal(c, Escape) {
			out = append(out, c)
			in = in[1:]
			continue
		}

		pendingEscSeqSet := ClusterSet{}
		pendingEscSeqSet = append(pendingEscSeqSet, c)
		if len(in) == 1 {
			// We have an escape character but no more input, so we should treat
			// it as a normal cluster.
			out = append(out, c)
			in = in[1:]
			continue
		}

		pendingEscSeqSet = append(pendingEscSeqSet, in[1])
		pendingEscSeqCluster := New(pendingEscSeqSet.Runes()...)
		switch {
		case !IsFeEscapeSequence(pendingEscSeqCluster):
			// We have an escape character and one more input, but it's not a
			// full escape sequence, so we should treat the escape character as
			// a normal cluster and continue processing the next cluster.
			out = append(out, c)
			in = in[1:]
			continue
		}

		nextCluster := New(in[2:].Runes()...)
		seqEndIdx := IndexOfSequenceEnd(nextCluster)
		if seqEndIdx == -1 {
			// We have the start of an escape sequence but we don't have the
			// full sequence yet, so we should treat the escape character as a
			// normal cluster and continue processing the next cluster.
			out = append(out, c)
			in = in[1:]
			continue
		}

		pendingEscSeqCluster = append(pendingEscSeqCluster, in[2:3+seqEndIdx].Runes()...)
		out = append(out, pendingEscSeqCluster)
		in = in[3+seqEndIdx:]
	}
	return out
}
