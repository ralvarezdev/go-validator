package mapper

import (
	"fmt"
	"log/slog"
	"reflect"
)

// DetectedField prints a detected field
//
// Parameters:
//
// - structTypeName: the name of the struct type
// - fieldName: the name of the field
// - fieldType: the type of the field
// - tag: the tag of the field
// - required: whether the field is required or not
// - logger: the logger to use
func DetectedField(
	structTypeName string,
	fieldName string,
	fieldType reflect.Type,
	tag string,
	required bool,
	logger *slog.Logger,
) {
	if logger != nil {
		logger.Debug(
			"detected field on struct type: "+structTypeName,
			slog.String("fieldName", fieldName),
			slog.Any("fieldType", fmt.Sprintf("%v", fieldType)),
			slog.String("fieldTag", tag),
			slog.Bool("fieldRequired", required),
		)
	}
}
