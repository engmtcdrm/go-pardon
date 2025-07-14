package gocliselect

type Item[T comparable] struct {
	Key       string
	Value     T
	SubSelect *Select[T]
}

func NewSelectItem[T comparable](key string, value T) Item[T] {
	return Item[T]{
		Key:   key,
		Value: value,
	}
}
