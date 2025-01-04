package validations

type (
	// Validations interface is an interface for struct fields validations
	Validations interface {
		HasFailed() bool
		AddFieldValidationError(
			fieldName string,
			validationError error,
		)
		SetNestedFieldsValidations(
			fieldName string,
			nestedValidations Validations,
		)
		GetFieldsValidations() *map[string][]error
		GetNestedFieldsValidations() *map[string]Validations
	}

	// DefaultValidations is a struct that holds the error messages for failed validations of a struct
	DefaultValidations struct {
		FieldsValidations       *map[string][]error
		NestedFieldsValidations *map[string]Validations
	}
)

// NewDefaultValidations creates a new DefaultValidations struct
func NewDefaultValidations() Validations {
	// Initialize the struct fields validations
	fieldsValidations := make(map[string][]error)
	nestedFieldsValidations := make(map[string]Validations)

	return &DefaultValidations{
		FieldsValidations:       &fieldsValidations,
		NestedFieldsValidations: &nestedFieldsValidations,
	}
}

// HasFailed returns true if there are failed validations
func (d *DefaultValidations) HasFailed() bool {
	// Check if there's a nested failed validations
	if d.NestedFieldsValidations != nil {
		for _, nestedValidation := range *d.NestedFieldsValidations {
			if nestedValidation.HasFailed() {
				return true
			}
		}
	}

	// Check if there are failed fields validations
	if d.FieldsValidations == nil {
		return false
	}
	return len(*d.FieldsValidations) > 0
}

// AddFieldValidationError adds a field validation error to the struct
func (d *DefaultValidations) AddFieldValidationError(
	fieldName string,
	validationError error,
) {
	// Check if the field name is already in the map
	fieldsValidations := *d.FieldsValidations
	if _, ok := fieldsValidations[fieldName]; !ok {
		fieldsValidations[fieldName] = []error{validationError}
	} else {
		// Append the validation error to the field name
		fieldsValidations[fieldName] = append(
			fieldsValidations[fieldName],
			validationError,
		)
	}
}

// SetNestedFieldsValidations sets the nested struct fields validations to the struct
func (d *DefaultValidations) SetNestedFieldsValidations(
	fieldName string,
	nestedValidations Validations,
) {
	(*d.NestedFieldsValidations)[fieldName] = nestedValidations
}

// GetFieldsValidations returns the fields validations errors
func (d *DefaultValidations) GetFieldsValidations() *map[string][]error {
	return d.FieldsValidations
}

// GetNestedFieldsValidations returns the nested struct fields validations
func (d *DefaultValidations) GetNestedFieldsValidations() *map[string]Validations {
	return d.NestedFieldsValidations
}
