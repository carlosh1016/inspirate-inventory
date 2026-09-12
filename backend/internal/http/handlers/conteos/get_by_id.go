package conteos

import (
	"net/http"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// GetByID handles GET /api/v1/conteos-inventario/:id. Admin-only, enforced
// by the router.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}

	conteo, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toConteoResponse(*conteo))
}
