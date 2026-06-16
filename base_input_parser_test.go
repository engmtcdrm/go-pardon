package pardon

import (
	"testing"
)

func TestThis(t *testing.T) {
	tt := TestTest{}
	tt.pendingInput = []byte("example input ❤️")
	tt.parseInputToGraphemeClusters()

	t.Logf("%v", tt.pendingInputClusterSet)
}
