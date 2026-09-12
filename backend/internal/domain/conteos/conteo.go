// Package conteos holds the pure entities for the monthly fragrance
// inventory count (conteo mensual de inventario): no I/O, no pgx/sqlc types.
package conteos

import (
	"time"

	"github.com/shopspring/decimal"
)

// EstadoConteo mirrors estado_conteo_enum.
type EstadoConteo string

const (
	EstadoAbierto EstadoConteo = "abierto"
	EstadoCerrado EstadoConteo = "cerrado"
)

// LimiteDiferenciaGramos is the maximum |gramos_fisico - gramos_sistema| a
// fragancia can be off by before it's flagged for review — beyond this,
// either the ventas were logged wrong or the physical count was off.
var LimiteDiferenciaGramos = decimal.NewFromInt(2)

// UsuarioBrief is the minimal (id, nombre_completo) needed to display who
// did something, without pulling in the full usuarios domain type.
type UsuarioBrief struct {
	ID             int64
	NombreCompleto string
}

// ConteoInventario is one month's fragrance count for a sede.
type ConteoInventario struct {
	ID         int64
	SedeID     int64
	Periodo    time.Time // date-only, first day of the month
	Estado     EstadoConteo
	CreadoPor  UsuarioBrief
	CerradoPor *UsuarioBrief
	CerradoAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Populated by GetByID, not by List.
	Items []ConteoItem
}

// ConteoItem is one fragancia's row within a ConteoInventario.
// SaldoInicial is that fragancia's gramos_fisico from the previous closed
// conteo (or its system stock at creation time, if there's no previous
// conteo). GramosSistema is the system's stock_actual total (bodega+vitrina)
// at the moment the conteo was created. GramosFisico is nil until an admin
// enters the physically-counted grams.
type ConteoItem struct {
	ID              int64
	ConteoID        int64
	FraganciaID     int64
	FraganciaNombre string
	SaldoInicial    decimal.Decimal
	GramosSistema   decimal.Decimal
	GramosFisico    *decimal.Decimal
}

// Diferencia is GramosFisico - GramosSistema, or nil while the fragancia
// hasn't been physically counted yet.
func (i ConteoItem) Diferencia() *decimal.Decimal {
	if i.GramosFisico == nil {
		return nil
	}
	d := i.GramosFisico.Sub(i.GramosSistema)
	return &d
}

// ExcedeLimite reports whether |Diferencia| is over LimiteDiferenciaGramos.
// Not yet counted (Diferencia == nil) never exceeds the limit.
func (i ConteoItem) ExcedeLimite() bool {
	d := i.Diferencia()
	if d == nil {
		return false
	}
	return d.Abs().GreaterThan(LimiteDiferenciaGramos)
}

// PorcentajeVendido is how much of SaldoInicial was sold this period,
// (SaldoInicial-GramosSistema)/SaldoInicial*100. Nil when SaldoInicial is
// zero (nothing to compare against — avoids a division by zero).
func (i ConteoItem) PorcentajeVendido() *decimal.Decimal {
	if i.SaldoInicial.IsZero() {
		return nil
	}
	pct := i.SaldoInicial.Sub(i.GramosSistema).Div(i.SaldoInicial).Mul(decimal.NewFromInt(100))
	return &pct
}
