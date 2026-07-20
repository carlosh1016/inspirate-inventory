package sesiones

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
)

// Entrada handles POST /api/v1/sesiones-laborales/entrada. Vendedora-only,
// enforced by the router — a vendedora clocks herself in, never someone
// else.
func (h *Handler) Entrada(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	sesion, err := h.service.Entrada(r.Context(), usecase.EntradaInput{SedeID: requester.SedeID, UsuarioID: requester.ID})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusCreated, toSesionResponse(*sesion))
}
