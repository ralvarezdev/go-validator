package validation

import (
	goreflect "github.com/ralvarezdev/go-reflect"
)

type (
	// StructValidations is a struct that holds the struct validations for the generated validations of a struct
	StructValidations struct {
		fieldName                *string
		reflection               *goreflect.Reflection
		fieldsValidations        *map[string]*FieldValidations
		nestedStructsValidations *map[string]*StructValidations
	}

	// FieldValidations is a struct that holds the field validations for the generated validations of a struct
	FieldValidations struct {
		errors *[]error
	}
)

// NewStructValidations creates a new StructValidations struct
//
// Parameters:
//
//   - instance: The instance of the struct to validate
//
// Returns:
//
//   - *StructValidations: The StructValidations struct
//   - error: An error if the instance is nil
func NewStructValidations(instance interface{}) (*StructValidations, error) {
	// Check if the instance is nil
	if instance == nil {
		return nil, ErrNilInstance
	}

	return &StructValidations{
		reflection: goreflect.NewDereferencedReflection(instance),
	}, nil
}

// NewNestedStructValidations creates a new nested StructValidations struct
//
// Parameters:
//
//   - fieldName: The name of the field that holds the nested struct
//   - instance: The instance of the nested struct to validate
//
// Returns:
//
//   - *StructValidations: The StructValidations struct
//   - error: An error if the field name is empty or the instance is nil
func NewNestedStructValidations(
	fieldName string,
	instance interface{},
) (*StructValidations, error) {
	// Check if the field name is empty or the instance is nil
	if instance == nil {
		return nil, ErrNilInstance
	}
	if fieldName == "" {
		return NewStructValidations(instance)
	}

	return &StructValidations{
		fieldName:  &fieldName,
		reflection: goreflect.NewDereferencedReflection(instance),
	}, nil
}

// GetReflection returns the reflection of the struct
//
// Returns:
//
//   - *goreflect.Reflection: The reflection of the struct
func (s *StructValidations) GetReflection() *goreflect.Reflection {
	return s.reflection
}

// GetStructTypeName returns the type name of the struct
//
// Returns:
//
//   - string: The type name of the struct
func (s *StructValidations) GetStructTypeName() string {
	return s.reflection.GetReflectedTypeName()
}

// HasFailed returns true if there are failed validations
//
// Returns:
//
//   - bool: True if there are failed validations, false otherwise
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
//
// Parameters:
//
//   - fieldName: The name of the field to add the validations to
//   - fieldValidations: The field validations to add
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
//
// Parameters:
//
//   - fieldName: The name of the field to add the validation error to
//   - validationError: The validation error to add
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
		fieldValidations = NewFieldValidations()
		(*s.fieldsValidations)[fieldName] = fieldValidations
	}

	// Append the validation error to the field name
	fieldValidations.AddValidationError(validationError)
}

// AddNestedStructValidations sets the nested struct fields validations to the struct
//
// Parameters:
//
//   - fieldName: The name of the field that holds the nested struct
//   - nestedStructValidations: The nested struct validations to add
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
//
// Returns:
//
//   - *map[string]*FieldValidations: The fields validations
func (s *StructValidations) GetFieldsValidations() *map[string]*FieldValidations {
	return s.fieldsValidations
}

// GetNestedStructsValidations returns the nested structs validations
//
// Returns:
//
//   - *map[string]*StructValidations: The nested structs validations
func (s *StructValidations) GetNestedStructsValidations() *map[string]*StructValidations {
	return s.nestedStructsValidations
}

// NewFieldValidations creates a new FieldValidations struct
//
// Returns:
//
//   - *FieldValidations: The FieldValidations struct
func NewFieldValidations() *FieldValidations {
	return &FieldValidations{}
}

// HasFailed returns true if there are failed validations
//
// Returns:
//
//   - bool: True if there are failed validations, false otherwise
func (f *FieldValidations) HasFailed() bool {
	if f.errors == nil {
		return false
	}
	return len(*f.errors) > 0
}

// AddValidationError adds a validation error to the field
//
// Parameters:
//
//   - validationError: The validation error to add
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
//
// Returns:
//
//   - *[]error: The field errors
func (f *FieldValidations) GetErrors() *[]error {
	return f.errors
}
