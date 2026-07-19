package usuarios

import (
	"net/http"
	"strconv"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
)

// List handles GET /api/v1/usuarios.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:     page,
		PageSize: pageSize,
		Sort:     q.Get("sort"),
		Q:        q.Get("q"),
		Rol:      q.Get("rol"),
		Activo:   q.Get("activo"),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]UsuarioResponse, len(result.Items))
	for i, u := range result.Items {
		items[i] = toUsuarioResponse(u)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}
