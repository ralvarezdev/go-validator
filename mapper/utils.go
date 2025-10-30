package mapper

import (
	"fmt"
	"reflect"
	"strings"
)

// IsProtobufGeneratedField checks if the field is a Protobuf generated field
//
// Parameters:
//
//  - fieldName: The name of the struct field
//
// Returns:
//
// - bool: true if the field is a Protobuf generated field
func IsProtobufGeneratedField(fieldName string) bool {
	return fieldName == State || fieldName == SizeCache || fieldName == UnknownFields
}

// GetJSONTagName returns the JSON tag name for a given struct field name
// 
// Parameters:
// 
// - jsonTag: The JSON tag of the struct field
// - fieldName: The name of the struct field
//
// Returns:
// 
// - string: JSON tag name
// - error: error if the JSON tag is empty
func GetJSONTagName(jsonTag, fieldName string) (string, error) {
	// Get the field name from the JSON tag
	tagParts := strings.Split(jsonTag, ",")
	if len(tagParts) == 0 {
		return "", fmt.Errorf(ErrEmptyJSONTag, fieldName)
	}
	return tagParts[0], nil
}

// GetProtobufTag gets the Protobuf tag number for a given struct field name
// 
// Parameters:
// 
// - structField: The struct field
// - fieldName: The name of the struct field
// 
// Returns:
// 
// - string: Protobuf tag
// - error: error if the Protobuf tag is not found
func GetProtobufTag(structField reflect.StructField, fieldName string) (string, error) {
	protobufTag := structField.Tag.Get(ProtobufTag)
	if protobufTag == "" {
		return "", fmt.Errorf(ErrProtobufTagNotFound, fieldName)
	}
	return protobufTag, nil
}

// IsProtobufOneOfField checks if the struct field is a Protobuf oneof field
// 
// Parameters:
// 
// - structField: The struct field
// 
// Returns:
// 
// - bool: true if the struct field is a Protobuf oneof field
func IsProtobufOneOfField(structField reflect.StructField) bool {
	return structField.Tag.Get(ProtobufOneOfTag) != ""
}

// GetProtobufTagName returns the Protobuf tag name for a given struct field name
//
// Parameters:
//
//   - protobufTag: The Protobuf tag of the struct field
//   - fieldName: The name of the struct field
//
// Returns:
//
//   - string: Protobuf tag name
//   
func GetProtobufTagName(protobufTag, fieldName string) (string, error) {
	for _, part := range strings.Split(protobufTag, ",") {
		if strings.HasPrefix(part, ProtobufNamePrefix) {
			return strings.TrimPrefix(part, ProtobufNamePrefix), nil
		}
	}
	return "", fmt.Errorf(ErrProtobufTagNameNotFound, fieldName)
}

// IsProtobufFieldOptional returns true if the Protobuf field is optional
//
// Parameters:
//
//   - protobufTag: The Protobuf tag of the struct field
//
// Returns:
//
//   - bool: true if the Protobuf field is optional
func IsProtobufFieldOptional(protobufTag string) bool {
	return strings.Contains(
		protobufTag,
		ProtobufOneOf,
	)
}
