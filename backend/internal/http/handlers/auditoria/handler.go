// Package auditoria contains the HTTP handlers for /api/v1/auditoria/* — two
// admin-only, read-only endpoints over the audit log.
package auditoria

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auditoria"
)

// Handler groups the auditoria HTTP handlers and their dependencies.
type Handler struct {
	service    *usecase.Service
	jwtManager jwt.Manager
}

// NewHandler builds a Handler.
func NewHandler(service *usecase.Service, jwtManager jwt.Manager) *Handler {
	return &Handler{service: service, jwtManager: jwtManager}
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
