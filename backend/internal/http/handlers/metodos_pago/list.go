package metodospago

import (
	"net/http"
	"strconv"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/metodos_pago"
)

// List handles GET /api/v1/metodos-pago.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	includeDeleted, _ := strconv.ParseBool(q.Get("include_deleted"))

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:           page,
		PageSize:       pageSize,
		Sort:           q.Get("sort"),
		Q:              q.Get("q"),
		Activo:         q.Get("activo"),
		IncludeDeleted: includeDeleted,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]MetodoPagoResponse, len(result.Items))
	for i, m := range result.Items {
		items[i] = toMetodoPagoResponse(m)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}
