package parser

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// AssertParsedValidations asserts parsed validations to map[string]interface{}
//
// Parameters:
//
//   - validations: The parsed validations to assert
//
// Returns:
//
//   - map[string]interface{}: The asserted parsed validations
//   - bool: Whether the assertion was successful
func AssertParsedValidations(validations interface{}) (
	map[string]interface{},
	bool,
) {
	// Check if there are no validations
	if validations == nil {
		return nil, false
	}

	// Assert the parsed validations to map[string]interface{}
	assertedValidations, ok := validations.(map[string]interface{})
	return assertedValidations, ok
}

// IsNestedParsedValidations checks if the parsed validations are nested
//
// Parameters:
//
//   - validations: The parsed validations to check
//   - field: The field to check for nested validations
//
// Returns:
//
//   - bool: Whether the parsed validations are nested
func IsNestedParsedValidations(
	validations map[string]interface{},
	field string,
) bool {
	// Check if there are no validations
	if validations == nil {
		return false
	}

	// Check if the field exists in the parsed validations
	nestedValidations, ok := validations[field]
	if !ok {
		return false
	}

	// Assert the nested validations to map[string]interface{}
	_, ok = nestedValidations.(map[string]interface{})
	return ok
}

// ParseValidationsToBadRequest converts parsed validations to a BadRequest
//
// Parameters:
//
//   - validations: The parsed validations to convert
//
// Returns:
//
//   - *errdetails.BadRequest: The converted validations to BadRequest
func ParseValidationsToBadRequest(validations interface{}) *errdetails.BadRequest {
	// Check if there are no validations
	if validations == nil {
		return nil
	}

	// Get the field

	br := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       field,
				Description: description,
			},
		},
	}

}
