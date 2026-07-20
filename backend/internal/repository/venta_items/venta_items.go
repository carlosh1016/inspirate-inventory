// Package ventaitems is the persistence port for the venta_items table.
// Rows are fully immutable — no Update, no SoftDelete.
package ventaitems

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// InsertParams is the full set of columns for one venta_item. Which of
// FraganciaID/VarianteEnvaseID/ProductoID/FeromonaProductoID/
// GramosFragancia are non-nil is dictated by TipoLinea (chk_venta_items_tipo_linea).
type InsertParams struct {
	VentaID            int64
	TipoLinea          string
	FraganciaID        *int64
	VarianteEnvaseID   *int64
	ProductoID         *int64
	FeromonaProductoID *int64
	GramosFragancia    *decimal.Decimal
	Cantidad           int32
	PrecioUnitario     decimal.Decimal
	Subtotal           decimal.Decimal
}

// Repository is the persistence port for venta_items, consumed by
// usecase/ventas. NewPostgres accepts generated.DBTX, so Insert can run
// inside CreateVenta's transaction alongside the venta it belongs to.
type Repository interface {
	Insert(ctx context.Context, params InsertParams) (generated.VentaItem, error)
	GetByVentaID(ctx context.Context, ventaID int64) ([]generated.GetVentaItemsByVentaIDRow, error)
}
