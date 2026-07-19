package productos

import (
	"net/http"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Get handles GET /api/v1/productos/:id.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	p, err := h.service.Get(r.Context(), id)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toProductoResponseFromGet(p))
}
