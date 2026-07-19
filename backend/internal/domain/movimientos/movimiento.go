// Package movimientos contains pure domain logic for movimientos_inventario:
// no I/O, no HTTP, no SQL. The orchestration engine that applies these
// (InventoryService) lives in usecase/movimientos, since it needs pgx.Tx in
// its signature and domain packages stay free of external dependencies.
package movimientos

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
)

// TipoMovimiento identifies the kind of inventory movement, matching the
// tipo_movimiento_enum in the schema.
type TipoMovimiento string

const (
	TipoEntradaMercancia      TipoMovimiento = "entrada_mercancia"
	TipoTrasladoBodegaVitrina TipoMovimiento = "traslado_bodega_vitrina"
	TipoVenta                 TipoMovimiento = "venta"
	TipoAjuste                TipoMovimiento = "ajuste"
	TipoDanado                TipoMovimiento = "danado"
	TipoDevolucion            TipoMovimiento = "devolucion"
	TipoCorreccion            TipoMovimiento = "correccion"
)

// Movimiento is the domain representation of one immutable row in
// movimientos_inventario. ItemNombre is populated by the engine that
// creates it (InventoryService already resolves the item's name to
// validate it exists) so HTTP responses can echo it back without an extra
// query.
type Movimiento struct {
	ID             int64
	SedeID         int64
	UsuarioID      int64
	TipoItem       stock.TipoItem
	ItemID         int64
	ItemNombre     string
	Tipo           TipoMovimiento
	Ubicacion      stock.Ubicacion
	Cantidad       decimal.Decimal
	StockAnterior  decimal.Decimal
	StockPosterior decimal.Decimal
	Motivo         *string
	VentaID        *int64
	CreatedAt      time.Time
}

// MovimientoInput is one movement to register. Cantidad carries an explicit
// sign chosen by the calling usecase (positive for entradas/traslado-in,
// negative for salidas/danado/traslado-out) — InventoryService only checks
// that the resulting stock never goes negative, it doesn't infer signs from
// Tipo.
type MovimientoInput struct {
	TipoItem  stock.TipoItem
	ItemID    int64
	Tipo      TipoMovimiento
	Ubicacion stock.Ubicacion
	Cantidad  decimal.Decimal
	Motivo    *string
	VentaID   *int64
}
