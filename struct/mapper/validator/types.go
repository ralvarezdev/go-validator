package validator

import (
	govalidatormappervalidation "github.com/ralvarezdev/go-validator/struct/mapper/validation"
)

type (
	// AuxiliaryValidatorFn is the type for the auxiliary validator function
	AuxiliaryValidatorFn func(
		toValidate interface{},
		validations *govalidatormappervalidation.StructValidations,
	)

	// ValidateFn is the type for the validate function
	ValidateFn func(toValidate interface{}) (interface{}, error)
)
