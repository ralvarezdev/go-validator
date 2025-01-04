package validations

type (
	// Validations interface is an interface for struct fields validations
	Validations interface {
		HasFailed() bool
		AddFailedFieldValidationError(
			fieldName string,
			validationError error,
		)
		SetNestedFieldsValidations(
			fieldName string,
			nestedValidations Validations,
		)
		GetFailedFieldsValidations() *map[string][]error
		GetNestedFieldsValidations() *map[string]Validations
	}

	// DefaultValidations is a struct that holds the error messages for failed validations of a struct
	DefaultValidations struct {
		FailedFieldsValidations *map[string][]error
		NestedFieldsValidations *map[string]Validations
	}
)

// NewDefaultValidations creates a new DefaultValidations struct
func NewDefaultValidations() *DefaultValidations {
	// Initialize the struct fields validations
	failedFieldsValidations := make(map[string][]error)
	nestedFieldsValidations := make(map[string]Validations)

	return &DefaultValidations{
		FailedFieldsValidations: &failedFieldsValidations,
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
	if d.FailedFieldsValidations == nil {
		return false
	}
	return len(*d.FailedFieldsValidations) > 0
}

// AddFailedFieldValidationError adds a failed field validation error to the struct
func (d *DefaultValidations) AddFailedFieldValidationError(
	fieldName string,
	validationError error,
) {
	// Check if the field name is already in the map
	failedFieldsValidations := *d.FailedFieldsValidations
	if _, ok := failedFieldsValidations[fieldName]; !ok {
		failedFieldsValidations[fieldName] = []error{validationError}
	} else {
		// Append the validation error to the field name
		failedFieldsValidations[fieldName] = append(
			failedFieldsValidations[fieldName],
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

// GetFailedFieldsValidations returns the failed fields validations
func (d *DefaultValidations) GetFailedFieldsValidations() *map[string][]error {
	return d.FailedFieldsValidations
}

// GetNestedFieldsValidations returns the nested struct fields validations
func (d *DefaultValidations) GetNestedFieldsValidations() *map[string]Validations {
	return d.NestedFieldsValidations
}
