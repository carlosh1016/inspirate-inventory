package modelosenvase

import (
	"net/http"
	"strconv"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/modelos_envase"
)

// List handles GET /api/v1/modelos-envase.
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
		Tipo:           q.Get("tipo"),
		Activo:         q.Get("activo"),
		IncludeDeleted: includeDeleted,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]ModeloEnvaseResponse, len(result.Items))
	for i, m := range result.Items {
		items[i] = toModeloEnvaseResponseFromList(m)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}
