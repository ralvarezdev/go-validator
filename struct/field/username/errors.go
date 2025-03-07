package username

import (
	"errors"
)

var (
	ErrMustBeAlphanumeric = errors.New("username must be have alphanumeric characters")
)
