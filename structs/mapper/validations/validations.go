package validations

type (
	// StructValidations is a struct that holds the struct validations for the generated validations of a struct
	StructValidations struct {
		fieldsValidations        *map[string]*FieldValidations
		nestedStructsValidations *map[string]*StructValidations
	}

	// FieldValidations is a struct that holds the field validations for the generated validations of a struct
	FieldValidations struct {
		errors *[]error
	}
)

// NewStructValidations creates a new StructValidations struct
func NewStructValidations() *StructValidations {
	return &StructValidations{}
}

// HasFailed returns true if there are failed validations
func (s *StructValidations) HasFailed() bool {
	// Check if there's a nested struct with failed validations
	if s.nestedStructsValidations != nil {
		for _, nestedStructValidation := range *s.nestedStructsValidations {
			if nestedStructValidation.HasFailed() {
				return true
			}
		}
	}
	// Check if there is a field with failed fields validations
	if s.fieldsValidations != nil {
		for _, fieldValidation := range *s.fieldsValidations {
			if fieldValidation.HasFailed() {
				return true
			}
		}
	}
	return false
}

// AddFieldValidations sets the fields validations to the struct
func (s *StructValidations) AddFieldValidations(
	fieldName string,
	fieldValidations *FieldValidations,
) {
	// Check if the field name is empty or the field validations is nil
	if fieldName == "" || fieldValidations == nil {
		return
	}

	// Check if the fields validations are nil
	if s.fieldsValidations == nil {
		fieldsValidations := make(map[string]*FieldValidations)
		s.fieldsValidations = &fieldsValidations
	}

	// Add the field validations to the struct
	(*s.fieldsValidations)[fieldName] = fieldValidations
}

// AddFieldValidationError adds a validation error to the field
func (s *StructValidations) AddFieldValidationError(
	fieldName string,
	validationError error,
) {
	// Check if the field name is empty or the validation error is nil
	if fieldName == "" || validationError == nil {
		return
	}

	// Check if the fields validations are nil
	if s.fieldsValidations == nil {
		fieldsValidations := make(map[string]*FieldValidations)
		s.fieldsValidations = &fieldsValidations
	}

	// Check if the field name is already in the map
	fieldValidations, ok := (*s.fieldsValidations)[fieldName]
	if !ok {
		(*s.fieldsValidations)[fieldName] = NewFieldValidations()
	}

	// Append the validation error to the field name
	fieldValidations.AddValidationError(validationError)
}

// AddNestedStructValidations sets the nested struct fields validations to the struct
func (s *StructValidations) AddNestedStructValidations(
	fieldName string,
	nestedStructValidations *StructValidations,
) {
	// Check if the nested struct validations is nil
	if nestedStructValidations == nil {
		return
	}

	// Check if the nested structs validations are nil
	if s.nestedStructsValidations == nil {
		nestedStructsValidations := make(map[string]*StructValidations)
		s.nestedStructsValidations = &nestedStructsValidations
	}

	// Add the nested struct validations to the struct
	(*s.nestedStructsValidations)[fieldName] = nestedStructValidations
}

// GetFieldsValidations returns the fields validations
func (s *StructValidations) GetFieldsValidations() *map[string]*FieldValidations {
	return s.fieldsValidations
}

// GetNestedStructsValidations returns the nested structs validations
func (s *StructValidations) GetNestedStructsValidations() *map[string]*StructValidations {
	return s.nestedStructsValidations
}

// NewFieldValidations creates a new FieldValidations struct
func NewFieldValidations() *FieldValidations {
	return &FieldValidations{}
}

// HasFailed returns true if there are failed validations
func (f *FieldValidations) HasFailed() bool {
	if f.errors == nil {
		return false
	}
	return len(*f.errors) > 0
}

// AddValidationError adds a validation error to the field
func (f *FieldValidations) AddValidationError(
	validationError error,
) {
	// Check if the validation error is nil
	if validationError == nil {
		return
	}

	// Check if the errors are nil
	if f.errors == nil {
		errors := make([]error, 0)
		f.errors = &errors
	}

	// Add the validation error to the field
	*f.errors = append(*f.errors, validationError)
}

// GetErrors returns the field errors
func (f *FieldValidations) GetErrors() *[]error {
	return f.errors
}
