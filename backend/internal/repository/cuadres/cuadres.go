// Package cuadres is the persistence port for cuadres_caja, consumed by
// usecase/cuadres and usecase/ventas (via CajaStatusService, read-only).
package cuadres

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no cuadre matches the lookup, and also when
// Cerrar's WHERE ... AND estado = 'abierto' guard matches zero rows (i.e.
// the cuadre was already closed, or never existed).
var ErrNotFound = errors.New("cuadre not found")

// ListFilter narrows and paginates ListPaginated results. Zero values (0,
// "", nil) mean "no filter" for each field.
type ListFilter struct {
	Page       int
	PageSize   int
	SedeID     int64
	Estado     string
	FechaDesde *time.Time
	FechaHasta *time.Time
}

// InsertParams is the full set of columns for a new (always abierto) cuadre.
type InsertParams struct {
	SedeID    int64
	Fecha     time.Time
	FondoBase decimal.Decimal
}

// CerrarParams is every value UpdateCuadreCerrar needs — all totals are
// pre-computed by usecase/cuadres before this call, never derived here.
type CerrarParams struct {
	ID                  int64
	TotalEfectivo       decimal.Decimal
	TotalNequi          decimal.Decimal
	TotalDaviplata      decimal.Decimal
	TotalTransferencia  decimal.Decimal
	TotalOtros          decimal.Decimal
	TotalPagos          decimal.Decimal
	TotalConsignaciones decimal.Decimal
	ValorTurno          decimal.Decimal
	SaldoCalculado      decimal.Decimal
	Observaciones       *string
	CerradoPorUsuarioID int64
}

// Repository is the persistence port for cuadres_caja.
type Repository interface {
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListCuadresPaginatedRow, int64, error)
	GetByID(ctx context.Context, id int64) (generated.GetCuadreByIDRow, error)
	GetBySedeFecha(ctx context.Context, sedeID int64, fecha time.Time) (generated.GetCuadreBySedeFechaRow, error)
	GetAbiertoAnterior(ctx context.Context, sedeID int64, fecha time.Time) (generated.CuadresCaja, error)
	Insert(ctx context.Context, params InsertParams) (generated.CuadresCaja, error)
	Cerrar(ctx context.Context, params CerrarParams) (generated.CuadresCaja, error)
	ExistsCerradoBySedeFecha(ctx context.Context, sedeID int64, fecha time.Time) (bool, error)
	TotalesPorMetodoEnFecha(ctx context.Context, sedeID int64, diaInicio, diaFin time.Time) (generated.GetTotalesPorMetodoEnFechaRow, error)
}
