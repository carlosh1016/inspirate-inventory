package ventas

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	ventasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/ventas"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// ListInput is the raw, unvalidated query params for GET /ventas. SedeID
// always comes from the requester's own claims. UsuarioID is forced by the
// HTTP handler to the requester's own id when the requester is a vendedora
// — this usecase doesn't know about roles, it just applies whatever filter
// it's given.
type ListInput struct {
	Page         int
	PageSize     int
	SedeID       int64
	UsuarioID    int64
	MetodoPagoID int64
	FechaDesde   *time.Time
	FechaHasta   *time.Time
	TotalMin     decimal.Decimal
	TotalMax     decimal.Decimal
	ConDescuento bool
}

// ListResult is a page of ventas plus pagination metadata.
type ListResult struct {
	Items    []generated.ListVentasPaginatedRow
	Total    int64
	Page     int
	PageSize int
}

// List returns a filtered, paginated page of ventas, newest first.
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

	items, total, err := s.Ventas.ListPaginated(ctx, ventasrepo.ListFilter{
		Page:         page,
		PageSize:     pageSize,
		SedeID:       in.SedeID,
		UsuarioID:    in.UsuarioID,
		MetodoPagoID: in.MetodoPagoID,
		FechaDesde:   in.FechaDesde,
		FechaHasta:   in.FechaHasta,
		TotalMin:     in.TotalMin,
		TotalMax:     in.TotalMax,
		ConDescuento: in.ConDescuento,
	})
	if err != nil {
		return ListResult{}, internalErr(err)
	}

	return ListResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}
