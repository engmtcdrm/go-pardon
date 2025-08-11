package pardon

// eval holds a value that can be statically set or dynamically computed.
type eval[T any] struct {
	val       T         // Static value
	fn        func(T) T // Optional transformation function
	defaultFn func(T) T // Fallback transformation function
}

// Get returns the evaluated value, applying fn or defaultFn if set.
func (d *eval[T]) Get() T {
	if d.fn != nil {
		return d.fn(d.val)
	}

	if d.defaultFn != nil {
		return d.defaultFn(d.val)
	}

	return d.val
}
