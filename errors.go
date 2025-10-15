package pardon

import (
	"errors"

	"github.com/engmtcdrm/go-pardon/internal/tui"
)

var (
	ErrUserAborted     = tui.ErrUserAborted
	ErrNoTitle         = errors.New("prompt requires a title")
	ErrNoSelectOptions = errors.New("select prompt requires at least one option")
	ErrNoValue         = errors.New("value must be set")
)
