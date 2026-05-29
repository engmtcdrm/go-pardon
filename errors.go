package pardon

import (
	"errors"
)

var (
	ErrNoPrompts       = errors.New("form requires at least one prompt")
	ErrNoSelectOptions = errors.New("select prompt requires at least one option")
	ErrNoTitle         = errors.New("prompt requires a title")
	ErrNoValue         = errors.New("value must be set")
	ErrUserAborted     = errors.New("user aborted")
)
