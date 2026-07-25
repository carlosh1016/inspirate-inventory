package reportes

import (
	"net/http"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Ventas handles GET /api/v1/reportes/ventas.
func (h *Handler) Ventas(w http.ResponseWriter, r *http.Request) {
	req, ok := requester(w, r)
	if !ok {
		return
	}
	params, err := h.parseReporteParams(r)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	ctx, cancel := withTimeout(r)
	defer cancel()

	data, err := h.service.GenerarVentas(ctx, req.SedeID, params)
	if err != nil {
		response.WriteError(w, r, mapGenerarErr(err))
		return
	}
	respondXLSX(w, r, filenameRango("ventas", params), data)
}
