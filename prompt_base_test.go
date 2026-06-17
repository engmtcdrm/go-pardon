package pardon

import (
	"testing"
)

type TestTest[T comparable] struct {
	PromptBase[T, *TestTest[T]]
}

func TestThis(t *testing.T) {
	tt := TestTest[string]{}
	tt.PendingInputBytes = []byte("example input ❤️")
	tt.ConvertBytesToGraphemeSet()

	t.Logf("%v", tt.PendingInputClusterSet)
}
