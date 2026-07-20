package sesiones

import (
	"net/http"
	"strconv"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
)

// Resumen handles GET /api/v1/sesiones-laborales/resumen. Admin-only,
// enforced by the router. fecha_desde/fecha_hasta are required query
// params (RFC3339) — usuario_id is optional.
func (h *Handler) Resumen(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fechaDesde := parseOptionalTime(q.Get("fecha_desde"))
	fechaHasta := parseOptionalTime(q.Get("fecha_hasta"))
	if fechaDesde == nil || fechaHasta == nil {
		response.WriteError(w, r, domainerrors.NewValidation(
			"Solicitud inválida", "Debes indicar fecha_desde y fecha_hasta.", nil,
		))
		return
	}
	usuarioID, _ := strconv.ParseInt(q.Get("usuario_id"), 10, 64)

	items, err := h.service.Resumen(r.Context(), usecase.ResumenInput{
		FechaDesde: *fechaDesde,
		FechaHasta: *fechaHasta,
		UsuarioID:  usuarioID,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	out := make([]ResumenSesionResponse, len(items))
	for i, item := range items {
		out[i] = toResumenSesionResponse(item)
	}

	response.WriteData(w, http.StatusOK, out)
}
