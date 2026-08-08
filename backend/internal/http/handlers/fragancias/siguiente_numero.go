package fragancias

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// SiguienteNumero handles GET /api/v1/fragancias/siguiente-numero. Suggests
// the next numero_genero for the requested genero — a UI default, not
// enforced: create/update still accept any explicit number.
func (h *Handler) SiguienteNumero(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	siguiente, err := h.service.SiguienteNumero(r.Context(), requester.SedeID, r.URL.Query().Get("genero"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, SiguienteNumeroResponse{Siguiente: siguiente})
}
