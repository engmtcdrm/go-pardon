package pardon

type funcs struct {
	answerFn func(string) string
	cursorFn func(string) string
	iconFn   func(string) string
	selectFn func(string) string
	titleFn  func(string) string
}

var defaultFuncs = funcs{
	answerFn: nil,
	cursorFn: nil,
	iconFn:   nil,
	selectFn: nil,
	titleFn:  nil,
}

func SetDefaultAnswerFunc(fn func(string) string) {
	defaultFuncs.answerFn = fn
}

func SetDefaultCursorFunc(fn func(string) string) {
	defaultFuncs.cursorFn = fn
}

func SetDefaultIconFunc(fn func(string) string) {
	defaultFuncs.iconFn = fn
}

func SetDefaultSelectFunc(fn func(string) string) {
	defaultFuncs.selectFn = fn
}

func SetDefaultTitleFunc(fn func(string) string) {
	defaultFuncs.titleFn = fn
}
