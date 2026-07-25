package reportes

import (
	"net/http"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Sesiones handles GET /api/v1/reportes/sesiones-laborales.
func (h *Handler) Sesiones(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.service.GenerarSesiones(ctx, req.SedeID, params)
	if err != nil {
		response.WriteError(w, r, mapGenerarErr(err))
		return
	}
	respondXLSX(w, r, filenameRango("sesiones-laborales", params), data)
}
