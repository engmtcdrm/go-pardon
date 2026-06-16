package pardon

import (
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/engmtcdrm/go-pardon/grapheme"
)

type BaseInputParser struct {
	pendingInput           []byte
	pendingInputRune       []rune
	pendingInputClusterSet grapheme.ClusterSet
	pendingEscSequence     grapheme.ClusterSet
}

// parseInputToGraphemeClusters parses pending input bytes to a
// [grapheme.ClusterSet] for further processing.
func (bp *BaseInputParser) parseInputToGraphemeClusters() (needMoreInput bool) {
	for len(bp.pendingInput) > 0 {
		r, width := utf8.DecodeRune(bp.pendingInput)
		// If we encounter an invalid UTF-8 sequence, we should wait for more input
		if r == utf8.RuneError {
			return true
		}

		bp.pendingInput = bp.pendingInput[width:]
		bp.pendingInputRune = append(bp.pendingInputRune, r)
	}

	pendingGraphemes := graphemes.FromString(string(bp.pendingInputRune))
	bp.pendingInputRune = nil

	for pendingGraphemes.Next() {
		cluster := grapheme.New([]rune(pendingGraphemes.Value())...)
		bp.pendingInputClusterSet = append(bp.pendingInputClusterSet, cluster)
	}

	return false
}

type InputHandler func(seq grapheme.Cluster) (done bool, err error)

func (bp *BaseInputParser) processEnter(c grapheme.Cluster, handler InputHandler) (done bool, err error) {
	return false, nil
}

func (bp *BaseInputParser) processEscapeSequence(c grapheme.Cluster, handler InputHandler) (done bool, err error) {
	defer func() {
		bp.pendingEscSequence = nil
	}()
	bp.pendingEscSequence = append(bp.pendingEscSequence, c)
	if len(bp.pendingInputClusterSet) == 1 {
		// We have an escape character but no more input, so we should wait
		// for more input before processing.
		return false, nil
	}
	nextCluster := bp.pendingInputClusterSet[1]
	bp.pendingEscSequence = append(bp.pendingEscSequence, nextCluster)
	pendingEscSequenceCluster := grapheme.New(bp.pendingEscSequence.Runes()...)
	if grapheme.IsFeEscapeSequence(pendingEscSequenceCluster) {
		if len(bp.pendingInputClusterSet) == 2 {
			// We have a complete escape sequence with only the escape character and the next rune,
			// so we should wait for more input before processing.
			return false, nil
		}

		nextCluster := grapheme.New(bp.pendingInputClusterSet[2:].Runes()...)
		seqEndIdx := grapheme.IndexOfSequenceEnd(nextCluster)
		if seqEndIdx == -1 {
			// We have the start of an escape sequence but we don't have the full sequence yet,
			// so we should wait for more input before processing.
			return false, nil
		}

		pendingEscSequenceCluster = append(pendingEscSequenceCluster, bp.pendingInputClusterSet[2:3+seqEndIdx].Runes()...)

		// Call the handler function with the complete escape sequence
		if handler != nil {
			handled, err := handler(pendingEscSequenceCluster)
			if handled {
				return true, err
			}
		}

		toRemove := 3 + seqEndIdx

		bp.pendingInputClusterSet = bp.pendingInputClusterSet[toRemove:]
	}
	return true, nil
}

type TestTest struct {
	BaseInputParser
}
