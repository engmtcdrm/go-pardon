package pardon

// Form represents a collection of prompts executed sequentially.
type Form struct {
	prompts []Prompt
}

// NewForm creates a new Form with the given prompts.
func NewForm(prompts ...Prompt) *Form {
	f := &Form{
		prompts: prompts,
	}
	return f
}

// Ask executes all prompts in sequence, stopping on the first error.
func (f *Form) Ask() error {
	for _, p := range f.prompts {
		if err := p.Ask(); err != nil {
			return err
		}
	}
	return nil
}
