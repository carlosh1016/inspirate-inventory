package stock

import (
	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// StockItemResponse is one row of the unified stock view.
type StockItemResponse struct {
	TipoItem     string `json:"tipo_item"`
	ItemID       int64  `json:"item_id"`
	Nombre       string `json:"nombre"`
	DetalleExtra string `json:"detalle_extra"`
	StockVitrina string `json:"stock_vitrina"`
	StockBodega  string `json:"stock_bodega"`
	StockTotal   string `json:"stock_total"`
	Minimo       string `json:"minimo"`
	BajoMinimo   bool   `json:"bajo_minimo"`
	Unidad       string `json:"unidad"`
}

// toStockItemResponse maps a row to its response shape. When ubicacion is
// "vitrina" or "bodega", the non-requested side is zeroed and total
// recalculated for display only — bajo_minimo keeps reflecting the real
// combined stock, since the filter doesn't redefine the business rule.
func toStockItemResponse(item generated.ListStockUnificadoRow, ubicacion string) StockItemResponse {
	vitrina, bodega, total := item.StockVitrina, item.StockBodega, item.StockTotal

	switch ubicacion {
	case "vitrina":
		bodega = decimal.Zero
		total = vitrina
	case "bodega":
		vitrina = decimal.Zero
		total = bodega
	}

	return StockItemResponse{
		TipoItem:     item.TipoItem,
		ItemID:       item.ItemID,
		Nombre:       item.Nombre,
		DetalleExtra: item.DetalleExtra,
		StockVitrina: vitrina.String(),
		StockBodega:  bodega.String(),
		StockTotal:   total.String(),
		Minimo:       item.Minimo.String(),
		BajoMinimo:   item.BajoMinimo,
		Unidad:       item.Unidad,
	}
}
