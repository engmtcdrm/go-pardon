package pardon

// Prompt is the interface implemented by all prompt types.
type Prompt interface {
	// Ask displays the prompt and waits for user input.
	Ask() error
}
