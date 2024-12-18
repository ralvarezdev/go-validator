package validations

import (
	"strings"
)

// Constants for the struct fields validations
const (
	Validations = "_validations"
	Errors      = "_errors"
)

// MapperValidations is a struct that holds the error messages for failed validations of a struct
type MapperValidations struct {
	FailedMapperValidations  *map[string][]error
	NestedMappersValidations *map[string]*MapperValidations
}

// NewMapperValidations creates a new MapperValidations struct
func NewMapperValidations() *MapperValidations {
	// Initialize the struct fields validations
	failedFieldsValidations := make(map[string][]error)
	nestedFieldsValidations := make(map[string]*MapperValidations)

	return &MapperValidations{
		FailedMapperValidations:  &failedFieldsValidations,
		NestedMappersValidations: &nestedFieldsValidations,
	}
}

// HasFailed returns true if there are failed validations
func (s *MapperValidations) HasFailed() bool {
	// Check if there's a nested failed validations
	if s.NestedMappersValidations != nil {
		for _, nestedStructFieldsValidations := range *s.NestedMappersValidations {
			if nestedStructFieldsValidations.HasFailed() {
				return true
			}
		}
	}

	// Check if there are failed fields validations
	if s.FailedMapperValidations == nil {
		return false
	}
	return len(*s.FailedMapperValidations) > 0
}

// AddFailedFieldValidationError adds a failed field validation error to the struct
func (s *MapperValidations) AddFailedFieldValidationError(validationName string, validationError error) {
	// Check if the field name is already in the map
	failedFieldsValidations := *s.FailedMapperValidations
	if _, ok := failedFieldsValidations[validationName]; !ok {
		failedFieldsValidations[validationName] = []error{validationError}
	} else {
		// Append the validation error to the field name
		failedFieldsValidations[validationName] = append(failedFieldsValidations[validationName], validationError)
	}
}

// SetNestedMapperValidations sets the nested struct fields validations to the struct
func (s *MapperValidations) SetNestedMapperValidations(
	validationName string,
	nestedStructFieldsValidations *MapperValidations,
) {
	(*s.NestedMappersValidations)[validationName] = nestedStructFieldsValidations
}

// GetLevelPadding returns the padding for the level
func (s *MapperValidations) GetLevelPadding(level int) string {
	var padding strings.Builder
	for i := 0; i < level; i++ {
		padding.WriteString("\t")
	}
	return padding.String()
}

// FailedValidationsMessage returns a formatted error message for MapperValidations
func (s *MapperValidations) FailedValidationsMessage(level int) *string {
	// Check if there are failed validations
	if !s.HasFailed() {
		return nil
	}

	// Get the padding for initial level, the fields, their properties and errors
	basePadding := s.GetLevelPadding(level)
	fieldPadding := s.GetLevelPadding(level + 1)
	fieldPropertiesPadding := s.GetLevelPadding(level + 2)
	fieldErrorsPadding := s.GetLevelPadding(level + 3)

	// Create the message
	var message strings.Builder
	message.WriteString(basePadding)
	message.WriteString(Validations)
	message.WriteString(": {\n")

	// Get the number of nested fields validations
	iteratedFields := make(map[string]bool)
	fieldsValidationsNumber := 0
	nestedFieldsValidationsNumber := 0
	iteratedFieldsValidationsNumber := 0
	iteratedNestedFieldsValidationsNumber := 0

	if s.FailedMapperValidations != nil {
		fieldsValidationsNumber = len(*s.FailedMapperValidations)
	}
	if s.NestedMappersValidations != nil {
		nestedFieldsValidationsNumber = len(*s.NestedMappersValidations)
	}

	// Iterate over all fields and their errors
	for field, fieldErrors := range *s.FailedMapperValidations {
		iteratedFieldsValidationsNumber++

		// Check if the field has no errors
		if len(fieldErrors) == 0 {
			continue
		}

		// Add field name
		message.WriteString(fieldPadding)
		message.WriteString(field)
		message.WriteString(": {\n")

		// Add field properties flag
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
		var nestedFieldValidations *MapperValidations
		ok := false
		if nestedFieldsValidationsNumber > 0 {
			nestedFieldValidations, ok = (*s.NestedMappersValidations)[field]
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
			nestedFieldValidationsMessage := nestedFieldValidations.FailedValidationsMessage(level + 1)

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
		for field, nestedFieldValidations := range *s.NestedMappersValidations {
			if _, ok := iteratedFields[field]; ok {
				continue
			}

			iteratedNestedFieldsValidationsNumber++
			nestedFieldValidationsMessage := nestedFieldValidations.FailedValidationsMessage(level + 1)

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

	return &messageString
}

// StringPtr returns a pointer to the failed validations message
func (s *MapperValidations) StringPtr() *string {
	// Return the failed validations message
	message := s.FailedValidationsMessage(0)

	// Replace all escaped characters
	if message != nil {
		*message = strings.ReplaceAll(*message, "\\t", "\t")
		*message = strings.ReplaceAll(*message, "\\n", "\n")
		return message
	}
	return nil
}
