package pardon

import "errors"

var (
	ErrNoSelectOptions = errors.New("select has no options to display")
	ErrUserAborted     = errors.New("user aborted")
	ErrNoTitle         = errors.New("title must be set")
	ErrNoValue         = errors.New("value must be set")
)
