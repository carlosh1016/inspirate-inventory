// Package stock contains pure domain logic for the unified stock view and
// the movement engine: no I/O, no HTTP, no SQL.
package stock

import "github.com/shopspring/decimal"

// TipoItem identifies which catalog table an item belongs to, matching the
// tipo_item_enum in the schema.
type TipoItem string

const (
	TipoItemFragancia      TipoItem = "fragancia"
	TipoItemVarianteEnvase TipoItem = "variante_envase"
	TipoItemProducto       TipoItem = "producto"
)

// Ubicacion identifies where a quantity of stock physically sits, matching
// the ubicacion_enum in the schema.
type Ubicacion string

const (
	UbicacionVitrina Ubicacion = "vitrina"
	UbicacionBodega  Ubicacion = "bodega"
)

// StockItem is one row of the unified stock view: a catalog item (fragancia,
// variante_envase or producto) alongside its current stock snapshot.
type StockItem struct {
	TipoItem     TipoItem
	ItemID       int64
	Nombre       string
	DetalleExtra string
	StockVitrina decimal.Decimal
	StockBodega  decimal.Decimal
	StockTotal   decimal.Decimal
	Minimo       decimal.Decimal
	BajoMinimo   bool
	Unidad       string
}
