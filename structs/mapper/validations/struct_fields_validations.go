package validations

import (
	"strings"
)

type (
	// Fields is a struct that holds the fields names for the generated validations of a struct
	Fields struct {
		Validations string
		Errors      string
	}

	// Validations interface is an interface for struct fields validations
	Validations interface {
		HasFailed() bool
		AddFailedFieldValidationError(
			validationName string,
			validationError error,
		)
		SetNestedValidations(
			validationName string,
			nestedValidations Validations,
		)
		GetFailedValidations() *map[string][]error
		GetNestedValidations() *map[string]*Validations
	}

	// DefaultValidations is a struct that holds the error messages for failed validations of a struct
	DefaultValidations struct {
		FailedValidations *map[string][]error
		NestedValidations *map[string]*DefaultValidations
	}

	// Generator is an interface for generating struct fields validations
	Generator interface {
		GetLevelPadding(level int) string
		GenerateFailedValidationsMessage(
			validations *Validations,
			level int,
		) (*string, error)
		Generate(validations *Validations) (*string, error)
	}

	// DefaultGenerator is a struct that holds the default fields names for the generated validations of a struct
	DefaultGenerator struct {
		Fields *Fields
	}
)

// NewFields creates a new Fields struct
func NewFields(validations, errors string) *Fields {
	return &Fields{
		Validations: validations,
		Errors:      errors,
	}
}

// DefaultFields is the default fields names for the generated validations of a struct
var DefaultFields = NewFields("$validations", "$errors")

// NewDefaultValidations creates a new DefaultValidations struct
func NewDefaultValidations() *DefaultValidations {
	// Initialize the struct fields validations
	failedFieldsValidations := make(map[string][]error)
	nestedFieldsValidations := make(map[string]*DefaultValidations)

	return &DefaultValidations{
		FailedValidations: &failedFieldsValidations,
		NestedValidations: &nestedFieldsValidations,
	}
}

// NewDefaultGenerator creates a new DefaultGenerator struct
func NewDefaultGenerator(fields *Fields) *DefaultGenerator {
	// Check if the fields are nil
	if fields == nil {
		fields = DefaultFields
	}

	return &DefaultGenerator{
		Fields: fields,
	}
}

// HasFailed returns true if there are failed validations
func (d *DefaultValidations) HasFailed() bool {
	// Check if there's a nested failed validations
	if d.NestedValidations != nil {
		for _, nestedValidation := range *d.NestedValidations {
			if nestedValidation.HasFailed() {
				return true
			}
		}
	}

	// Check if there are failed fields validations
	if d.FailedValidations == nil {
		return false
	}
	return len(*d.FailedValidations) > 0
}

// AddFailedFieldValidationError adds a failed field validation error to the struct
func (d *DefaultValidations) AddFailedFieldValidationError(
	validationName string,
	validationError error,
) {
	// Check if the field name is already in the map
	failedFieldsValidations := *d.FailedValidations
	if _, ok := failedFieldsValidations[validationName]; !ok {
		failedFieldsValidations[validationName] = []error{validationError}
	} else {
		// Append the validation error to the field name
		failedFieldsValidations[validationName] = append(
			failedFieldsValidations[validationName],
			validationError,
		)
	}
}

// SetNestedValidations sets the nested struct fields validations to the struct
func (d *DefaultValidations) SetNestedValidations(
	validationName string,
	nestedValidations *DefaultValidations,
) {
	(*d.NestedValidations)[validationName] = nestedValidations
}

// GetFailedValidations returns the failed fields validations
func (d *DefaultValidations) GetFailedValidations() *map[string][]error {
	return d.FailedValidations
}

// GetNestedValidations returns the nested struct fields validations
func (d *DefaultValidations) GetNestedValidations() *map[string]*DefaultValidations {
	return d.NestedValidations
}

// GetLevelPadding returns the padding for the level
func (d *DefaultGenerator) GetLevelPadding(level int) string {
	var padding strings.Builder
	for i := 0; i < level; i++ {
		padding.WriteString("\t")
	}
	return padding.String()
}

// GenerateFailedValidationsMessage returns a formatted error message for DefaultGenerator
func (d *DefaultGenerator) GenerateFailedValidationsMessage(
	validations *Validations,
	level int,
) (*string, error) {
	// Check if the validations are nil
	if validations == nil {
		return nil, NilValidationsError
	}

	// Check if there are failed validations
	validationsStruct := *validations
	if !validationsStruct.HasFailed() {
		return nil, nil
	}

	// Get the padding for initial level, the fields, their properties and errors
	basePadding := d.GetLevelPadding(level)
	fieldPadding := d.GetLevelPadding(level + 1)
	fieldPropertiesPadding := d.GetLevelPadding(level + 2)
	fieldErrorsPadding := d.GetLevelPadding(level + 3)

	// Create the message
	var message strings.Builder
	message.WriteString(basePadding)
	message.WriteString(d.Fields.Validations)
	message.WriteString(": {\n")

	// Get the number of nested fields validations
	iteratedFields := make(map[string]bool)
	fieldsValidationsNumber := 0
	nestedFieldsValidationsNumber := 0
	iteratedFieldsValidationsNumber := 0
	iteratedNestedFieldsValidationsNumber := 0

	if validationsStruct.GetFailedValidations() != nil {
		fieldsValidationsNumber = len(*validationsStruct.GetFailedValidations())
	}
	if validationsStruct.GetNestedValidations() != nil {
		nestedFieldsValidationsNumber = len(*validationsStruct.GetNestedValidations())
	}

	// Iterate over all fields and their errors
	var nestedValidations *map[string]*Validations
	for field, fieldErrors := range *validationsStruct.GetFailedValidations() {
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
		message.WriteString(d.Fields.Errors)
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
		var nestedFieldValidations *Validations
		ok := false
		if nestedFieldsValidationsNumber > 0 {
			nestedValidations = validationsStruct.GetNestedValidations()
			nestedFieldValidations, ok = (*nestedValidations)[field]
		}

		// Add comma if not it does not have nested fields
		message.WriteString(fieldPropertiesPadding)
		if !ok || !(*nestedFieldValidations).HasFailed() {
			if ok {
				iteratedNestedFieldsValidationsNumber++
			}

			message.WriteString("]\n")
		} else {
			iteratedNestedFieldsValidationsNumber++
			nestedFieldValidationsMessage, err := d.GenerateFailedValidationsMessage(
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
		for field, nestedFieldValidations := range *validationsStruct.GetNestedValidations() {
			if _, ok := iteratedFields[field]; ok {
				continue
			}

			iteratedNestedFieldsValidationsNumber++
			nestedFieldValidationsMessage, err := d.GenerateFailedValidationsMessage(
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

// Generate returns a pointer to the failed validations message
func (d *DefaultGenerator) Generate(validations *Validations) (
	message *string,
	err error,
) {
	// Return the failed validations message
	message, err = d.GenerateFailedValidationsMessage(validations, 0)
	if err != nil {
		return nil, err
	}

	// Replace all escaped characters
	if message != nil {
		*message = strings.ReplaceAll(*message, "\\t", "\t")
		*message = strings.ReplaceAll(*message, "\\n", "\n")
		return message, nil
	}
	return nil, nil
}
