package movimientos

import (
	"context"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	movimientosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// ListInput is the raw, unvalidated query params for GET /movimientos.
type ListInput struct {
	Page       int
	PageSize   int
	SedeID     int64
	TipoItem   string
	ItemID     int64
	Tipo       string
	UsuarioID  int64
	Ubicacion  string
	VentaID    int64
	FechaDesde *time.Time
	FechaHasta *time.Time
}

// ListResult is a page of movimientos plus pagination metadata.
type ListResult struct {
	Items    []generated.ListMovimientosPaginatedRow
	Total    int64
	Page     int
	PageSize int
}

// List returns a filtered, paginated page of movimientos_inventario,
// newest first.
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

	items, total, err := s.Movimientos.ListPaginated(ctx, movimientosrepo.ListFilter{
		Page:       page,
		PageSize:   pageSize,
		SedeID:     in.SedeID,
		TipoItem:   in.TipoItem,
		ItemID:     in.ItemID,
		Tipo:       in.Tipo,
		UsuarioID:  in.UsuarioID,
		Ubicacion:  in.Ubicacion,
		VentaID:    in.VentaID,
		FechaDesde: in.FechaDesde,
		FechaHasta: in.FechaHasta,
	})
	if err != nil {
		return ListResult{}, internalErr(err)
	}

	return ListResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}
