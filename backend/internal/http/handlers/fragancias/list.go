package fragancias

import (
	"net/http"
	"strconv"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/fragancias"
)

// List handles GET /api/v1/fragancias.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	stockBajo, _ := strconv.ParseBool(q.Get("stock_bajo"))
	includeDeleted, _ := strconv.ParseBool(q.Get("include_deleted"))

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:           page,
		PageSize:       pageSize,
		Sort:           q.Get("sort"),
		Q:              q.Get("q"),
		SedeID:         requester.SedeID,
		Genero:         q.Get("genero"),
		Activo:         q.Get("activo"),
		StockBajo:      stockBajo,
		IncludeDeleted: includeDeleted,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]FraganciaResponse, len(result.Items))
	for i, f := range result.Items {
		items[i] = toFraganciaResponseFromList(f)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}
