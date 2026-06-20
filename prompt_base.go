package pardon

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/engmtcdrm/go-pardon/grapheme"
	"golang.org/x/term"
)

// ClusterHandler defines a function type that handles a [grapheme.Cluster].
// It returns a boolean indicating whether processing is done and an error if
// any occurred.
type ClusterHandler func(c grapheme.Cluster) (done bool, err error)

// PromptBase is a generic base struct for creating terminal prompts.
// It provides common functionality such as handling input and output,
// managing the prompt's icon, title, and answer, and processing escape
// sequences.
type PromptBase[T comparable, Self any] struct {
	// Self holds a pointer to the concrete self struct so fluent
	// methods on the base can return the self type.
	Self Self

	// Out is the output writer for the terminal, typically [os.Stdout].
	Out io.Writer

	// In is the terminal input reader.
	In TerminalInput

	icon   eval[string]
	title  eval[string]
	answer eval[string]
	value  *T

	// The prompt string to display to the user.
	//
	// This should not be modified outside the struct. It is intentionally
	// exported to allow for easier testing.
	Prompt grapheme.ClusterSet

	// PendingInputBytes holds the raw input bytes that have been read from the
	// terminal but not yet processed.
	//
	// This should not be modified outside the struct. It is intentionally
	// exported to allow for easier testing.
	PendingInputBytes []byte

	// PendingInputRunes holds the decoded runes from
	// [PromptBase.PendingInputBytes].
	//
	// This should not be modified outside the struct. It is intentionally
	// exported to allow for easier testing.
	PendingInputRunes []rune

	// PendingInputClusterSet holds the grapheme clusters parsed from
	// [PromptBase.PendingInputRunes].
	//
	// This should not be modified outside the struct. It is intentionally
	// exported to allow for easier testing.
	PendingInputClusterSet grapheme.ClusterSet
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

// Icon sets a static icon for the prompt.
func (bp *PromptBase[T, Self]) Icon(s string) Self {
	bp.icon.val = s
	bp.icon.fn = nil
	return bp.Self
}

// IconFunc sets a dynamic icon function for the prompt.
func (bp *PromptBase[T, Self]) IconFunc(fn func(string) string) Self {
	bp.icon.fn = fn
	return bp.Self
}

// Title sets a static title for the prompt.
func (bp *PromptBase[T, Self]) Title(title string) Self {
	bp.title.val = title
	bp.title.fn = nil
	return bp.Self
}

// TitleFunc sets a dynamic title function for the prompt.
func (bp *PromptBase[T, Self]) TitleFunc(fn func(string) string) Self {
	bp.title.fn = fn
	return bp.Self
}

// Value sets a default value for the prompt.
func (bp *PromptBase[T, Self]) Value(value *T) Self {
	bp.value = value
	return bp.Self
}

// ConvertBytesToGraphemeSet parses pending input bytes to a
// [grapheme.ClusterSet] for further processing.
func (bp *PromptBase[T, Self]) ConvertBytesToGraphemeSet() (needMoreInput bool) {
	for len(bp.PendingInputBytes) > 0 {
		r, width := utf8.DecodeRune(bp.PendingInputBytes)
		// If we encounter an invalid UTF-8 sequence, we should wait for more
		// input before processing.
		if r == utf8.RuneError {
			return true
		}

		bp.PendingInputBytes = bp.PendingInputBytes[width:]
		bp.PendingInputRunes = append(bp.PendingInputRunes, r)
	}

	pendingGraphemes := graphemes.FromString(string(bp.PendingInputRunes))
	bp.PendingInputRunes = nil

	for pendingGraphemes.Next() {
		cluster := grapheme.New([]rune(pendingGraphemes.Value())...)
		bp.PendingInputClusterSet = append(bp.PendingInputClusterSet, cluster)
	}

	return false
}

// GetTerminalSize returns the width and height of the terminal in columns and
// rows. If the terminal size cannot be determined, it returns default values of
// 80 columns and 25 rows.
func (bp *PromptBase[T, Self]) GetTerminalSize() (width int, height int) {
	termWidth := 80  // Default width
	termHeight := 25 // Default height

	f, ok := bp.Out.(*os.File)
	if !ok {
		return termWidth, termHeight
	}

	width, height, err := term.GetSize(int(f.Fd()))
	if err == nil {
		return width, height
	}

	return termWidth, termHeight
}

func (bp *PromptBase[T, Self]) ProcessEscapeSequence(c grapheme.Cluster, handler ClusterHandler) (done bool, err error) {
	pendingEscSeqSet := grapheme.ClusterSet{}

	// Handle first grapheme of the escape sequence.
	pendingEscSeqSet = append(pendingEscSeqSet, c)
	if len(bp.PendingInputClusterSet) == 1 {
		// We have an escape character but no more input, so we should wait for
		// more input before processing.
		return false, nil
	}

	// Handle second grapheme of the escape sequence.
	pendingEscSeqSet = append(pendingEscSeqSet, bp.PendingInputClusterSet[1])
	pendingEscSeqCluster := grapheme.New(pendingEscSeqSet.Runes()...)
	switch {
	case !grapheme.IsFeEscapeSequence(pendingEscSeqCluster):
		return true, nil
	case len(bp.PendingInputClusterSet) == 2:
		// We have an escape character and one more input, but it's not a full
		// escape sequence, so we should wait for more input before processing.
		return false, nil
	}

	// Handle any additional graphemes that may be part of the escape sequence.
	nextCluster := grapheme.New(bp.PendingInputClusterSet[2:].Runes()...)
	seqEndIdx := grapheme.IndexOfSequenceEnd(nextCluster)
	if seqEndIdx == -1 {
		// We have the start of an escape sequence but we don't have the full
		// sequence yet, so we should wait for more input before processing.
		return false, nil
	}

	pendingEscSeqCluster = append(pendingEscSeqCluster, bp.PendingInputClusterSet[2:3+seqEndIdx].Runes()...)

	if handler != nil {
		handled, err := handler(pendingEscSeqCluster)
		if handled {
			return true, err
		}
	}

	bp.PendingInputClusterSet = bp.PendingInputClusterSet[3+seqEndIdx:]

	return true, nil
}

func (bp *PromptBase[T, Self]) buildAndSetPrompt() {
	bp.Prompt = grapheme.ClusterSetFromString(fmt.Sprintf("%s%s ", bp.icon.Get(), bp.title.Get()))
}
