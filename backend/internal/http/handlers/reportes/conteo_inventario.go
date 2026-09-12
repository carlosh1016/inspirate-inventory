package reportes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/reportes"
)

// ConteoInventario handles GET /api/v1/reportes/conteos-inventario/:id.
func (h *Handler) ConteoInventario(w http.ResponseWriter, r *http.Request) {
	req, ok := requester(w, r)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.WriteError(w, r, domainerrors.NewValidation("Solicitud inválida", "El identificador debe ser numérico.", nil))
		return
	}

	ctx, cancel := withTimeout(r)
	defer cancel()

	data, periodo, err := h.service.GenerarConteoInventario(ctx, req.SedeID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrConteoNoEncontrado) {
			response.WriteError(w, r, domainerrors.NewNotFound("Conteo no encontrado", "El conteo de inventario solicitado no existe."))
			return
		}
		response.WriteError(w, r, mapGenerarErr(err))
		return
	}

	// periodo is a calendar date (DATE column, no time-of-day meaning) — format
	// its Y/M/D verbatim, no timezone conversion (would shift the day).
	filename := "conteo-inventario-" + periodo.Format(fechaLayout) + ".xlsx"
	respondXLSX(w, r, filename, data)
}
