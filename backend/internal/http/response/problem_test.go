package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

func TestWriteErrorWithDomainError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/usuarios", nil)

	err := domainerrors.NewValidation(
		"Datos inválidos",
		"Revisa los campos marcados",
		map[string][]string{"correo": {"Debe ser un correo válido"}},
	)
	response.WriteError(w, r, err)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/problem+json") {
		t.Fatalf("unexpected content-type: %q", ct)
	}

	var problem response.Problem
	if decodeErr := json.Unmarshal(w.Body.Bytes(), &problem); decodeErr != nil {
		t.Fatalf("invalid json: %v", decodeErr)
	}
	if problem.Type != "https://api.inspirate.co/errors/validation_error" {
		t.Fatalf("unexpected type: %q", problem.Type)
	}
	if problem.Status != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected status field: %d", problem.Status)
	}
	if problem.Instance != "/api/v1/usuarios" {
		t.Fatalf("unexpected instance: %q", problem.Instance)
	}
	if len(problem.Errors["correo"]) != 1 {
		t.Fatalf("expected field errors for correo, got %+v", problem.Errors)
	}
}

func TestWriteErrorWithGenericErrorDoesNotLeakDetail(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios/1", nil)

	response.WriteError(w, r, errors.New("pq: connection refused on 10.0.0.5"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

	var problem response.Problem
	if err := json.Unmarshal(w.Body.Bytes(), &problem); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if strings.Contains(problem.Detail, "10.0.0.5") || strings.Contains(problem.Detail, "connection refused") {
		t.Fatalf("internal error details leaked to client: %q", problem.Detail)
	}
}
