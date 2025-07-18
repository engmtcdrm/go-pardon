package pardon

type promptInterface[T any] interface {
	Title(title T) promptInterface[T]
	TitleFunc(fn func() T) promptInterface[T]
	Value(value *T) promptInterface[T]
	Icon(value *T) promptInterface[T]
	IconFunc(fn func() T) promptInterface[T]
	Validate(fn func(T) error) promptInterface[T]
	AnswerFunc(fn func(T) T) promptInterface[T]
	Ask() error
}
