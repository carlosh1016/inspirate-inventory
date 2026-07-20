package ventas

import (
	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/envases"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/productos"
)

// PricingInput carries everything Calculate needs to price one venta_item.
// Which pointer fields are set depends on TipoLinea — the caller (usecase)
// already validated that coherence before building this.
type PricingInput struct {
	TipoLinea        TipoLinea
	ModeloEnvase     *envases.ModeloEnvase
	Producto         *productos.Producto
	FeromonaProducto *productos.Producto
	Cantidad         int32
}

// PricingResult is the priced outcome of one line: PrecioUnitario is what
// goes to venta_items.precio_unitario, SubtotalLinea to venta_items.subtotal.
type PricingResult struct {
	PrecioUnitario decimal.Decimal
	PrecioFeromona decimal.Decimal
	SubtotalLinea  decimal.Decimal
}

// PricingService calculates the unit price and line subtotal for one
// venta_item, given its already-resolved catalog entities. Pure: no I/O.
type PricingService interface {
	Calculate(input PricingInput) (PricingResult, error)
}
