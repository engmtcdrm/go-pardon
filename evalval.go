package pardon

// evalVal is a struct that holds a value and an optional function to retrieve it dynamically.
type evalVal[T any] struct {
	val       T
	fn        func(T) T
	defaultFn func(T) T
}

// Get retrieves the value.
// If a function is set, it calls that function to get the value; otherwise, it
// returns the stored value.
func (d *evalVal[T]) Get() T {
	if d.fn != nil {
		return d.fn(d.val)
	}

	if d.defaultFn != nil {
		return d.defaultFn(d.val)
	}

	return d.val
}
