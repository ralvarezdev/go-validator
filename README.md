# go-validator

Struct validator for Go based on reflection. It builds a mapper of a struct from its `json` tags or protobuf fields, checks required fields, collects per-field errors and converts them into a JSON-friendly map or gRPC `BadRequest` error details.

Module `github.com/ralvarezdev/go-validator` (Go 1.24.0). Docs: <https://pkg.go.dev/github.com/ralvarezdev/go-validator>.

## Installation

```bash
go get github.com/ralvarezdev/go-validator
```

Depends on `go-reflect`, `go-strings` (both by the same author) and `google.golang.org/genproto/googleapis/rpc`.

## Concepts

- **Mapper** (`mapper`) — reflection-based description of a struct. A field is required unless its JSON tag is missing or contains `omitempty`. `NewJSONGenerator(logger)` and `NewProtobufGenerator(logger)` implement `Generator` (`NewMapper`, `NewMapperWithNoError`).
- **Validations** (`mapper/validation`) — `StructValidations` and `FieldValidations` hold the errors collected for an instance, including nested structs.
- **Validator** (`mapper/validator`) — `DefaultValidator` checks that required fields are initialised. `NewDefaultService(rawParser, endParser, validator, birthdateOptions, passwordOptions, logger)` returns the `Service`.
- **Parsers** (`mapper/parser`) — a `RawParser` flattens validations and an `EndParser` renders them. Implementations in `parser/json` and `parser/grpc` (`ErrorDetails.GetBadRequest()`).
- **Field validators** — `Email`, `Username`, `Birthdate`, `Password` on `Service`, configured with `BirthdateOptions{MinimumAge, MaximumAge}` and `PasswordOptions{MinimumLength, MinimumSpecialCount, MinimumNumbersCount, MinimumCapsCount}`. Error values live in `field/`.
- **Running** — `Service.Validate(mapper, auxiliaryValidatorFns...)` and `Service.CreateValidateFn(mapper, cache, auxiliaryValidatorFns...)`.

## Usage

```go
type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname,omitempty"` // optional
}

generator := govalidatormapper.NewJSONGenerator(nil)
m, err := generator.NewMapper(CreateUser{})
```

Pass the mapper to a `validator.Service` (see `mapper/validator/service.go`). No end-to-end example is included in the repository.

## Project structure

```
field/        error values for birthdate, mail, password, username
mapper/       Mapper, JSON and protobuf generators
  parser/     parser interfaces; json/ and grpc/ implementations
  validation/ StructValidations and FieldValidations
  validator/  Validator, Service and DefaultService
```

The repository has no tests; `.golangci.yml` is provided for linting.

## License

GNU General Public License v3.0.
