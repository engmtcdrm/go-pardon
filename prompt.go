package gocliselect

type PromptInterface[T any] interface {
	Title(title T) PromptInterface[T]
	TitleFunc(fn func() T) PromptInterface[T]
	Value(value *T) PromptInterface[T]
	Icon(value *T) PromptInterface[T]
	IconFunc(fn func() T) PromptInterface[T]
	Validate(fn func(T) error) PromptInterface[T]
	Ask() error
}
