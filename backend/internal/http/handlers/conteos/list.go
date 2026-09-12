package conteos

import (
	"net/http"
	"strconv"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/conteos"
)

// List handles GET /api/v1/conteos-inventario. Admin-only, enforced by the
// router.
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
		Page:     page,
		PageSize: pageSize,
		SedeID:   requester.SedeID,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]ConteoListItemResponse, len(result.Conteos))
	for i, c := range result.Conteos {
		items[i] = toConteoListItemResponse(c)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}
