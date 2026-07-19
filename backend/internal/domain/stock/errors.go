package stock

import "errors"

// ErrItemNoEncontrado is returned when a movimiento references a catalog
// item that doesn't exist or is soft-deleted.
var ErrItemNoEncontrado = errors.New("item no encontrado o eliminado")

// ErrStockInsuficiente is returned when applying a movimiento would leave
// stock_actual.cantidad below zero for one or more items.
var ErrStockInsuficiente = errors.New("stock insuficiente")

// StockInsuficienteItem describes one item that failed a stock-sufficiency
// check. A slice of these is carried as DomainError.Extra so the frontend
// can map the shortfall to the right form field instead of parsing a
// message string.
type StockInsuficienteItem struct {
	TipoItem   string `json:"tipo_item"`
	ItemID     int64  `json:"item_id"`
	Nombre     string `json:"nombre"`
	Requerido  string `json:"requerido"`
	Disponible string `json:"disponible"`
}
