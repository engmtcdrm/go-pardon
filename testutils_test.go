package pardon

import "testing"

type funcType int

const (
	answerFn funcType = iota
	cursorFn
	iconFn
	selectFn
	titleFn
)

// changeDefaultFunc is a helper function to temporarily change one of the
// default functions for testing purposes. It saves the original function,
// replaces it with the new function, and restores the original function after
// the test completes using t.Cleanup.
func changeDefaultFunc(t *testing.T, fnType funcType, newFn func(string) string) {
	var originalFn func(string) string

	switch fnType {
	case answerFn:
		originalFn = defaultFuncs.answerFn
		defaultFuncs.answerFn = newFn
	case cursorFn:
		originalFn = defaultFuncs.cursorFn
		defaultFuncs.cursorFn = newFn
	case iconFn:
		originalFn = defaultFuncs.iconFn
		defaultFuncs.iconFn = newFn
	case selectFn:
		originalFn = defaultFuncs.selectFn
		defaultFuncs.selectFn = newFn
	case titleFn:
		originalFn = defaultFuncs.titleFn
		defaultFuncs.titleFn = newFn
	default:
		t.Fatalf("Unknown function type: %v", fnType)
	}

	t.Cleanup(func() {
		switch fnType {
		case answerFn:
			defaultFuncs.answerFn = originalFn
		case cursorFn:
			defaultFuncs.cursorFn = originalFn
		case iconFn:
			defaultFuncs.iconFn = originalFn
		case selectFn:
			defaultFuncs.selectFn = originalFn
		case titleFn:
			defaultFuncs.titleFn = originalFn
		}
	})
}
