package parser

import (
	govalidatormappervalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
)

type (
	// Parser is an interface to parse struct fields validations
	Parser interface {
		ParseValidations(rootStructValidations *govalidatormappervalidations.StructValidations) (
			interface{},
			error,
		)
	}
)
