package mapper

import (
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
			"Detected field on struct type",
			slog.String("struct_type", structTypeName),
			slog.String("field_name", fieldName),
			slog.Any("field_type", fieldType.String()),
			slog.String("field_tag", tag),
			slog.Bool("field_is_required", required),
		)
	}
}
