package sesiones

import (
	"net/http"
	"strconv"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
)

// List handles GET /api/v1/sesiones-laborales. A vendedora only sees her
// own sesiones — usuario_id is forced to the requester's own id regardless
// of what the query string asks for.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	usuarioID, _ := strconv.ParseInt(q.Get("usuario_id"), 10, 64)
	if requester.Rol == RolVendedora {
		usuarioID = requester.ID
	}
	abiertas, _ := strconv.ParseBool(q.Get("abiertas"))

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:       page,
		PageSize:   pageSize,
		SedeID:     requester.SedeID,
		UsuarioID:  usuarioID,
		FechaDesde: parseOptionalTime(q.Get("fecha_desde")),
		FechaHasta: parseOptionalTime(q.Get("fecha_hasta")),
		Abiertas:   abiertas,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]SesionResponse, len(result.Sesiones))
	for i, s := range result.Sesiones {
		items[i] = toSesionResponse(s)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}

// parseOptionalTime parses an RFC3339 query param, returning nil for an
// empty or malformed value (treated as "no filter", consistent with
// /ventas' and /movimientos' identical helper).
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
