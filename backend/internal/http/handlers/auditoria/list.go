package auditoria

import (
	"net/http"
	"strconv"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auditoria"
)

// List handles GET /api/v1/auditoria. Admin-only (enforced by the router).
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	usuarioID, _ := strconv.ParseInt(q.Get("usuario_id"), 10, 64)

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:          page,
		PageSize:      pageSize,
		UsuarioID:     usuarioID,
		Accion:        q.Get("accion"),
		TablaAfectada: q.Get("tabla_afectada"),
		FechaDesde:    parseOptionalTime(q.Get("fecha_desde")),
		FechaHasta:    parseOptionalTime(q.Get("fecha_hasta")),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteList(w, http.StatusOK, toEventoResponses(result.Eventos), result.Total, result.Page, result.PageSize)
}

func parseOptionalTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}
