package pardon

// Form is a struct that holds a list of prompts and their answers.
type Form struct {
	prompts []Prompt
}

// NewForm creates a new Form with the provided prompts.
func NewForm(prompts ...Prompt) *Form {
	f := &Form{
		prompts: prompts,
	}
	return f
}

// Ask iterates through each prompt in the form and asks for user input.
func (f *Form) Ask() error {
	for _, p := range f.prompts {
		if err := p.Ask(); err != nil {
			return err
		}
	}
	return nil
}
