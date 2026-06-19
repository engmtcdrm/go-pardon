package grapheme

// ClusterSet represents a set of [Cluster]s.
type ClusterSet []Cluster

// Bytes returns the byte representation of the ClusterSet.
func (cs ClusterSet) Bytes() []byte {
	var result []byte
	for _, rk := range cs {
		result = append(result, rk.Bytes()...)
	}

	return result
}

func (cs ClusterSet) Len() int {
	return len(cs)
}

// Runes returns the rune representation of the ClusterSet.
func (cs ClusterSet) Runes() []rune {
	var result []rune
	for _, rk := range cs {
		result = append(result, rk...)
	}

	return result
}

// String returns the string representation of the ClusterSet.
func (cs ClusterSet) String() string {
	return string(cs.Runes())
}

func (cs ClusterSet) VisualLen() int {
	total := 0
	for _, c := range cs {
		total += c.VisualLen()
	}

	return total
}
