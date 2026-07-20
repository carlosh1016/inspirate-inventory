package ventas

// ItemErrorDetalle describes what's wrong with one line of a
// CreateVenta request, by its position in the original request array.
// Used for structural coherence failures, unknown/inactive catalog
// references, and the feromona-wrong-categoria check.
type ItemErrorDetalle struct {
	Index  int    `json:"index"`
	Motivo string `json:"motivo"`
}

// ItemErrorExtra wraps ItemErrorDetalle entries as DomainError.Extra —
// an object with an "items" key, matching the shape the spec defines for
// stock-insuficiente errors.
type ItemErrorExtra struct {
	Items []ItemErrorDetalle `json:"items"`
}

// StockInsuficienteVentaItem is one line that failed a stock-sufficiency
// check while registering a venta's movimientos. Unlike
// domainstock.StockInsuficienteItem (which InventoryService returns and has
// no positional information), this carries Index so the frontend can map
// the shortfall back to the exact line in the request that caused it.
type StockInsuficienteVentaItem struct {
	Index      int    `json:"index"`
	TipoItem   string `json:"tipo_item"`
	ItemID     int64  `json:"item_id"`
	Nombre     string `json:"nombre"`
	Requerido  string `json:"requerido"`
	Disponible string `json:"disponible"`
	Unidad     string `json:"unidad"`
}

// StockInsuficienteExtra wraps StockInsuficienteVentaItem entries as
// DomainError.Extra.
type StockInsuficienteExtra struct {
	Items []StockInsuficienteVentaItem `json:"items"`
}
