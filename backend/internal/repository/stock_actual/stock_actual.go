// Package stockactual is the persistence port for stock_actual: initializing
// a new item's stock rows, reading totals, the unified stock view (M8), and
// the locked read/write pair movimientos (M9) use to apply changes safely
// under concurrency.
package stockactual

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// Tipo_item values, matching the tipo_item_enum in the schema. Shared here
// so catalog usecases (fragancias, envases, catalog items) don't hardcode
// their own copies of these strings.
const (
	TipoItemFragancia      = "fragancia"
	TipoItemVarianteEnvase = "variante_envase"
	TipoItemProducto       = "producto"
)

// Ubicacion values, matching the ubicacion_enum in the schema.
const (
	UbicacionVitrina = "vitrina"
	UbicacionBodega  = "bodega"
)

// ListUnificadoFilter narrows and paginates ListUnificado results. An empty
// TipoItemFilter means "all three types".
type ListUnificadoFilter struct {
	Page             int
	PageSize         int
	SedeID           int64
	TipoItemFilter   string
	StockBajo        bool
	StockCero        bool
	IncludeInactivos bool
}

// Repository is the persistence port for stock_actual, consumed by
// usecases across the catalog modules and by movimientos/stock (Tanda 3).
type Repository interface {
	// InitializeStock creates the vitrina and bodega rows for a new item
	// with quantity zero. Safe to call more than once (ON CONFLICT DO
	// NOTHING).
	InitializeStock(ctx context.Context, sedeID int64, tipoItem string, itemID int64) error
	// GetStockTotal returns the current vitrina and bodega quantities for
	// an item (zero if no rows exist yet).
	GetStockTotal(ctx context.Context, sedeID int64, tipoItem string, itemID int64) (vitrina, bodega decimal.Decimal, err error)
	// ListUnificado returns a filtered, paginated page of the unified stock
	// view across fragancias, variantes_envase and productos.
	ListUnificado(ctx context.Context, filter ListUnificadoFilter) ([]generated.ListStockUnificadoRow, int64, error)
	// GetForUpdate locks and reads the stock_actual row for one
	// (sede, tipo_item, item, ubicacion), for use inside a transaction
	// before applying a movimiento. found is false if no row exists yet
	// (shouldn't normally happen — InitializeStock always creates both
	// rows up front — but is handled defensively).
	GetForUpdate(ctx context.Context, sedeID int64, tipoItem string, itemID int64, ubicacion string) (row generated.StockActual, found bool, err error)
	// SetCantidad writes the new absolute quantity for a (sede, tipo_item,
	// item, ubicacion), creating the row if it's somehow missing.
	SetCantidad(ctx context.Context, sedeID int64, tipoItem string, itemID int64, ubicacion string, cantidad decimal.Decimal) error
}
