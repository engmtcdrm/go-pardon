package keys

func IsFeEscapeSequence(key Key) bool {
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
