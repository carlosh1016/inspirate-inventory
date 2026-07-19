// Package envases contains pure domain logic for modelos_envase and
// variantes_envase: no I/O, no HTTP, no SQL.
package envases

import (
	"time"

	"github.com/shopspring/decimal"
)

// ModeloEnvase is the domain representation of a modelo_envase, including
// how many live (non-deleted) child rows currently reference it.
type ModeloEnvase struct {
	ID                 int64
	Tipo               string
	TamanoOz           decimal.Decimal
	EquivGramos        decimal.Decimal
	PrecioSolo         decimal.Decimal
	PrecioConFragancia decimal.Decimal
	PrecioRecarga      decimal.Decimal
	Activo             bool
	VariantesActivas   int64
	DeletedAt          *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// PuedeEliminarse reports whether the modelo has no child rows referencing
// it — the only condition under which it may be soft-deleted.
func (m ModeloEnvase) PuedeEliminarse() bool {
	return m.VariantesActivas == 0
}
