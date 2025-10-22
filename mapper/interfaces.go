package mapper

type (
	// Generator is an interface for creating a mapper
	Generator interface {
		NewMapper(structInstance any) (*Mapper, error)
		NewMapperWithNoError(structInstance any) *Mapper
	}
)
