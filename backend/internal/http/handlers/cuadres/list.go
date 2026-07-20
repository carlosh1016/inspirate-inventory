package cuadres

import (
	"net/http"
	"strconv"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/cuadres"
)

// List handles GET /api/v1/cuadres-caja. Admin-only, enforced by the router.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:       page,
		PageSize:   pageSize,
		SedeID:     requester.SedeID,
		Estado:     q.Get("estado"),
		FechaDesde: parseOptionalDate(q.Get("fecha_desde")),
		FechaHasta: parseOptionalDate(q.Get("fecha_hasta")),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]CuadreListItemResponse, len(result.Cuadres))
	for i, c := range result.Cuadres {
		items[i] = toCuadreListItemResponse(c)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}

// parseOptionalDate parses a "YYYY-MM-DD" query param, returning nil for an
// empty or malformed value (treated as "no filter").
func parseOptionalDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}
