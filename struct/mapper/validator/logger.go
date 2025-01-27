package validator

import (
	"fmt"
	gologgermode "github.com/ralvarezdev/go-logger/mode"
	gologgermodenamed "github.com/ralvarezdev/go-logger/mode/named"
	"reflect"
)

// Logger is the structs mapper validator logger
type Logger struct {
	logger gologgermodenamed.Logger
}

// NewLogger creates a new structs mapper validator logger
func NewLogger(header string, modeLogger gologgermode.Logger) (*Logger, error) {
	// Initialize the mode named logger
	namedLogger, err := gologgermodenamed.NewDefaultLogger(header, modeLogger)
	if err != nil {
		return nil, err
	}

	return &Logger{logger: namedLogger}, nil
}

// InitializedField prints the initialized field on debug mode
func (l *Logger) InitializedField(
	structTypeName,
	fieldName string,
	fieldType reflect.Type,
	fieldValue interface{},
	required bool,
) {
	l.logger.Debug(
		"detected initialized field on struct type: "+structTypeName,
		"field name: "+fieldName,
		fmt.Sprintf("field type: '%v'", fieldType),
		fmt.Sprintf("field value: '%v'", fieldValue),
		fmt.Sprintf("field is required: '%v'", required),
	)
}

// UninitializedField prints the uninitialized field on debug mode
func (l *Logger) UninitializedField(
	structTypeName, fieldName string,
	fieldType reflect.Type,
	required bool,
) {
	l.logger.Debug(
		"detected uninitialized field on struct type: "+structTypeName,
		"field name: "+fieldName,
		fmt.Sprintf("field type: '%v'", fieldType),
		fmt.Sprintf("field is required: '%v'", required),
	)
}

// FieldTagNameNotFound prints the field tag name not found on debug mode
func (l *Logger) FieldTagNameNotFound(
	structTypeName, fieldName string,
	fieldType reflect.Type,
	fieldValue interface{},
) {
	l.logger.Debug(
		"field tag name not found on struct type: "+structTypeName,
		"field name: "+fieldName,
		fmt.Sprintf("field type: '%v'", fieldType),
		fmt.Sprintf("field value: '%v'", fieldValue),
	)
}
