package pardon

// Option represents a selectable key-value pair for use in selection prompts.
type Option[T comparable] struct {
	Key   string // Display label
	Value T      // Associated value
}

// NewOption creates a new Option with the given key and value.
func NewOption[T comparable](key string, value T) Option[T] {
	return Option[T]{
		Key:   key,
		Value: value,
	}
}
