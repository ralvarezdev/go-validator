package grpc

import (
	"fmt"

	govalidatormapperparser "github.com/ralvarezdev/go-validator/struct/mapper/parser"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type (
	// ErrorDetails is the struct for the error details wrapper
	ErrorDetails struct {
		fieldViolations []*errdetails.BadRequest_FieldViolation
	}

	// DefaultEndParser is the default implementation of the EndParser interface
	DefaultEndParser struct{}
)

// NewErrorDetails adds the root struct parsed validations to the error details
//
// Parameters:
//
//   - structParsedValidations: The root struct parsed validations to add
//   - parentFieldName: The parent field name to prefix to the field names
//   - fieldsViolations: The existing field violations to add to the error details
//
// Returns:
//
//   - error: An error if the root struct parsed validations are nil or if the fields are already in the error details
func NewErrorDetails(
	structParsedValidations *govalidatormapperparser.StructParsedValidations,
	parentFieldName *string,
	fieldsViolations []*errdetails.BadRequest_FieldViolation,
) (*ErrorDetails, error) {
	// Check if the root struct parsed validations are nil
	if structParsedValidations == nil {
		return nil, govalidatormapperparser.ErrNilStructParsedValidations
	}

	// Create the error details
	var e *ErrorDetails
	if fieldsViolations != nil {
		e = &ErrorDetails{
			fieldViolations: []*errdetails.BadRequest_FieldViolation{},
		}
	} else {
		e = &ErrorDetails{
			fieldViolations: fieldsViolations,
		}
	}

	// Add the struct parsed validations fields
	fieldsParsedValidations := structParsedValidations.GetFields()
	if fieldsParsedValidations != nil {
		for fieldName, fieldParsedValidations := range fieldsParsedValidations {
			// Prefix the field name with the parent field name if it exists
			if parentFieldName != nil {
				fieldName = fmt.Sprintf("%s.%s", *parentFieldName, fieldName)
			}

			// Add the field parsed validations
			if err := e.AddField(
				fieldName,
				fieldParsedValidations,
			); err != nil {
				return nil, err
			}
		}
	}

	// Add the struct parsed validations nested structs
	nestedStructsParsedValidations := structParsedValidations.GetNestedStructs()
	if nestedStructsParsedValidations != nil {
		for fieldName, nestedStructParsedValidations := range nestedStructsParsedValidations {
			// Prefix the field name with the parent field name if it exists
			if parentFieldName != nil {
				fieldName = fmt.Sprintf("%s.%s", *parentFieldName, fieldName)
			}

			// Add the nested struct parsed validations
			if err := e.AddNestedStruct(
				fieldName,
				nestedStructParsedValidations,
			); err != nil {
				return nil, err
			}
		}
	}

	return e, nil
}

// AddField adds a field parsed validations to the error details
//
// Parameters:
//
//   - fieldName: The name of the field
//   - fieldParsedValidations: The field parsed validations to add
//
// Returns:
//
//   - error: An error if the field name is already in the error details
func (e *ErrorDetails) AddField(
	fieldName string,
	fieldParsedValidations *govalidatormapperparser.FieldParsedValidations,
) error {
	if e == nil {
		return ErrNilErrorDetails
	}

	// Check if the field name is empty or the field parsed validations are nil
	if fieldName == "" || fieldParsedValidations == nil {
		return nil
	}

	// Check if the fields violations are nil
	if e.fieldViolations == nil {
		e.fieldViolations = []*errdetails.BadRequest_FieldViolation{}
	}

	// Add the field parsed validations to the error details
	for _, err := range fieldParsedValidations.GetErrors() {
		e.fieldViolations = append(
			e.fieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       fieldName,
				Description: err,
			},
		)
	}
	return nil
}

// AddNestedStruct adds a nested struct parsed validations to the error details
//
// Parameters:
//
//   - fieldName: The name of the field that holds the nested struct
//   - structParsedValidations: The struct parsed validations to add
//
// Returns:
//
//   - error: An error if the struct name is already in the error details
func (e *ErrorDetails) AddNestedStruct(
	fieldName string,
	structParsedValidations *govalidatormapperparser.StructParsedValidations,
) error {
	if e == nil {
		return ErrNilErrorDetails
	}

	// Check if the struct name is empty or the struct parsed validations are nil
	if structParsedValidations == nil {
		return nil
	}

	// Check if the fields violations are nil
	if e.fieldViolations == nil {
		e.fieldViolations = []*errdetails.BadRequest_FieldViolation{}
	}

	// Get the struct error details
	_, err := NewErrorDetails(
		structParsedValidations,
		&fieldName,
		e.fieldViolations,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetBadRequest gets the BadRequest from the ErrorDetails
//
// Returns:
//
//   - *errdetails.BadRequest: The BadRequest from the ErrorDetails
func (e *ErrorDetails) GetBadRequest() *errdetails.BadRequest {
	if e == nil {
		return nil
	}
	return &errdetails.BadRequest{
		FieldViolations: e.fieldViolations,
	}
}

// NewDefaultEndParser creates a new DefaultEndParser
//
// Returns:
//
//   - DefaultEndParser: The new DefaultEndParser
func NewDefaultEndParser() DefaultEndParser {
	return DefaultEndParser{}
}

// ParseValidations parses the validations into a BadRequest
//
// Parameters:
//
//   - structValidations: The root struct validations
//
// Returns:
//
//   - interface{}: The parsed validations
//   - error: An error if the root struct validations are nil or if there was an error generating the BadRequest
func (d DefaultEndParser) ParseValidations(structParsedValidations *govalidatormapperparser.StructParsedValidations) (
	interface{},
	error,
) {
	// Check if the root struct parsed validations are nil
	if structParsedValidations == nil {
		return nil, govalidatormapperparser.ErrNilStructParsedValidations
	}

	// Convert the parsed validations to a BadRequest
	errorDetails, err := NewErrorDetails(structParsedValidations, nil, nil)
	if err != nil {
		return nil, err
	}
	return errorDetails.GetBadRequest(), nil
}
