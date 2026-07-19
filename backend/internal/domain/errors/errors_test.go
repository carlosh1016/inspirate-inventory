package errors_test

import (
	stderrors "errors"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
)

func TestConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *domainerrors.DomainError
		code domainerrors.Code
	}{
		{"validation", domainerrors.NewValidation("t", "d", map[string][]string{"correo": {"inválido"}}), domainerrors.CodeValidation},
		{"not_found", domainerrors.NewNotFound("t", "d"), domainerrors.CodeNotFound},
		{"conflict", domainerrors.NewConflict("t", "d"), domainerrors.CodeConflict},
		{"business_rule", domainerrors.NewBusinessRule("t", "d"), domainerrors.CodeBusinessRule},
		{"unauthorized", domainerrors.NewUnauthorized("t", "d"), domainerrors.CodeUnauthorized},
		{"forbidden", domainerrors.NewForbidden("t", "d"), domainerrors.CodeForbidden},
		{"rate_limit", domainerrors.NewRateLimit("t", "d"), domainerrors.CodeRateLimit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Fatalf("expected code %q, got %q", tt.code, tt.err.Code)
			}
			if tt.err.Error() != "d" {
				t.Fatalf("expected Error() to return detail %q, got %q", "d", tt.err.Error())
			}
		})
	}
}

func TestErrorFallsBackToTitle(t *testing.T) {
	err := &domainerrors.DomainError{Code: domainerrors.CodeInternal, Title: "solo titulo"}
	if err.Error() != "solo titulo" {
		t.Fatalf("expected Error() to fall back to title, got %q", err.Error())
	}
}

func TestNewInternalWrapsAndUnwraps(t *testing.T) {
	cause := stderrors.New("db connection refused")
	err := domainerrors.NewInternal("Error interno", "Ocurrió un error inesperado", cause)

	if err.Code != domainerrors.CodeInternal {
		t.Fatalf("expected CodeInternal, got %q", err.Code)
	}
	if !stderrors.Is(err, cause) {
		t.Fatalf("expected errors.Is to unwrap to cause")
	}
}
