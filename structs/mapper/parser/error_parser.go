package parser

import (
	govalidatorvalidations "github.com/ralvarezdev/go-validator/structs/mapper/validations"
	"strings"
)

var (
	// Fields is the fields flag
	Fields = "$fields"

	// Errors is the errors flag
	Errors = "$errors"
)

type (
	// ErrorParser is the struct for the error parser
	ErrorParser struct{}
)

// NewErrorParser creates a new ErrorParser struct
func NewErrorParser() *ErrorParser {
	return &ErrorParser{}
}

// GetLevelPadding returns the padding for the level
func (e *ErrorParser) GetLevelPadding(level int) string {
	var padding strings.Builder
	for i := 0; i < level; i++ {
		padding.WriteString("\t")
	}
	return padding.String()
}

// GenerateFailedValidationsMessage returns a formatted error message for ErrorParser
func (e *ErrorParser) GenerateFailedValidationsMessage(
	validations govalidatorvalidations.Validations,
	level int,
) (*string, error) {
	// Check if the validations are nil
	if validations == nil {
		return nil, govalidatorvalidations.ErrNilValidations
	}

	// Check if there are failed validations
	if !validations.HasFailed() {
		return nil, nil
	}

	// Get the padding for initial level, the fields, their properties and errors
	basePadding := e.GetLevelPadding(level)
	fieldPadding := e.GetLevelPadding(level + 1)
	fieldPropertiesPadding := e.GetLevelPadding(level + 2)
	fieldErrorsPadding := e.GetLevelPadding(level + 3)

	// Create the message and add the fields flag
	var message strings.Builder
	message.WriteString(basePadding)
	message.WriteString(Fields)
	message.WriteString(": {\n")

	// Get the number of nested fields validations
	iteratedFields := make(map[string]bool)
	fieldsValidationsNumber := 0
	nestedFieldsValidationsNumber := 0
	iteratedFieldsValidationsNumber := 0
	iteratedNestedFieldsValidationsNumber := 0

	if validations.GetFailedFieldsValidations() != nil {
		fieldsValidationsNumber = len(*validations.GetFailedFieldsValidations())
	}
	if validations.GetNestedFieldsValidations() != nil {
		nestedFieldsValidationsNumber = len(*validations.GetNestedFieldsValidations())
	}

	// Iterate over all fields and their errors
	var nestedValidations *map[string]govalidatorvalidations.Validations
	for field, fieldErrors := range *validations.GetFailedFieldsValidations() {
		iteratedFieldsValidationsNumber++

		// Check if the field has no errors
		if len(fieldErrors) == 0 {
			continue
		}

		// Add field name
		message.WriteString(fieldPadding)
		message.WriteString(field)
		message.WriteString(": {\n")

		// Add field errors flag
		message.WriteString(fieldPropertiesPadding)
		message.WriteString(Errors)
		message.WriteString(": [\n")

		// Iterate over all errors for the field
		iteratedFields[field] = true
		for index, err := range fieldErrors {
			message.WriteString(fieldErrorsPadding)
			message.WriteString(err.Error())

			// Add comma if not the last error
			if index < len(fieldErrors)-1 {
				message.WriteString(",\n")
			} else {
				message.WriteString("\n")
			}
		}

		// Get the nested fields validations for the field if it has any
		var nestedFieldValidations govalidatorvalidations.Validations
		ok := false
		if nestedFieldsValidationsNumber > 0 {
			nestedValidations = validations.GetNestedFieldsValidations()
			nestedFieldValidations, ok = (*nestedValidations)[field]
		}

		// Add comma if not it does not have nested fields
		message.WriteString(fieldPropertiesPadding)
		if !ok || !nestedFieldValidations.HasFailed() {
			if ok {
				iteratedNestedFieldsValidationsNumber++
			}

			message.WriteString("]\n")
		} else {
			iteratedNestedFieldsValidationsNumber++
			nestedFieldValidationsMessage, err := e.GenerateFailedValidationsMessage(
				nestedFieldValidations,
				level+1,
			)
			if err != nil {
				return nil, err
			}

			// Add nested fields errors
			if nestedFieldValidationsMessage != nil {
				message.WriteString("],\n")
				message.WriteString(*nestedFieldValidationsMessage)
			}
		}

		// Add comma if is not the last field
		message.WriteString(fieldPadding)
		if iteratedFieldsValidationsNumber < fieldsValidationsNumber || iteratedNestedFieldsValidationsNumber < nestedFieldsValidationsNumber {
			message.WriteString("},\n")
		} else {
			message.WriteString("}\n")
		}
	}

	// Iterate over all nested fields validations
	if iteratedNestedFieldsValidationsNumber < nestedFieldsValidationsNumber {
		for field, nestedFieldValidations := range *validations.GetNestedFieldsValidations() {
			if _, ok := iteratedFields[field]; ok {
				continue
			}

			iteratedNestedFieldsValidationsNumber++
			nestedFieldValidationsMessage, err := e.GenerateFailedValidationsMessage(
				nestedFieldValidations,
				level+1,
			)
			if err != nil {
				return nil, err
			}

			// Add field name
			message.WriteString(fieldPadding)
			message.WriteString(field)
			message.WriteString(": {\n")

			// Add nested fields errors
			message.WriteString(fieldPropertiesPadding)
			message.WriteString(*nestedFieldValidationsMessage)

			// Add comma if is not the last field
			message.WriteString(fieldPadding)
			if iteratedNestedFieldsValidationsNumber < nestedFieldsValidationsNumber {
				message.WriteString("},\n")
			} else {
				message.WriteString("}\n")
			}
		}
	}

	// Add closing bracket
	message.WriteString(basePadding)
	message.WriteString("}")

	// Get message string
	messageString := message.String()

	return &messageString, nil
}

// ParseValidations parses the validations and returns the failed validations message
func (e *ErrorParser) ParseValidations(validations govalidatorvalidations.Validations) (
	interface{},
	error,
) {
	// Return the failed validations message
	parsedValidations, err := e.GenerateFailedValidationsMessage(validations, 0)
	if err != nil {
		return nil, err
	}

	// Replace all escaped characters
	if parsedValidations != nil {
		*parsedValidations = strings.ReplaceAll(*parsedValidations, "\\t", "\t")
		*parsedValidations = strings.ReplaceAll(*parsedValidations, "\\n", "\n")
		return parsedValidations, nil
	}
	return nil, nil
}
