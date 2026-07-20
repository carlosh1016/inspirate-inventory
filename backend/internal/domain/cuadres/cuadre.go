// Package cuadres holds the pure entities for the daily cash register
// close-out (cuadre de caja): no I/O, no pgx/sqlc types.
package cuadres

import (
	"time"

	"github.com/shopspring/decimal"
)

// EstadoCuadre mirrors estado_cuadre_enum.
type EstadoCuadre string

const (
	EstadoAbierto EstadoCuadre = "abierto"
	EstadoCerrado EstadoCuadre = "cerrado"
)

// UsuarioBrief is the minimal (id, nombre_completo) needed to display who
// did something, without pulling in the full usuarios domain type.
type UsuarioBrief struct {
	ID             int64
	NombreCompleto string
}

// Cuadre is one day's cash register close-out for a sede. TotalEfectivo..
// TotalOtros and SaldoCalculado are live-computed while Estado is
// EstadoAbierto (see usecase/cuadres.TotalesService), and frozen at the
// values written by Cerrar once Estado is EstadoCerrado.
type Cuadre struct {
	ID                  int64
	SedeID              int64
	Fecha               time.Time // date-only, midnight in Colombia — no time-of-day meaning
	Estado              EstadoCuadre
	FondoBase           decimal.Decimal
	TotalEfectivo       decimal.Decimal
	TotalNequi          decimal.Decimal
	TotalDaviplata      decimal.Decimal
	TotalTransferencia  decimal.Decimal
	TotalOtros          decimal.Decimal
	TotalPagos          decimal.Decimal
	TotalConsignaciones decimal.Decimal
	ValorTurno          decimal.Decimal
	SaldoCalculado      decimal.Decimal
	Observaciones       *string
	CerradoPorUsuarioID *int64
	CerradoAt           *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time

	// Populated by Get/GetHoy, not by List.
	Pagos          []PagoCaja
	Consignaciones []Consignacion
	CerradoPor     *UsuarioBrief
}

// PagoCaja is one operative cash payout registered against an open Cuadre
// (e.g. buying office supplies out of the till).
type PagoCaja struct {
	ID           int64
	CuadreCajaID int64
	UsuarioID    int64
	Concepto     string
	Monto        decimal.Decimal
	CreatedAt    time.Time
	Usuario      *UsuarioBrief
}

// Consignacion is one bank deposit registered against an open Cuadre.
type Consignacion struct {
	ID           int64
	CuadreCajaID int64
	UsuarioID    int64
	Monto        decimal.Decimal
	Banco        *string
	Referencia   *string
	CreatedAt    time.Time
	Usuario      *UsuarioBrief
}
