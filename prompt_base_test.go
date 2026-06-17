package pardon

import (
	"testing"
)

type TestTest[T comparable] struct {
	PromptBase[T, *TestTest[T]]
}

func TestThis(t *testing.T) {
	tt := TestTest[string]{}
	tt.pendingInputBytes = []byte("example input ❤️")
	tt.parseInputToGraphemeSet()

	t.Logf("%v", tt.pendingInputClusterSet)
}
