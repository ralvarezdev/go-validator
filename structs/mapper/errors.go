package mapper

var (
	ErrMissingProtobufTag                 = "missing protobuf tag: %s"
	ErrMissingProtobufTagName             = "missing protobuf tag name: %s"
	ErrMissingJSONTag                     = "missing json tag: %s"
	ErrMissingJSONTagLooksLikeProtocField = "missing json tag, looks like a protoc field: %s. If it is a protoc field, use ProtobufGenerator instead"
	ErrEmptyJSONTag                       = "empty json tag: %s"
)
