package reportes

import (
	"net/http"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Cuadres handles GET /api/v1/reportes/cuadres-caja (solo cuadres cerrados).
func (h *Handler) Cuadres(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.service.GenerarCuadres(ctx, req.SedeID, params)
	if err != nil {
		response.WriteError(w, r, mapGenerarErr(err))
		return
	}
	respondXLSX(w, r, filenameRango("cuadres-caja", params), data)
}
