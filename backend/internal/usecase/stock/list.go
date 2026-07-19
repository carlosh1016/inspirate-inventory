package stock

import (
	"context"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// ListInput is the raw, unvalidated query params for GET /stock (and, with
// StockBajo forced true by the handler, GET /stock/alertas). Ubicacion
// filtering is a display-only transform the HTTP handler applies to the
// result — this usecase always returns the real vitrina/bodega/total.
type ListInput struct {
	Page             int
	PageSize         int
	SedeID           int64
	TipoItem         string // "" | "fragancia" | "variante_envase" | "producto"
	StockBajo        bool
	StockCero        bool
	IncludeInactivos bool
}

// ListResult is a page of unified stock items plus pagination metadata.
type ListResult struct {
	Items    []generated.ListStockUnificadoRow
	Total    int64
	Page     int
	PageSize int
}

// List returns a filtered, paginated page of the unified stock view.
func (s *Service) List(ctx context.Context, in ListInput) (ListResult, error) {
	page := in.Page
	if page < 1 {
		page = defaultPage
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	items, total, err := s.StockActual.ListUnificado(ctx, repo.ListUnificadoFilter{
		Page:             page,
		PageSize:         pageSize,
		SedeID:           in.SedeID,
		TipoItemFilter:   in.TipoItem,
		StockBajo:        in.StockBajo,
		StockCero:        in.StockCero,
		IncludeInactivos: in.IncludeInactivos,
	})
	if err != nil {
		return ListResult{}, internalErr(err)
	}

	return ListResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}
