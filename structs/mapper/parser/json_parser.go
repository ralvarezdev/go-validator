package parser

import (
	govalidatorvalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
)

type (
	// JSONParsedValidations is a struct that holds the JSON parsed validations
	JSONParsedValidations struct {
		Fields *map[string]*JSONParsedValidations `json:"$fields,omitempty"`
		Errors *[]string                          `json:"$errors,omitempty"`
	}

	// JSONParser is a struct that holds the JSON parser
	JSONParser struct{}
)

// NewJSONParsedValidations creates a new JSONParsedValidations struct
func NewJSONParsedValidations() *JSONParsedValidations {
	return &JSONParsedValidations{}
}

// AddFieldParsedValidations adds a field parsed validations to the JSON parsed validations
func (j *JSONParsedValidations) AddFieldParsedValidations(
	field string,
	fieldParsedValidations *JSONParsedValidations,
) {
	if j.Fields == nil {
		j.Fields = &map[string]*JSONParsedValidations{}
	}
	(*j.Fields)[field] = fieldParsedValidations
}

// GetFieldParsedValidations returns the field parsed validations from the JSON parsed validations
func (j *JSONParsedValidations) GetFieldParsedValidations(field string) *JSONParsedValidations {
	if j.Fields == nil {
		return nil
	}
	return (*j.Fields)[field]
}

// AddErrors adds errors to the JSON parsed validations
func (j *JSONParsedValidations) AddErrors(errors *[]error) {
	if j.Errors == nil {
		j.Errors = &[]string{}
	}

	// Iterate over all errors and add them to the JSON parsed validations
	for _, err := range *errors {
		*j.Errors = append(*j.Errors, err.Error())
	}
}

// NewJSONParser creates a new JSONParser struct
func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

// GenerateJSONParsedValidations returns a
func (j *JSONParser) GenerateJSONParsedValidations(
	validations govalidatorvalidations.Validations,
	jsonParsedValidations *JSONParsedValidations,
) error {
	// Check if the validations or JSON parsed validations are nil
	if validations == nil {
		return govalidatorvalidations.ErrNilValidations
	}
	if jsonParsedValidations == nil {
		return ErrNilJSONParsedValidations
	}

	// Check if there are failed validations
	if !validations.HasFailed() {
		return nil
	}

	// Get the fields validations
	failedFieldsValidations := *validations.GetFailedFieldsValidations()
	nestedFieldsValidations := *validations.GetNestedFieldsValidations()

	// Iterate over all fields and their errors
	var nestedJSONParsedValidations *JSONParsedValidations
	for field, fieldErrors := range failedFieldsValidations {
		// Check if the field has no errors
		if len(fieldErrors) == 0 {
			continue
		}

		// Initialize the JSON parsed validations
		nestedJSONParsedValidations = NewJSONParsedValidations()

		// Set the fields errors if there are any
		nestedJSONParsedValidations.AddErrors(&fieldErrors)
		jsonParsedValidations.AddFieldParsedValidations(
			field,
			nestedJSONParsedValidations,
		)
	}

	// Iterate over all nested fields validations
	for field, nestedFieldValidations := range nestedFieldsValidations {
		// Check if the nested field validations are nil
		if nestedFieldValidations == nil {
			continue
		}

		// Check if the given field is already in the JSON parsed validations
		nestedJSONParsedValidations = jsonParsedValidations.GetFieldParsedValidations(field)
		if nestedJSONParsedValidations == nil {
			nestedJSONParsedValidations = NewJSONParsedValidations()
			jsonParsedValidations.AddFieldParsedValidations(
				field,
				nestedJSONParsedValidations,
			)
		}

		// Generate the nested JSON parsed validations
		err := j.GenerateJSONParsedValidations(
			nestedFieldValidations,
			nestedJSONParsedValidations,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseValidations parses the validations into JSON
func (j *JSONParser) ParseValidations(validations govalidatorvalidations.Validations) (
	interface{},
	error,
) {
	// Initialize the JSON validations
	jsonParsedValidations := NewJSONParsedValidations()

	// Generate the JSON parsed validations
	err := j.GenerateJSONParsedValidations(validations, jsonParsedValidations)
	if err != nil {
		return nil, err
	}
	return jsonParsedValidations, nil
}
