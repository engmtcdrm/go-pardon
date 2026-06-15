package pardon

import (
	"unicode/utf8"

	// "github.com/clipperhouse/uax29/graphemes"
	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/engmtcdrm/go-pardon/keys"
)

type BaseInputParser struct {
	pendingInput         []byte
	pendingInputRune     []rune
	pendingInputRuneKeys keys.Keys
	pendingEscSequence   keys.Keys
}

// parseInputToRuneKeys parses pending input bytes to rune keys for further
// processing.
func (bp *BaseInputParser) parseInputToRuneKeys() (needMoreInput bool) {
	for len(bp.pendingInput) > 0 {
		r, width := utf8.DecodeRune(bp.pendingInput)
		// If we encounter an invalid UTF-8 sequence, we should wait for more input
		if r == utf8.RuneError {
			return true
		}

		bp.pendingInput = bp.pendingInput[width:]
		bp.pendingInputRune = append(bp.pendingInputRune, r)
	}

	gr2 := graphemes.FromString(string(bp.pendingInputRune))
	bp.pendingInputRune = nil

	for gr2.Next() {
		runeKey := keys.New([]rune(gr2.Value())...)
		bp.pendingInputRuneKeys = append(bp.pendingInputRuneKeys, runeKey)
	}

	// gr := uniseg.NewGraphemes(string(bp.pendingInputRune))
	// bp.pendingInputRune = nil

	// Store each grapheme cluster as a rune key for further processing
	// for gr.Next() {
	// 	runeKey := keys.New(gr.Runes()...)
	// 	bp.pendingInputRuneKeys = append(bp.pendingInputRuneKeys, runeKey)
	// }

	return false
}

type KeyHandler func(seq keys.Key) (done bool, err error)

func (bp *BaseInputParser) processEnter(r keys.Key, handler KeyHandler) (done bool, err error) {
	return false, nil
}

func (bp *BaseInputParser) processEscapeSequence(r keys.Key, handler KeyHandler) (done bool, err error) {
	bp.pendingEscSequence = append(bp.pendingEscSequence, r)
	if len(bp.pendingInputRuneKeys) == 1 {
		// We have an escape character but no more input, so we should wait
		// for more input before processing.
		return false, nil
	}
	nextRune := bp.pendingInputRuneKeys[1]
	bp.pendingEscSequence = append(bp.pendingEscSequence, nextRune)
	pendingEscSequenceKey := keys.New(keys.RuneKeysToRunes(bp.pendingEscSequence)...)
	if keys.IsFeEscapeSequenceRune(pendingEscSequenceKey) {
		if len(bp.pendingInputRuneKeys) == 2 {
			// We have a complete escape sequence with only the escape character and the next rune,
			// so we should wait for more input before processing.
			return false, nil
		}

		tempRunes := keys.New(keys.RuneKeysToRunes(bp.pendingInputRuneKeys[2:])...)
		seqEndIdx := keys.IndexOfSequenceEndRune(tempRunes)
		if seqEndIdx == -1 {
			// We have the start of an escape sequence but we don't have the full sequence yet,
			// so we should wait for more input before processing.
			return false, nil
		}

		pendingEscSequenceKey = append(pendingEscSequenceKey, keys.RuneKeysToRunes(bp.pendingInputRuneKeys[2:3+seqEndIdx])...)

		// Call the handler function with the complete escape sequence
		if handler != nil {
			handled, err := handler(pendingEscSequenceKey)
			if handled {
				return true, err
			}
		}

		toRemove := 3 + seqEndIdx

		bp.pendingInputRuneKeys = bp.pendingInputRuneKeys[toRemove:]
	}
	return false, nil
}

type TestTest struct {
	BaseInputParser
}
