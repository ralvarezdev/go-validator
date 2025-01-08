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
	structName,
	fieldName string,
	fieldType reflect.Type,
	fieldValue interface{},
	required bool,
) {
	l.logger.Debug(
		fmt.Sprintf("detected field on struct '%v'", structName),
		fmt.Sprintf("field '%v' is initialized", fieldName),
		fmt.Sprintf("type: '%v'", fieldType),
		fmt.Sprintf("value: '%v'", fieldValue),
		fmt.Sprintf("required: '%v'", required),
	)
}

// UninitializedField prints the uninitialized field on debug mode
func (l *Logger) UninitializedField(
	structName, fieldName string,
	fieldType reflect.Type,
	required bool,
) {
	l.logger.Debug(
		fmt.Sprintf("detected field on struct '%v'", structName),
		fmt.Sprintf("field '%v' is uninitialized", fieldName),
		fmt.Sprintf("type: '%v'", fieldType),
		fmt.Sprintf("required: '%v'", required),
	)
}
