// Package movimientos is the persistence port for movimientos_inventario.
// Rows here are immutable — there is no Update or SoftDelete.
package movimientos

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ListFilter narrows and paginates ListPaginated results. Zero values (0,
// "", nil) mean "no filter" for each field.
type ListFilter struct {
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

// InsertParams is the full set of columns for one immutable movimiento row.
type InsertParams struct {
	SedeID         int64
	UsuarioID      int64
	TipoItem       string
	ItemID         int64
	Tipo           string
	Ubicacion      string
	Cantidad       decimal.Decimal
	StockAnterior  decimal.Decimal
	StockPosterior decimal.Decimal
	Motivo         *string
	VentaID        *int64
}

// Repository is the persistence port for movimientos_inventario, consumed
// by usecase/movimientos. NewPostgres accepts generated.DBTX, so Insert can
// run inside InventoryService's transaction.
type Repository interface {
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListMovimientosPaginatedRow, int64, error)
	Insert(ctx context.Context, params InsertParams) (generated.MovimientosInventario, error)
}
