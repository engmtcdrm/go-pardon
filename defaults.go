package pardon

// icons holds the default icon strings used across different prompt types.
// These provide consistent visual cues for different types of prompts.
type icons struct {
	Alert        string // Icon for alert/warning prompts
	QuestionMark string // Icon for question prompts and generic confirmations
	Password     string // Icon for password input prompts
}

// Icons contains default icons for different prompt types.
var Icons = icons{
	Alert:        "[!] ",
	QuestionMark: "[?] ",
	Password:     "🔒 ",
}

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

// SetDefaultAnswerFunc sets the global default answer transformation function.
func SetDefaultAnswerFunc(fn func(string) string) {
	defaultFuncs.answerFn = fn
}

// SetDefaultCursorFunc sets the global default cursor formatting function.
func SetDefaultCursorFunc(fn func(string) string) {
	defaultFuncs.cursorFn = fn
}

// SetDefaultIconFunc sets the global default icon transformation function.
func SetDefaultIconFunc(fn func(string) string) {
	defaultFuncs.iconFn = fn
}

// SetDefaultSelectFunc sets the global default selection formatting function.
func SetDefaultSelectFunc(fn func(string) string) {
	defaultFuncs.selectFn = fn
}

// SetDefaultTitleFunc sets the global default title transformation function.
func SetDefaultTitleFunc(fn func(string) string) {
	defaultFuncs.titleFn = fn
}
