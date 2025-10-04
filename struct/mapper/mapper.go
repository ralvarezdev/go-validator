package mapper

import (
	"reflect"
)

type (
	// Mapper is a map of fields to validate from a struct
	Mapper struct {
		// structInstance is the instance of the struct
		structInstance interface{}

		// fields key is the field name and value is the tag name
		fields *map[string]string

		// requiredFields key is the field name and value is a boolean to determine if the field is required
		requiredFields *map[string]bool

		// nestedMappers key is the field name of the nested struct and value is the nested mapper
		nestedMappers *map[string]*Mapper
	}
)

// NewMapper creates a new mapper
//
// Parameters:
//
//   - structInstance: instance of the struct to create the mapper from
//
// Returns:
//
//   - *Mapper: instance of the mapper
func NewMapper(structInstance interface{}) *Mapper {
	return &Mapper{structInstance: structInstance}
}

// GetStructInstance returns the instance of the struct
//
// Returns:
//
//   - interface{}: instance of the struct
func (m *Mapper) GetStructInstance() interface{} {
	if m == nil {
		return nil
	}
	return m.structInstance
}

// Type returns the type of the struct instance
//
// Returns:
//
//   - reflect.Type: type of the struct instance
func (m *Mapper) Type() reflect.Type {
	if m == nil {
		return nil
	}
	return reflect.TypeOf(m.structInstance)
}

// GetFieldsTagName returns the fields of the mapper
//
// Returns:
//
//   - *map[string]string: map of fields where key is the field name and value is the tag name
func (m *Mapper) GetFieldsTagName() *map[string]string {
	if m == nil {
		return nil
	}
	return m.fields
}

// GetFieldTagName returns the tag name of a field
//
// Parameters:
//
//   - fieldName: name of the field
//
// Returns:
//
//   - string: tag name of the field
//   - bool: true if the field exists, false otherwise
func (m *Mapper) GetFieldTagName(fieldName string) (
	string,
	bool,
) {
	if m == nil {
		return "", false
	}

	// Check if the fields map is nil
	if m.fields == nil {
		return "", false
	}

	// Check if the field name exists in the map
	fieldTagName, ok := (*m.fields)[fieldName]
	return fieldTagName, ok
}

// AddFieldTagName adds a field tag name to the mapper
//
// Parameters:
//
//   - fieldName: name of the field
//   - fieldTagName: tag name of the field
func (m *Mapper) AddFieldTagName(fieldName, fieldTagName string) {
	if m == nil {
		return
	}

	// Initialize the fields map if it is nil
	if m.fields == nil {
		m.fields = &map[string]string{}
	}

	// Add the field tag name to the map
	(*m.fields)[fieldName] = fieldTagName
}

// GetRequiredFields returns the required fields of the mapper
//
// Returns:
//
//   - *map[string]bool: map of required fields where key is the field name and value is a boolean to determine if the field is required
func (m *Mapper) GetRequiredFields() *map[string]bool {
	if m == nil {
		return nil
	}
	return m.requiredFields
}

// IsFieldRequired returns if a field is required
//
// Parameters:
//
//   - fieldName: name of the field
//
// Returns:
//
//   - bool: true if the field is required, false otherwise
//   - bool: true if the field exists, false otherwise
func (m *Mapper) IsFieldRequired(fieldName string) (bool, bool) {
	if m == nil {
		return false, false
	}

	// Check if the required fields map is nil
	if m.requiredFields == nil {
		return false, false
	}

	// Check if the field name exists in the map
	isFieldRequired, ok := (*m.requiredFields)[fieldName]
	return isFieldRequired, ok
}

// SetFieldIsRequired sets if a field is required
//
// Parameters:
//
//   - fieldName: name of the field
//   - required: true if the field is required, false otherwise
func (m *Mapper) SetFieldIsRequired(fieldName string, required bool) {
	if m == nil {
		return
	}

	// Initialize the required fields map if it is nil
	if m.requiredFields == nil {
		m.requiredFields = &map[string]bool{}
	}

	// Set if the field is required
	(*m.requiredFields)[fieldName] = required
}

// HasFieldsValidations returns if the mapper has fields
//
// Returns:
//
//   - bool: true if the mapper has fields, false otherwise
func (m *Mapper) HasFieldsValidations() bool {
	if m == nil {
		return false
	}
	return m.fields != nil
}

// GetNestedMappers returns the nested mappers of the mapper
//
// Returns:
//
//   - *map[string]*Mapper: map of nested mappers where key is the field name of the nested struct and value is the nested mapper
func (m *Mapper) GetNestedMappers() *map[string]*Mapper {
	if m == nil {
		return nil
	}
	return m.nestedMappers
}

// GetFieldNestedMapper returns the nested mapper of a field
//
// Parameters:
//
//   - fieldName: name of the field
//
// Returns:
//
//   - *Mapper: nested mapper of the field, or nil if the field does not exist or has no nested mapper
func (m *Mapper) GetFieldNestedMapper(fieldName string) *Mapper {
	if m == nil {
		return nil
	}

	// Check if the nested mappers map is nil
	if m.nestedMappers == nil {
		return nil
	}

	return (*m.nestedMappers)[fieldName]
}

// AddFieldNestedMapper adds a nested mapper to the mapper
//
// Parameters:
//
//   - fieldName: name of the field
//   - nestedMapper: nested mapper to add
func (m *Mapper) AddFieldNestedMapper(fieldName string, nestedMapper *Mapper) {
	if m == nil {
		return
	}

	// Initialize the nested mappers map if it is nil
	if m.nestedMappers == nil {
		m.nestedMappers = &map[string]*Mapper{}
	}

	// Add the nested mapper to the map
	(*m.nestedMappers)[fieldName] = nestedMapper
}
