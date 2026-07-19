package usuarios

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
)

// Activate handles PATCH /api/v1/usuarios/:id/activar.
func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	err := h.service.Activate(r.Context(), usecase.ActivateInput{
		TargetID:    id,
		RequesterID: requester.ID,
		IP:          middleware.IPFromContext(r.Context()),
		UserAgent:   middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteNoContent(w)
}
