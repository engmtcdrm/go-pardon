package pardon

// eval is a struct that holds a value, an optional function to evaluate, and an optional default function.
type eval[T any] struct {
	val       T
	fn        func(T) T
	defaultFn func(T) T
}

// Get returns the evaluated value based on the provided function or the default function.
// If no function is provided, it returns the original value.
func (d *eval[T]) Get() T {
	if d.fn != nil {
		return d.fn(d.val)
	}

	if d.defaultFn != nil {
		return d.defaultFn(d.val)
	}

	return d.val
}
