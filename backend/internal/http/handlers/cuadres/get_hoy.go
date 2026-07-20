package cuadres

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// GetHoy handles GET /api/v1/cuadres-caja/hoy. Returns {"data": null} (not
// a 404) when no cuadre exists yet for today — the frontend interprets
// that as "open the day".
func (h *Handler) GetHoy(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	cuadre, err := h.service.GetHoy(r.Context(), requester.SedeID)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	if cuadre == nil {
		response.WriteData(w, http.StatusOK, nil)
		return
	}

	response.WriteData(w, http.StatusOK, toCuadreResponse(*cuadre))
}
