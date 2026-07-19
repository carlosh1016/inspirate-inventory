package stock

import "net/http"

// Alertas handles GET /api/v1/stock/alertas — equivalent to
// GET /api/v1/stock?stock_bajo=true, exposed as its own endpoint so the
// frontend doesn't need to know the equivalence.
func (h *Handler) Alertas(w http.ResponseWriter, r *http.Request) {
	h.list(w, r, true)
}
