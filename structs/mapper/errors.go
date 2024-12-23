package mapper

var (
	MissingProtobufTagError                 = "missing protobuf tag: %s"
	MissingProtobufTagNameError             = "missing protobuf tag name: %s"
	DuplicateProtobufTagNameError           = "duplicate protobuf tag name: %s"
	MissingJSONTagError                     = "missing json tag: %s"
	MissingJSONTagLooksLikeProtocFieldError = "missing json tag, looks like a protoc field: %s. If it is a protoc field, use ProtobufGenerator instead"
	EmptyJSONTagError                       = "empty json tag: %s"
	DuplicateJSONTagNameError               = "duplicate json tag name: %s"
)
