package pardon

import (
	"io"
	"os"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/engmtcdrm/go-pardon/grapheme"
)

type PromptBase[T comparable, Self any] struct {
	// Self holds a pointer to the concrete self struct so fluent
	// methods on the base can return the self type.
	Self Self

	// Out is the output writer for the terminal, typically [os.Stdout].
	Out io.Writer

	// In is the terminal input reader.
	In *TerminalInput

	icon   eval[string]
	title  eval[string]
	answer eval[string]
	prompt string
	value  *T

	pendingInputBytes      []byte
	pendingInputRunes      []rune
	pendingInputClusterSet grapheme.ClusterSet
}

// NewPromptBase creates and initializes a new PromptBase instance with the
// given default value.
func NewPromptBase[T comparable, Self any](value *T) PromptBase[T, Self] {
	return PromptBase[T, Self]{
		Self:   zeroParent[Self](),
		In:     NewTerminalInput(),
		Out:    os.Stdout,
		icon:   eval[string]{val: Icons.QuestionMark, fn: nil, defaultFn: defaultFuncs.iconFn},
		title:  eval[string]{val: "", fn: nil, defaultFn: defaultFuncs.titleFn},
		answer: eval[string]{val: "", fn: nil, defaultFn: defaultFuncs.answerFn},
		value:  value,
	}
}

// AnswerFunc sets a function to transform the final answer being displayed.
func (bp *PromptBase[T, Self]) AnswerFunc(fn func(string) string) Self {
	bp.answer.fn = fn
	return bp.Self
}

// Icon sets a static icon for the confirmation prompt.
func (bp *PromptBase[T, Self]) Icon(s string) Self {
	bp.icon.val = s
	bp.icon.fn = nil
	return bp.Self
}

// IconFunc sets a dynamic icon function for the confirmation prompt.
func (bp *PromptBase[T, Self]) IconFunc(fn func(string) string) Self {
	bp.icon.fn = fn
	return bp.Self
}

// Title sets a static title for the confirmation prompt.
func (bp *PromptBase[T, Self]) Title(title string) Self {
	bp.title.val = title
	bp.title.fn = nil
	return bp.Self
}

// TitleFunc sets a dynamic title function for the confirmation prompt.
func (bp *PromptBase[T, Self]) TitleFunc(fn func(string) string) Self {
	bp.title.fn = fn
	return bp.Self
}

// Value sets a default value for the confirmation prompt.
func (bp *PromptBase[T, Self]) Value(value *T) Self {
	bp.value = value
	return bp.Self
}

// parseInputToGraphemeSet parses pending input bytes to a [grapheme.ClusterSet]
// for further processing.
func (bp *PromptBase[T, Self]) parseInputToGraphemeSet() (needMoreInput bool) {
	for len(bp.pendingInputBytes) > 0 {
		r, width := utf8.DecodeRune(bp.pendingInputBytes)
		// If we encounter an invalid UTF-8 sequence, we should wait for more input
		if r == utf8.RuneError {
			return true
		}

		bp.pendingInputBytes = bp.pendingInputBytes[width:]
		bp.pendingInputRunes = append(bp.pendingInputRunes, r)
	}

	pendingGraphemes := graphemes.FromString(string(bp.pendingInputRunes))
	bp.pendingInputRunes = nil

	for pendingGraphemes.Next() {
		cluster := grapheme.New([]rune(pendingGraphemes.Value())...)
		bp.pendingInputClusterSet = append(bp.pendingInputClusterSet, cluster)
	}

	return false
}

type InputHandler func(c grapheme.Cluster) (done bool, err error)

func (bp *PromptBase[T, Self]) processEscapeSequence(c grapheme.Cluster, handler InputHandler) (done bool, err error) {
	pendingEscSeqSet := grapheme.ClusterSet{}

	// Handle first grapheme of the escape sequence.
	pendingEscSeqSet = append(pendingEscSeqSet, c)
	if len(bp.pendingInputClusterSet) == 1 {
		// We have an escape character but no more input, so we should wait for
		// more input before processing.
		return false, nil
	}

	// Handle second grapheme of the escape sequence.
	pendingEscSeqSet = append(pendingEscSeqSet, bp.pendingInputClusterSet[1])
	pendingEscSeqCluster := grapheme.New(pendingEscSeqSet.Runes()...)
	switch {
	case !grapheme.IsFeEscapeSequence(pendingEscSeqCluster):
		return true, nil
	case len(bp.pendingInputClusterSet) == 2:
		// We have an escape character and one more input, but it's not a full
		// escape sequence, so we should wait for more input before processing.
		return false, nil
	}

	// Handle any additional graphemes that may be part of the escape sequence.
	nextCluster := grapheme.New(bp.pendingInputClusterSet[2:].Runes()...)
	seqEndIdx := grapheme.IndexOfSequenceEnd(nextCluster)
	if seqEndIdx == -1 {
		// We have the start of an escape sequence but we don't have the full
		// sequence yet, so we should wait for more input before processing.
		return false, nil
	}

	pendingEscSeqCluster = append(pendingEscSeqCluster, bp.pendingInputClusterSet[2:3+seqEndIdx].Runes()...)

	if handler != nil {
		handled, err := handler(pendingEscSeqCluster)
		if handled {
			return true, err
		}
	}

	bp.pendingInputClusterSet = bp.pendingInputClusterSet[3+seqEndIdx:]

	return true, nil
}

// zeroParent returns the zero value for the generic type P. This is used to
// initialize the Self field in the BasePrompt struct to a zero value of the
// concrete parent type.
func zeroParent[P any]() P {
	var p P
	return p
}
