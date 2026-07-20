package sesiones

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Salida handles POST /api/v1/sesiones-laborales/salida. Vendedora-only,
// enforced by the router — a vendedora closes her own currently open
// sesion.
func (h *Handler) Salida(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	sesion, err := h.service.Salida(r.Context(), requester.ID)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toSesionResponse(*sesion))
}
