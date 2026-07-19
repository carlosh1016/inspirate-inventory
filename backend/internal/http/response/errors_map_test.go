package response_test

import (
	"net/http"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

func TestStatusForCode(t *testing.T) {
	tests := []struct {
		code   domainerrors.Code
		status int
	}{
		{domainerrors.CodeValidation, http.StatusUnprocessableEntity},
		{domainerrors.CodeBusinessRule, http.StatusUnprocessableEntity},
		{domainerrors.CodeUnauthorized, http.StatusUnauthorized},
		{domainerrors.CodeForbidden, http.StatusForbidden},
		{domainerrors.CodeNotFound, http.StatusNotFound},
		{domainerrors.CodeConflict, http.StatusConflict},
		{domainerrors.CodeRateLimit, http.StatusTooManyRequests},
		{domainerrors.CodeInternal, http.StatusInternalServerError},
		{domainerrors.Code("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		if got := response.StatusForCode(tt.code); got != tt.status {
			t.Errorf("StatusForCode(%q) = %d, want %d", tt.code, got, tt.status)
		}
	}
}
