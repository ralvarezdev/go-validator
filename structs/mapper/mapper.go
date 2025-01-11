package mapper

import (
	"fmt"
	goreflect "github.com/ralvarezdev/go-reflect"
	"reflect"
	"strings"
)

// Protobuf fields generated by the protoc compiler
const (
	State              = "state"
	SizeCache          = "sizeCache"
	UnknownFields      = "unknownFields"
	ProtobufTag        = "protobuf"
	ProtobufOneOf      = "oneof"
	ProtobufNamePrefix = "name="
	JSONTag            = "json"
	JSONOmitempty      = "omitempty"
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

	// Generator is an interface for creating a mapper
	Generator interface {
		NewMapper(structInstance interface{}) (*Mapper, error)
	}

	// ProtobufGenerator is a generator for Protobuf mappers
	ProtobufGenerator struct {
		logger *Logger
	}

	// JSONGenerator is a generator for JSON mappers
	JSONGenerator struct {
		logger *Logger
	}
)

// NewMapper creates a new mapper
func NewMapper(structInstance interface{}) *Mapper {
	return &Mapper{structInstance: structInstance}
}

// GetStructInstance returns the instance of the struct
func (m *Mapper) GetStructInstance() interface{} {
	return m.structInstance
}

// Type returns the type of the struct instance
func (m *Mapper) Type() reflect.Type {
	return reflect.TypeOf(m.structInstance)
}

// GetFieldsTagName returns the fields of the mapper
func (m *Mapper) GetFieldsTagName() *map[string]string {
	return m.fields
}

// GetFieldTagName returns the tag name of a field
func (m *Mapper) GetFieldTagName(fieldName string) (
	string,
	bool,
) {
	// Check if the fields map is nil
	if m.fields == nil {
		return "", false
	}

	// Check if the field name exists in the map
	fieldTagName, ok := (*m.fields)[fieldName]
	return fieldTagName, ok
}

// AddFieldTagName adds a field tag name to the mapper
func (m *Mapper) AddFieldTagName(fieldName, fieldTagName string) {
	// Initialize the fields map if it is nil
	if m.fields == nil {
		m.fields = &map[string]string{}
	}

	// Add the field tag name to the map
	(*m.fields)[fieldName] = fieldTagName
}

// GetRequiredFields returns the required fields of the mapper
func (m *Mapper) GetRequiredFields() *map[string]bool {
	return m.requiredFields
}

// IsFieldRequired returns if a field is required
func (m *Mapper) IsFieldRequired(fieldName string) (bool, bool) {
	// Check if the required fields map is nil
	if m.requiredFields == nil {
		return false, false
	}

	// Check if the field name exists in the map
	isFieldRequired, ok := (*m.requiredFields)[fieldName]
	return isFieldRequired, ok
}

// SetFieldIsRequired sets if a field is required
func (m *Mapper) SetFieldIsRequired(fieldName string, required bool) {
	// Initialize the required fields map if it is nil
	if m.requiredFields == nil {
		m.requiredFields = &map[string]bool{}
	}

	// Set if the field is required
	(*m.requiredFields)[fieldName] = required
}

// HasFieldsValidations returns if the mapper has fields
func (m *Mapper) HasFieldsValidations() bool {
	return m.fields != nil
}

// GetNestedMappers returns the nested mappers of the mapper
func (m *Mapper) GetNestedMappers() *map[string]*Mapper {
	return m.nestedMappers
}

// GetFieldNestedMapper returns the nested mapper of a field
func (m *Mapper) GetFieldNestedMapper(fieldName string) *Mapper {
	// Check if the nested mappers map is nil
	if m.nestedMappers == nil {
		return nil
	}

	return (*m.nestedMappers)[fieldName]
}

// AddFieldNestedMapper adds a nested mapper to the mapper
func (m *Mapper) AddFieldNestedMapper(fieldName string, nestedMapper *Mapper) {
	// Initialize the nested mappers map if it is nil
	if m.nestedMappers == nil {
		m.nestedMappers = &map[string]*Mapper{}
	}

	// Add the nested mapper to the map
	(*m.nestedMappers)[fieldName] = nestedMapper
}

// NewProtobufGenerator creates a new Protobuf generator
func NewProtobufGenerator(logger *Logger) *ProtobufGenerator {
	return &ProtobufGenerator{
		logger: logger,
	}
}

// NewJSONGenerator creates a new JSON generator
func NewJSONGenerator(logger *Logger) *JSONGenerator {
	return &JSONGenerator{
		logger: logger,
	}
}

