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
		ValidateRequiredFields(
			rootStructValidations *govalidatormappervalidations.StructValidations,
			request interface{},
			mapper *govalidatormapper.Mapper,
		) error
		ParseValidations(rootStructValidations *govalidatormappervalidations.StructValidations) (
			interface{},
			error,
		)
		RunAndParseValidations(
			getValidationsFn func(*govalidatormappervalidations.StructValidations) error,
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

// ValidateRequiredFields validates the required fields
func (d *DefaultService) ValidateRequiredFields(
	rootStructValidations *govalidatormappervalidations.StructValidations,
	request interface{},
	mapper *govalidatormapper.Mapper,
) error {
	return d.validator.ValidateRequiredFields(
		rootStructValidations,
		request,
		mapper,
	)
}

// ParseValidations parses the validations
func (d *DefaultService) ParseValidations(
	rootStructValidations *govalidatormappervalidations.StructValidations,
) (interface{}, error) {
	// Check if there are any failed validations
	if !rootStructValidations.HasFailed() {
		return nil, nil
	}

	// Get the parsed validations from the validations
	parsedValidations, err := d.parser.ParseValidations(rootStructValidations)
	if err != nil {
		return nil, err
	}
	return parsedValidations, nil
}

// RunAndParseValidations runs and parses the validations
func (d *DefaultService) RunAndParseValidations(
	getValidationsFn func(*govalidatormappervalidations.StructValidations) error,
) (interface{}, error) {
	// Initialize struct fields validations
	rootStructValidations := govalidatormappervalidations.NewStructValidations()

	// Get the validations
	err := getValidationsFn(rootStructValidations)
	if err != nil {
		return nil, err
	}

	// Parse the validations
	return d.ParseValidations(rootStructValidations)
}
