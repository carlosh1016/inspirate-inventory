package reportes

import (
	"net/http"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Movimientos handles GET /api/v1/reportes/movimientos.
func (h *Handler) Movimientos(w http.ResponseWriter, r *http.Request) {
	req, ok := requester(w, r)
	if !ok {
		return
	}
	params, err := h.parseReporteParams(r)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	extra := parseMovimientosFiltros(r)
	ctx, cancel := withTimeout(r)
	defer cancel()

	data, err := h.service.GenerarMovimientos(ctx, req.SedeID, params, extra)
	if err != nil {
		response.WriteError(w, r, mapGenerarErr(err))
		return
	}
	respondXLSX(w, r, filenameRango("movimientos", params), data)
}
