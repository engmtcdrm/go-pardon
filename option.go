package gocliselect

type Option[T comparable] struct {
	Key   string
	Value T
}

func NewOption[T comparable](key string, value T) Option[T] {
	return Option[T]{
		Key:   key,
		Value: value,
	}
}