// NewMapper creates the fields to validate from a Protobuf compiled struct
func (p *ProtobufGenerator) NewMapper(structInstance interface{}) (
	*Mapper,
	error,
) {
	// Reflection of data
	reflectedType := goreflect.GetDereferencedType(structInstance)

	// Get the struct type name
	structTypeName := goreflect.GetTypeName(reflectedType)

	// Initialize the root map of fields and the map of nested mappers
	rootMapper := NewMapper(structInstance)

	// Reflection of the type of data
	for i := 0; i < reflectedType.NumField(); i++ {
		// Get the field type through reflection
		structField := reflectedType.Field(i)
		fieldType := structField.Type
		fieldName := structField.Name

		// Check if the field is a protoc generated field
		if fieldName == State || fieldName == SizeCache || fieldName == UnknownFields {
			continue
		}

		// Get the Protobuf tag of the field
		protobufTag := structField.Tag.Get(ProtobufTag)
		if protobufTag == "" {
			return nil, fmt.Errorf(ErrProtobufTagNotFound, fieldName)
		}

		// Get the field name from the Protobuf tag
		var protobufName string
		for _, part := range strings.Split(protobufTag, ",") {
			if strings.HasPrefix(part, ProtobufNamePrefix) {
				protobufName = strings.TrimPrefix(part, ProtobufNamePrefix)
				break
			}
		}

		// Check if the field name is empty
		if protobufName == "" {
			return nil, fmt.Errorf(ErrProtobufTagNameNotFound, fieldName)
		}

		// Add the field to the fields map
		rootMapper.AddFieldTagName(fieldName, protobufName)

		// Check if the field is a pointer
		if fieldType.Kind() == reflect.Ptr {
			// Dereference the pointer
			fieldType = fieldType.Elem()

			// Check if the element type is not a struct and the tag to determine if it contains 'oneof', which means it is an optional struct field
			if fieldType.Kind() != reflect.Struct || strings.Contains(
				protobufTag,
				ProtobufOneOf,
			) {
				// Set field as not required
				rootMapper.SetFieldIsRequired(fieldName, false)

				// Print field
				if p.logger != nil {
					p.logger.DetectedField(
						structTypeName,
						fieldName,
						fieldType,
						protobufTag,
						false,
					)
				}
				continue
			}

			// Create a new Mapper for the nested struct field
			fieldNestedMapper, err := p.NewMapper(
				reflect.New(fieldType).Interface(),
			)
			if err != nil {
				return nil, err
			}

			// Add the nested fields to the map
			rootMapper.AddFieldNestedMapper(fieldName, fieldNestedMapper)
		}

		// Set field as required
		rootMapper.SetFieldIsRequired(fieldName, true)

		// Print field
		if p.logger != nil {
			p.logger.DetectedField(
				structTypeName,
				fieldName,
				fieldType,
				protobufTag,
				true,
			)
		}
	}

	return rootMapper, nil
}

// NewMapper creates the fields to validate from a JSON struct
func (j *JSONGenerator) NewMapper(structInstance interface{}) (
	*Mapper,
	error,
) {
	// Reflection of data
	reflectedType := goreflect.GetDereferencedType(structInstance)

	// Get the struct type name
	structTypeName := reflectedType.Name()

	// Initialize the root map of fields and the map of nested mappers
	rootMapper := NewMapper(structInstance)

	// Reflection of the type of data
	var jsonTag string
	var jsonName string
	for i := 0; i < reflectedType.NumField(); i++ {
		// Get the field type through reflection
		structField := reflectedType.Field(i)
		fieldType := structField.Type
		fieldTag := structField.Tag
		fieldName := structField.Name

		// Get the JSON tag of the field
		jsonTag = fieldTag.Get(JSONTag)

		// Get the field name from the JSON tag
		tagParts := strings.Split(jsonTag, ",")
		if len(tagParts) == 0 {
			return nil, fmt.Errorf(ErrEmptyJSONTag, fieldName)
		}
		jsonName = tagParts[0]

		// Add field tag name to the map and set the field as parsed
		rootMapper.AddFieldTagName(fieldName, jsonName)

		// Check if the JSON tag is unassigned or if it contains 'omitempty', which means it is an optional field
		if jsonTag == "-" || strings.Contains(jsonTag, JSONOmitempty) {
			// Set field name as not required
			rootMapper.SetFieldIsRequired(fieldName, false)

			// Print field
			if j.logger != nil {
				j.logger.DetectedField(
					structTypeName,
					fieldName,
					fieldType,
					jsonTag,
					false,
				)
			}
			continue
		}

		// Set field name as required
		rootMapper.SetFieldIsRequired(fieldName, true)

		// Dereference the pointer
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// Check if the element type is a struct
		if fieldType.Kind() == reflect.Struct {
			// Create a new Mapper for the nested struct field
			fieldNestedMapper, err := j.NewMapper(
				reflect.New(fieldType).Interface(),
			)
			if err != nil {
				return nil, err
			}

			// Add the nested fields to the map
			rootMapper.AddFieldNestedMapper(fieldName, fieldNestedMapper)
		}

		// Print field
		if j.logger != nil {
			j.logger.DetectedField(
				structTypeName,
				fieldName,
				fieldType,
				jsonTag,
				true,
			)
		}
	}

	return rootMapper, nil
}
