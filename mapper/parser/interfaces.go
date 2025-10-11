package parser

import (
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/mapper/validation"
)

type (
	// RawParser is an interface to generate parsed validations from struct fields validations
	RawParser interface {
		ParseValidations(
			structValidations *govalidatormappervalidation.StructValidations,
			dest *StructParsedValidations,
		) error
	}

	// EndParser is an interface to parse the root struct validations into a final format
	EndParser interface {
		ParseValidations(structParsedValidations *StructParsedValidations) (
			interface{},
			error,
		)
	}
)
