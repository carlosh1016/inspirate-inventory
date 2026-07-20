// Package ventas is the persistence port for the ventas table. Rows are
// immutable except observaciones — there is no SoftDelete.
package ventas

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no venta matches the lookup.
var ErrNotFound = errors.New("venta not found")

// ListFilter narrows and paginates ListPaginated results. Zero values (0,
// nil, false) mean "no filter" for each field.
type ListFilter struct {
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

// InsertParams is the full set of columns for a new venta.
type InsertParams struct {
	SedeID         int64
	UsuarioID      int64
	MetodoPagoID   int64
	Subtotal       decimal.Decimal
	DescuentoPct   decimal.Decimal
	DescuentoMonto decimal.Decimal
	Total          decimal.Decimal
	Observaciones  *string
}

// Repository is the persistence port for ventas, consumed by usecase/ventas.
// NewPostgres accepts generated.DBTX, so Insert can run inside CreateVenta's
// transaction alongside venta_items and the movimientos it generates.
type Repository interface {
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListVentasPaginatedRow, int64, error)
	GetByID(ctx context.Context, id int64) (generated.GetVentaByIDRow, error)
	Insert(ctx context.Context, params InsertParams) (generated.Venta, error)
	UpdateObservaciones(ctx context.Context, id int64, observaciones *string) (generated.Venta, error)
	// ResumenHoy and its two companions back GET /ventas/hoy/resumen.
	// diaInicio/diaFin bound the [inicio, fin) window for "today" in the
	// caller's chosen timezone (America/Bogota) — computed by the usecase,
	// not here.
	ResumenHoy(ctx context.Context, sedeID int64, diaInicio, diaFin time.Time) (generated.GetResumenVentasHoyRow, error)
	VentasPorVendedoraHoy(ctx context.Context, sedeID int64, diaInicio, diaFin time.Time) ([]generated.GetVentasPorVendedoraHoyRow, error)
	TopFraganciasHoy(ctx context.Context, sedeID int64, diaInicio, diaFin time.Time) ([]generated.GetTopFraganciasHoyRow, error)
}
