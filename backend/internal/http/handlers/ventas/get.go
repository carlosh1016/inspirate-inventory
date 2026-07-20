package ventas

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Get handles GET /api/v1/ventas/:id. A vendedora may only see her own
// ventas — if she requests someone else's, this responds 404 (not 403) so
// as not to reveal that the venta exists.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	result, err := h.service.Get(r.Context(), id)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	if requester.Rol == RolVendedora && result.Venta.UsuarioID != requester.ID {
		response.WriteError(w, r, domainerrors.NewNotFound("Venta no encontrada", "La venta solicitada no existe."))
		return
	}

	response.WriteData(w, http.StatusOK, toVentaDetalladaResponse(result.Venta, result.MovimientosGenerados))
}
