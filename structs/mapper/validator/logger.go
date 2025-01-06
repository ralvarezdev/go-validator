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

// PrintField prints the field on debug mode
func (l *Logger) PrintField(
	fieldName string,
	fieldType reflect.Type,
	fieldValue interface{},
) {
	l.logger.Debug(
		fmt.Sprintf("field '%v'", fieldName),
		fmt.Sprintf("type: '%v'", fieldType),
		fmt.Sprintf("value: '%v'", fieldValue),
	)
}

// UninitializedField prints the uninitialized field on debug mode
func (l *Logger) UninitializedField(fieldName string) {
	l.logger.Debug(
		fmt.Sprintf("field '%v' is uninitialized", fieldName),
	)
}
