package json

import (
	govalidatormappervalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
	govalidatormapper "github.com/ralvarezdev/go-validator/structs/mapper/validator"
)

type (
	// ParsedValidations interface for the JSON parsed validations
	ParsedValidations interface {
		AddFieldParsedValidations(
			field string,
			fieldParsedValidations ParsedValidations,
		)
		GetFieldParsedValidations(field string) ParsedValidations
		AddErrors(errors *[]error)
	}

	// DefaultParsedValidations is the struct for the default JSON parsed validations
	DefaultParsedValidations struct {
		Fields *map[string]ParsedValidations `json:"$fields,omitempty"`
		Errors *[]string                     `json:"$errors,omitempty"`
	}

	// Parser is a struct that holds the JSON parser
	Parser struct {
		newParsedValidationsFn func() ParsedValidations
	}
)

// NewDefaultParsedValidations creates a new DefaultParsedValidations struct
func NewDefaultParsedValidations() ParsedValidations {
	return &DefaultParsedValidations{}
}

// AddFieldParsedValidations adds a field parsed validations to the JSON parsed validations
func (d *DefaultParsedValidations) AddFieldParsedValidations(
	field string,
	fieldParsedValidations ParsedValidations,
) {
	if d.Fields == nil {
		d.Fields = &map[string]ParsedValidations{}
	}
	(*d.Fields)[field] = fieldParsedValidations
}

// GetFieldParsedValidations returns the field parsed validations from the JSON parsed validations
func (d *DefaultParsedValidations) GetFieldParsedValidations(field string) ParsedValidations {
	if d.Fields == nil {
		return nil
	}
	return (*d.Fields)[field]
}

// AddErrors adds errors to the JSON parsed validations
func (d *DefaultParsedValidations) AddErrors(errors *[]error) {
	if d.Errors == nil {
		d.Errors = &[]string{}
	}

	// Iterate over all errors and add them to the JSON parsed validations
	for _, err := range *errors {
		*d.Errors = append(*d.Errors, err.Error())
	}
}

// NewParser creates a new Parser struct
func NewParser(newParsedValidationsFn func() ParsedValidations) *Parser {
	return &Parser{newParsedValidationsFn: newParsedValidationsFn}
}

// GenerateParsedValidations returns a
func (p *Parser) GenerateParsedValidations(
	validations govalidatormappervalidations.Validations,
	parsedValidations ParsedValidations,
) error {
	// Check if the validations or parsed validations are nil
	if validations == nil {
		return govalidatormapper.ErrNilValidations
	}
	if parsedValidations == nil {
		return ErrNilParsedValidations
	}

	// Check if there are failed validations
	if !validations.HasFailed() {
		return nil
	}

	// Get the fields validations
	fieldsValidations := *validations.GetFieldsValidations()
	nestedFieldsValidations := *validations.GetNestedFieldsValidations()

	// Iterate over all fields and their errors
	var nestedParsedValidations ParsedValidations
	for field, fieldErrors := range fieldsValidations {
		// Check if the field has no errors
		if len(fieldErrors) == 0 {
			continue
		}

		// Initialize the JSON parsed validations
		nestedParsedValidations = p.newParsedValidationsFn()

		// Set the fields errors if there are any
		nestedParsedValidations.AddErrors(&fieldErrors)
		parsedValidations.AddFieldParsedValidations(
			field,
			nestedParsedValidations,
		)
	}

	// Iterate over all nested fields validations
	for field, nestedFieldValidations := range nestedFieldsValidations {
		// Check if the nested field validations are nil
		if nestedFieldValidations == nil {
			continue
		}

		// Check if the given field is already in the JSON parsed validations
		nestedParsedValidations = parsedValidations.GetFieldParsedValidations(field)
		if nestedParsedValidations == nil {
			nestedParsedValidations = p.newParsedValidationsFn()
			parsedValidations.AddFieldParsedValidations(
				field,
				nestedParsedValidations,
			)
		}

		// Generate the nested JSON parsed validations
		err := p.GenerateParsedValidations(
			nestedFieldValidations,
			nestedParsedValidations,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseValidations parses the validations into JSON
func (p *Parser) ParseValidations(validations govalidatormappervalidations.Validations) (
	interface{},
	error,
) {
	// Initialize the parsed validations
	parsedValidations := p.newParsedValidationsFn()

	// Generate the JSON parsed validations
	err := p.GenerateParsedValidations(
		validations,
		parsedValidations,
	)
	if err != nil {
		return nil, err
	}
	return parsedValidations, nil
}
