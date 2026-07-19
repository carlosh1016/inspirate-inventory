// Package stockactual is the persistence port for stock_actual needed by
// catalog usecases in this tanda: initializing a new item's stock rows
// (vitrina + bodega, both zero) and reading its total. Applying movements
// against stock is Tanda 3.
package stockactual

import (
	"context"

	"github.com/shopspring/decimal"
)

// Tipo_item values, matching the tipo_item_enum in the schema. Shared here
// so catalog usecases (fragancias, envases, catalog items) don't hardcode
// their own copies of these strings.
const (
	TipoItemFragancia      = "fragancia"
	TipoItemVarianteEnvase = "variante_envase"
	TipoItemProducto       = "producto"
)

// Repository is the persistence port for stock_actual, consumed by
// usecases across the catalog modules (fragancias, envases, catalog items).
type Repository interface {
	// InitializeStock creates the vitrina and bodega rows for a new item
	// with quantity zero. Safe to call more than once (ON CONFLICT DO
	// NOTHING).
	InitializeStock(ctx context.Context, sedeID int64, tipoItem string, itemID int64) error
	// GetStockTotal returns the current vitrina and bodega quantities for
	// an item (zero if no rows exist yet).
	GetStockTotal(ctx context.Context, sedeID int64, tipoItem string, itemID int64) (vitrina, bodega decimal.Decimal, err error)
}
