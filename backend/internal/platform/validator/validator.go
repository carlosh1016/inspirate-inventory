// Package validator wraps go-playground/validator/v10 to turn struct-tag
// validation failures into a domain-level *DomainError with Spanish,
// per-field messages keyed by json tag (not the Go field name).
package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	playground "github.com/go-playground/validator/v10"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
)

// Validator validates structs using `validate:"..."` tags.
type Validator struct {
	v *playground.Validate
}

// New builds a Validator whose field errors are keyed by the struct's json
// tag, so error payloads match the snake_case field names clients send.
func New() *Validator {
	v := playground.New(playground.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{v: v}
}

// Validate runs struct validation. On failure it returns a *DomainError
// (Code=CodeValidation) with one Spanish message per invalid field. Errors
// that aren't validation failures (e.g. an unsupported type) are returned
// as-is so they surface as bugs, not user-facing 422s.
func (val *Validator) Validate(s any) error {
	err := val.v.Struct(s)
	if err == nil {
		return nil
	}

	var fieldErrs playground.ValidationErrors
	if !errors.As(err, &fieldErrs) {
		return err
	}

	fields := make(map[string][]string, len(fieldErrs))
	for _, fe := range fieldErrs {
		fields[fe.Field()] = append(fields[fe.Field()], messageFor(fe))
	}

	return domainerrors.NewValidation("Datos inválidos", "Revisa los campos marcados.", fields)
}

func messageFor(fe playground.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Este campo es obligatorio."
	case "email":
		return "Debe ser un correo electrónico válido."
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Debe tener al menos %s caracteres.", fe.Param())
		}
		return fmt.Sprintf("Debe ser al menos %s.", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Debe tener como máximo %s caracteres.", fe.Param())
		}
		return fmt.Sprintf("Debe ser como máximo %s.", fe.Param())
	case "oneof":
		return fmt.Sprintf("Debe ser uno de: %s.", fe.Param())
	default:
		return "Valor inválido."
	}
}
