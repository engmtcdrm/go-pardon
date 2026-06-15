package keys

func IsFeEscapeSequenceRune(key Key) bool {
	if len(key) != 2 {
		return false
	}

	if key[0] != escape {
		return false
	}

	if key[1] < 0x40 || key[1] > 0x5F {
		return false
	}

	return true
}

func IsSequenceEndRune(key Key) bool {
	if len(key) != 1 {
		return false
	}

	if key[0] < 0x40 || key[0] > 0x7E {
		return false
	}

	return true
}

func HasSequenceEndRune(key Key) bool {
	for _, b := range key {
		if b >= 0x40 && b <= 0x7E {
			return true
		}
	}

	return false
}

func IndexOfSequenceEndRune(key Key) int {
	for i, b := range key {
		if b >= 0x40 && b <= 0x7E {
			return i
		}
	}

	return -1
}
