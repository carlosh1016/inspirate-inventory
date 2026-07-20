// Package consignaciones is the persistence port for consignaciones,
// consumed by usecase/cuadres.
package consignaciones

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no consignacion matches the lookup.
var ErrNotFound = errors.New("consignacion not found")

// InsertParams is the full set of columns for a new consignacion.
type InsertParams struct {
	CuadreCajaID int64
	UsuarioID    int64
	Monto        decimal.Decimal
	Banco        *string
	Referencia   *string
}

// Repository is the persistence port for consignaciones.
type Repository interface {
	Insert(ctx context.Context, params InsertParams) (generated.Consignacione, error)
	GetByCuadre(ctx context.Context, cuadreCajaID int64) ([]generated.GetConsignacionesByCuadreRow, error)
	GetByID(ctx context.Context, id int64) (generated.Consignacione, error)
	Delete(ctx context.Context, id int64) error
	GetTotalByCuadre(ctx context.Context, cuadreCajaID int64) (decimal.Decimal, error)
}
