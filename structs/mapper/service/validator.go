package service

import (
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper"
	govalidatormapperparser "github.com/ralvarezdev/go-validator/structs/mapper/parser"
	govalidatormappervalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/structs/mapper/validator"
)

type (
	// Service interface for the validator service
	Service interface {
		ValidateNilFields(
			validations govalidatormappervalidations.Validations,
			request interface{},
			mapper *govalidatormapper.Mapper,
		) error
		ParseValidations(validations govalidatormappervalidations.Validations) (
			interface{},
			error,
		)
		RunAndParseValidations(
			getValidationsFn func(govalidatormappervalidations.Validations) error,
		) (interface{}, error)
	}

	// DefaultService struct
	DefaultService struct {
		parser    govalidatormapperparser.Parser
		validator govalidatormappervalidator.Validator
	}
)

// NewDefaultService creates a new default validator service
func NewDefaultService(
	parser govalidatormapperparser.Parser,
	validator govalidatormappervalidator.Validator,
) (*DefaultService, error) {
	// Check if the parser or the validator is nil
	if parser == nil {
		return nil, govalidatormapperparser.ErrNilParser
	}
	if validator == nil {
		return nil, govalidatormappervalidator.ErrNilValidator
	}

	return &DefaultService{
		parser:    parser,
		validator: validator,
	}, nil
}

// ValidateNilFields validates the nil fields
func (d *DefaultService) ValidateNilFields(
	validations govalidatormappervalidations.Validations,
	request interface{},
	mapper *govalidatormapper.Mapper,
) error {
	return d.validator.ValidateNilFields(
		validations,
		request,
		mapper,
	)
}

// ParseValidations parses the validations
func (d *DefaultService) ParseValidations(
	validations govalidatormappervalidations.Validations,
) (interface{}, error) {
	// Check if there are any failed validations
	if !validations.HasFailed() {
		return nil, nil
	}

	// Get the parsed validations from the validations
	parsedValidations, err := d.parser.ParseValidations(validations)
	if err != nil {
		return nil, err
	}
	return parsedValidations, nil
}

// RunAndParseValidations runs and parses the validations
func (d *DefaultService) RunAndParseValidations(
	getValidationsFn func(govalidatormappervalidations.Validations) error,
) (interface{}, error) {
	// Initialize struct fields validations
	validations := govalidatormappervalidations.NewDefaultValidations()

	// Get the validations
	err := getValidationsFn(validations)
	if err != nil {
		return nil, err
	}

	// Parse the validations
	return d.ParseValidations(validations)
}
