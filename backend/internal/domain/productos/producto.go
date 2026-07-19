// Package productos contains pure domain logic for productos: no I/O, no
// HTTP, no SQL.
package productos

import (
	"time"

	"github.com/shopspring/decimal"
)

// Categoria is the producto's kind.
type Categoria string

const (
	CategoriaCrema    Categoria = "crema"
	CategoriaFeromona Categoria = "feromona"
	CategoriaHogar    Categoria = "hogar"
	CategoriaRegalo   Categoria = "regalo"
	CategoriaOtro     Categoria = "otro"
)

// Producto is the domain representation of a producto, including its
// current stock snapshot (read alongside the row, never written directly by
// this module — stock changes only happen through movimientos, Tanda 3).
type Producto struct {
	ID           int64
	SedeID       int64
	Nombre       string
	Categoria    Categoria
	Precio       decimal.Decimal
	StockMinimo  int32
	Activo       bool
	StockVitrina decimal.Decimal
	StockBodega  decimal.Decimal
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// StockTotal is vitrina + bodega.
func (p Producto) StockTotal() decimal.Decimal {
	return p.StockVitrina.Add(p.StockBodega)
}

// PuedeEliminarse reports whether the producto has zero stock everywhere —
// the only condition under which it may be soft-deleted.
func (p Producto) PuedeEliminarse() bool {
	return p.StockTotal().IsZero()
}
