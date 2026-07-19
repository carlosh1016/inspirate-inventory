package stock

import (
	"net/http"
	"strconv"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/stock"
)

// List handles GET /api/v1/stock.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.list(w, r, false)
}

// list serves both List and Alertas: forceStockBajo lets Alertas reuse this
// handler with stock_bajo pinned to true, matching the shared-usecase
// pattern described in the plan (same List, different default).
func (h *Handler) list(w http.ResponseWriter, r *http.Request, forceStockBajo bool) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	stockBajo, _ := strconv.ParseBool(q.Get("stock_bajo"))
	stockCero, _ := strconv.ParseBool(q.Get("stock_cero"))
	includeInactivos, _ := strconv.ParseBool(q.Get("include_inactivos"))
	if forceStockBajo {
		stockBajo = true
	}

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:             page,
		PageSize:         pageSize,
		SedeID:           requester.SedeID,
		TipoItem:         q.Get("tipo_item"),
		StockBajo:        stockBajo,
		StockCero:        stockCero,
		IncludeInactivos: includeInactivos,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	ubicacion := q.Get("ubicacion")
	items := make([]StockItemResponse, len(result.Items))
	for i, it := range result.Items {
		items[i] = toStockItemResponse(it, ubicacion)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}
