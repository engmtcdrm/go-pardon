package gocliselect

import "errors"

var (
	ErrNoSelectItems = errors.New("select has no items to display")
	ErrUserAborted   = errors.New("user aborted")
	ErrNoTitle       = errors.New("title must be set")
	ErrNoValue       = errors.New("value must be set")
)
