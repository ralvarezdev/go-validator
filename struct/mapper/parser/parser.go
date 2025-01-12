package parser

import (
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/struct/mapper/validation"
)

type (
	// Parser is an interface to parse struct fields validations
	Parser interface {
		ParseValidations(rootStructValidations *govalidatormappervalidation.StructValidations) (
			interface{},
			error,
		)
	}
)
