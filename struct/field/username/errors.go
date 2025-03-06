package username

import (
	"errors"
)

var (
	ErrFoundWhitespaces = errors.New("username cannot contain whitespaces")
)
