// Package conteos is the persistence port for conteos_inventario and
// conteo_inventario_items, consumed by usecase/conteos.
package conteos

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no conteo/item matches the lookup, and also
// when Cerrar's WHERE ... AND estado = 'abierto' guard matches zero rows
// (i.e. the conteo was already closed, or never existed).
var ErrNotFound = errors.New("conteo not found")

// ListFilter narrows and paginates ListPaginated results.
type ListFilter struct {
	Page     int
	PageSize int
	SedeID   int64
}

// Repository is the persistence port for conteos_inventario.
type Repository interface {
	Insert(ctx context.Context, sedeID int64, periodo time.Time, creadoPorUsuarioID int64) (generated.ConteosInventario, error)
	// GenerarItems creates one conteo_inventario_items row per active
	// fragancia of sedeID, resolving saldo_inicial/gramos_sistema in a
	// single INSERT ... SELECT (see db/queries/conteo_inventario.sql).
	GenerarItems(ctx context.Context, conteoID, sedeID int64, periodo time.Time) error
	GetBySedePeriodo(ctx context.Context, sedeID int64, periodo time.Time) (generated.ConteosInventario, error)
	GetByID(ctx context.Context, id int64) (generated.GetConteoByIDRow, error)
	ListItems(ctx context.Context, conteoID int64) ([]generated.ListConteoItemsRow, error)
	UpdateItemFisico(ctx context.Context, id, conteoID int64, gramosFisico decimal.Decimal) (generated.UpdateConteoItemFisicoRow, error)
	Cerrar(ctx context.Context, id, cerradoPorUsuarioID int64) (generated.ConteosInventario, error)
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListConteosPaginatedRow, int64, error)
}
