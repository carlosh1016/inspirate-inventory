package conteos

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/conteos"
)

// Cerrar handles POST /api/v1/conteos-inventario/:id/cerrar. Admin-only,
// enforced by the router. A conteo is immutable once cerrado — there is no
// reopen endpoint, so a second call against the same id returns 409.
func (h *Handler) Cerrar(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}

	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	conteo, err := h.service.Cerrar(r.Context(), usecase.CerrarInput{
		TargetID:    id,
		RequesterID: requester.ID,
		IP:          middleware.IPFromContext(r.Context()),
		UserAgent:   middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toConteoResponse(*conteo))
}
