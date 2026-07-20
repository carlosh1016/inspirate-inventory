// Package pagoscaja is the persistence port for pagos_caja, consumed by
// usecase/cuadres.
package pagoscaja

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no pago matches the lookup.
var ErrNotFound = errors.New("pago_caja not found")

// InsertParams is the full set of columns for a new pago.
type InsertParams struct {
	CuadreCajaID int64
	UsuarioID    int64
	Concepto     string
	Monto        decimal.Decimal
}

// Repository is the persistence port for pagos_caja.
type Repository interface {
	Insert(ctx context.Context, params InsertParams) (generated.PagosCaja, error)
	GetByCuadre(ctx context.Context, cuadreCajaID int64) ([]generated.GetPagosByCuadreRow, error)
	GetByID(ctx context.Context, id int64) (generated.PagosCaja, error)
	Delete(ctx context.Context, id int64) error
	GetTotalByCuadre(ctx context.Context, cuadreCajaID int64) (decimal.Decimal, error)
}
