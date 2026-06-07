package pardon

// eval holds a value that can be statically set or dynamically computed.
type eval[T any] struct {
	// Static value.
	val T

	// Optional transformation function.
	fn func(T) T

	// Fallback transformation function.
	defaultFn func(T) T
}

// Get returns the evaluated value, applying the transformation function if it
// exists. If no transformation function is set, it returns the static value.
func (d *eval[T]) Get() T {
	if d.fn != nil {
		return d.fn(d.val)
	}

	if d.defaultFn != nil {
		return d.defaultFn(d.val)
	}

	return d.val
}
