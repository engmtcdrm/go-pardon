package pardon

type Prompt interface {
	// Ask prompts the user for input and stores the result.
	Ask() error
}
