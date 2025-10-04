package mapper

type (
	// Generator is an interface for creating a mapper
	Generator interface {
		NewMapper(structInstance interface{}) (*Mapper, error)
		NewMapperWithNoError(structInstance interface{}) *Mapper
	}
)
