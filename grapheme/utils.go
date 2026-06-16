package grapheme

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

func IsSequenceEnd(gc Cluster) bool {
	if len(gc) != 1 {
		return false
	}

	if gc[0] < 0x40 || gc[0] > 0x7E {
		return false
	}

	return true
}

func HasSequenceEnd(gc Cluster) bool {
	for _, b := range gc {
		if b >= 0x40 && b <= 0x7E {
			return true
		}
	}

	return false
}

func IndexOfSequenceEnd(gc Cluster) int {
	for i, b := range gc {
		if b >= 0x40 && b <= 0x7E {
			return i
		}
	}

	return -1
}
