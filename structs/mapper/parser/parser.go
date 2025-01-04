package parser

import (
	govalidatorvalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
)

type (
	// Parser is an interface to parse struct fields validations
	Parser interface {
		ParseValidations(validations govalidatorvalidations.Validations) (
			interface{},
			error,
		)
	}
)
