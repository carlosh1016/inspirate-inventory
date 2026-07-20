// Package sesiones contains the HTTP handlers for
// /api/v1/sesiones-laborales/*.
package sesiones

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/validator"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
)

// RolVendedora identifies the restricted role: List is scoped to a
// vendedora's own sesiones.
const RolVendedora = "vendedora"

// Handler groups the sesiones laborales HTTP handlers and their
// dependencies.
type Handler struct {
	service    *usecase.Service
	jwtManager jwt.Manager
	validator  *validator.Validator
}

// NewHandler builds a Handler.
func NewHandler(service *usecase.Service, jwtManager jwt.Manager, v *validator.Validator) *Handler {
	return &Handler{service: service, jwtManager: jwtManager, validator: v}
}

func parseIDParam(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteError(w, r, domainerrors.NewValidation(
			"Solicitud inválida", "El identificador debe ser numérico.", nil,
		))
		return 0, false
	}
	return id, true
}

func badRequestBodyErr() error {
	return domainerrors.NewValidation("Solicitud inválida", "El cuerpo de la solicitud no es un JSON válido.", nil)
}
