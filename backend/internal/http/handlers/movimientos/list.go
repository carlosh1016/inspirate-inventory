package movimientos

import (
	"net/http"
	"strconv"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
)

// List handles GET /api/v1/movimientos.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	itemID, _ := strconv.ParseInt(q.Get("item_id"), 10, 64)
	usuarioID, _ := strconv.ParseInt(q.Get("usuario_id"), 10, 64)
	ventaID, _ := strconv.ParseInt(q.Get("venta_id"), 10, 64)

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:       page,
		PageSize:   pageSize,
		SedeID:     requester.SedeID,
		TipoItem:   q.Get("tipo_item"),
		ItemID:     itemID,
		Tipo:       q.Get("tipo"),
		UsuarioID:  usuarioID,
		Ubicacion:  q.Get("ubicacion"),
		VentaID:    ventaID,
		FechaDesde: parseOptionalTime(q.Get("fecha_desde")),
		FechaHasta: parseOptionalTime(q.Get("fecha_hasta")),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]MovimientoResponse, len(result.Items))
	for i, row := range result.Items {
		items[i] = toMovimientoResponseFromList(row)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}

// parseOptionalTime parses an RFC3339 query param, returning nil for an
// empty or malformed value (treated as "no filter" rather than a 422 —
// consistent with how the other list endpoints ignore unparsable filters).
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
