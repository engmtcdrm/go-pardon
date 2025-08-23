package pardon

import "errors"

var (
	ErrUserAborted     = errors.New("user aborted")
	ErrNoTitle         = errors.New("prompt requires a title")
	ErrNoSelectOptions = errors.New("select prompt requires at least one option")
	ErrNoValue         = errors.New("value must be set")
)
