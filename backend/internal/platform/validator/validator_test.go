package validator_test

import (
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/validator"
)

type sampleRequest struct {
	NombreCompleto string `json:"nombre_completo" validate:"required,min=3,max=200"`
	Correo         string `json:"correo" validate:"required,email"`
	Rol            string `json:"rol" validate:"required,oneof=admin vendedora"`
}

func TestValidateSuccess(t *testing.T) {
	v := validator.New()
	err := v.Validate(sampleRequest{NombreCompleto: "Ana Pérez", Correo: "ana@inspirate.co", Rol: "admin"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateFailureUsesJSONFieldNames(t *testing.T) {
	v := validator.New()
	err := v.Validate(sampleRequest{NombreCompleto: "A", Correo: "not-an-email", Rol: "gerente"})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	var domainErr *domainerrors.DomainError
	if de, ok := err.(*domainerrors.DomainError); ok {
		domainErr = de
	} else {
		t.Fatalf("expected *DomainError, got %T", err)
	}

	if domainErr.Code != domainerrors.CodeValidation {
		t.Fatalf("expected CodeValidation, got %q", domainErr.Code)
	}

	for _, field := range []string{"nombre_completo", "correo", "rol"} {
		if len(domainErr.Fields[field]) == 0 {
			t.Errorf("expected error message for field %q, got none (fields=%+v)", field, domainErr.Fields)
		}
	}
}
