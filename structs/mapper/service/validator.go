package service

import (
	goflagmode "github.com/ralvarezdev/go-flags/mode"
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper"
	govalidatormapperparser "github.com/ralvarezdev/go-validator/structs/mapper/parser"
	govalidatormappervalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
	govalidatormappervalidator "github.com/ralvarezdev/go-validator/structs/mapper/validator"
)

type (
	// Service interface for the validator service
	Service interface {
		ModeFlag() *goflagmode.Flag
		ValidateNilFields(
			request interface{},
			mapper *govalidatormapper.Mapper,
		) (
			govalidatormappervalidations.Validations,
			error,
		)
		ParseValidations(validations govalidatormappervalidations.Validations) (
			interface{},
			error,
		)
		RunAndParseValidations(
			getValidationsFn func() (
				govalidatormappervalidations.Validations,
				error,
			),
		) (interface{}, error)
	}

	// DefaultService struct
	DefaultService struct {
		mode      *goflagmode.Flag
		parser    govalidatormapperparser.Parser
		validator govalidatormappervalidator.Validator
	}
)

// NewDefaultService creates a new default validator service
func NewDefaultService(
	parser govalidatormapperparser.Parser,
	validator govalidatormappervalidator.Validator,
	mode *goflagmode.Flag,
) (*DefaultService, error) {
	// Check if the parser or the validator is nil
	if parser == nil {
		return nil, govalidatormapperparser.ErrNilParser
	}
	if validator == nil {
		return nil, govalidatormappervalidations.ErrNilValidator
	}

	return &DefaultService{
		mode:      mode,
		parser:    parser,
		validator: validator,
	}, nil
}

// ModeFlag returns the mode flag
func (d *DefaultService) ModeFlag() *goflagmode.Flag {
	return d.mode
}

// ValidateNilFields validates the nil fields
func (d *DefaultService) ValidateNilFields(
	request interface{},
	mapper *govalidatormapper.Mapper,
) (govalidatormappervalidations.Validations, error) {
	return d.validator.ValidateNilFields(
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
	getValidationsFn func() (govalidatormappervalidations.Validations, error),
) (interface{}, error) {
	// Get the validations
	validations, err := getValidationsFn()
	if err != nil {
		return nil, err
	}

	// Parse the validations
	return d.ParseValidations(validations)
}
