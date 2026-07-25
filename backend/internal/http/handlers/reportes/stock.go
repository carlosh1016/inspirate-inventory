package reportes

import (
	"net/http"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// Stock handles GET /api/v1/reportes/stock (snapshot actual, sin rango).
func (h *Handler) Stock(w http.ResponseWriter, r *http.Request) {
	req, ok := requester(w, r)
	if !ok {
		return
	}
	params, err := h.parseStockParams(r)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	ctx, cancel := withTimeout(r)
	defer cancel()

	data, err := h.service.GenerarStock(ctx, req.SedeID, params)
	if err != nil {
		response.WriteError(w, r, mapGenerarErr(err))
		return
	}
	filename := "stock-actual-" + time.Now().In(h.loc).Format(fechaLayout) + ".xlsx"
	respondXLSX(w, r, filename, data)
}
