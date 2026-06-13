package keys

func IsFeEscapeSequence(bytes []byte) bool {
	if len(bytes) != 2 {
		return false
	}

	if bytes[0] != Escape {
		return false
	}

	if bytes[1] < 0x40 || bytes[1] > 0x5F {
		return false
	}

	return true
}
