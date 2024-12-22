package service

import (
	"errors"
)

var (
	NilValidatorError   = errors.New("validator cannot be nil")
	NilMessageError     = errors.New("message cannot be nil")
	NilValidationsError = errors.New("validations cannot be nil")
)
