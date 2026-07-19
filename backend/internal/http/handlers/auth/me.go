package auth

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Me handles GET /api/v1/auth/me.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized(
			"No autenticado",
			"Debes iniciar sesión para continuar.",
		))
		return
	}

	usuario, err := h.service.Me(r.Context(), user.ID)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toUsuarioResponse(usuario))
}
