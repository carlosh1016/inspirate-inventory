package envases

import (
	"time"

	"github.com/shopspring/decimal"
)

// VarianteEnvase is the domain representation of a variante_envase,
// including its current stock snapshot (read alongside the row, never
// written directly by this module — stock changes only happen through
// movimientos, Tanda 3).
type VarianteEnvase struct {
	ID             int64
	ModeloEnvaseID int64
	SedeID         int64
	Color          string
	StockMinimo    int32
	Activo         bool
	StockVitrina   decimal.Decimal
	StockBodega    decimal.Decimal
	DeletedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// StockTotal is vitrina + bodega.
func (v VarianteEnvase) StockTotal() decimal.Decimal {
	return v.StockVitrina.Add(v.StockBodega)
}

// PuedeEliminarse reports whether the variante has zero stock everywhere —
// the only condition under which it may be soft-deleted.
func (v VarianteEnvase) PuedeEliminarse() bool {
	return v.StockTotal().IsZero()
}
