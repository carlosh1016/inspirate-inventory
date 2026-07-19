// Package fragancias contains pure domain logic for fragancias: no I/O, no
// HTTP, no SQL.
package fragancias

import (
	"time"

	"github.com/shopspring/decimal"
)

// Genero is the fragancia's target gender.
type Genero string

const (
	GeneroMasculina Genero = "masculina"
	GeneroFemenina  Genero = "femenina"
)

// Fragancia is the domain representation of a fragancia, including its
// current stock snapshot (read alongside the row, never written directly by
// this module — stock changes only happen through movimientos, Tanda 3).
type Fragancia struct {
	ID                int64
	SedeID            int64
	NombreComercial   string
	NombreAlternativo *string
	Genero            Genero
	GramosMinimo      decimal.Decimal
	Activo            bool
	StockVitrina      decimal.Decimal
	StockBodega       decimal.Decimal
	DeletedAt         *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// StockTotal is vitrina + bodega.
func (f Fragancia) StockTotal() decimal.Decimal {
	return f.StockVitrina.Add(f.StockBodega)
}

// PuedeEliminarse reports whether the fragancia has zero stock everywhere —
// the only condition under which it may be soft-deleted.
func (f Fragancia) PuedeEliminarse() bool {
	return f.StockTotal().IsZero()
}
