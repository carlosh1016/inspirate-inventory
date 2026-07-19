package response

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
)

// StatusForCode maps a domain error Code to the HTTP status it represents.
func StatusForCode(code domainerrors.Code) int {
	switch code {
	case domainerrors.CodeValidation, domainerrors.CodeBusinessRule:
		return http.StatusUnprocessableEntity
	case domainerrors.CodeUnauthorized:
		return http.StatusUnauthorized
	case domainerrors.CodeForbidden:
		return http.StatusForbidden
	case domainerrors.CodeNotFound:
		return http.StatusNotFound
	case domainerrors.CodeConflict:
		return http.StatusConflict
	case domainerrors.CodeRateLimit:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
