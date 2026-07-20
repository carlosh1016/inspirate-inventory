package cuadres

import (
	"time"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/cuadres"
)

// UsuarioBriefResponse is the (id, nombre_completo) of who did something.
type UsuarioBriefResponse struct {
	ID             int64  `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
}

func toUsuarioBriefResponse(u *domaincuadres.UsuarioBrief) *UsuarioBriefResponse {
	if u == nil {
		return nil
	}
	return &UsuarioBriefResponse{ID: u.ID, NombreCompleto: u.NombreCompleto}
}

// PagoCajaResponse is one pago registered against a cuadre.
type PagoCajaResponse struct {
	ID        int64                 `json:"id"`
	Concepto  string                `json:"concepto"`
	Monto     string                `json:"monto"`
	Usuario   *UsuarioBriefResponse `json:"usuario"`
	CreatedAt time.Time             `json:"created_at"`
}

func toPagoCajaResponse(p domaincuadres.PagoCaja) PagoCajaResponse {
	return PagoCajaResponse{
		ID:        p.ID,
		Concepto:  p.Concepto,
		Monto:     p.Monto.String(),
		Usuario:   toUsuarioBriefResponse(p.Usuario),
		CreatedAt: p.CreatedAt,
	}
}

// ConsignacionResponse is one consignación registered against a cuadre.
type ConsignacionResponse struct {
	ID         int64                 `json:"id"`
	Monto      string                `json:"monto"`
	Banco      *string               `json:"banco"`
	Referencia *string               `json:"referencia"`
	Usuario    *UsuarioBriefResponse `json:"usuario"`
	CreatedAt  time.Time             `json:"created_at"`
}

func toConsignacionResponse(c domaincuadres.Consignacion) ConsignacionResponse {
	return ConsignacionResponse{
		ID:         c.ID,
		Monto:      c.Monto.String(),
		Banco:      c.Banco,
		Referencia: c.Referencia,
		Usuario:    toUsuarioBriefResponse(c.Usuario),
		CreatedAt:  c.CreatedAt,
	}
}

// CuadreResponse is the full response shape for GET (hoy/by-id), POST
// abrir, and POST cerrar.
type CuadreResponse struct {
	ID                  int64                  `json:"id"`
	Fecha               string                 `json:"fecha"`
	Estado              string                 `json:"estado"`
	FondoBase           string                 `json:"fondo_base"`
	TotalEfectivo       string                 `json:"total_efectivo"`
	TotalNequi          string                 `json:"total_nequi"`
	TotalDaviplata      string                 `json:"total_daviplata"`
	TotalTransferencia  string                 `json:"total_transferencia"`
	TotalOtros          string                 `json:"total_otros"`
	TotalVentas         string                 `json:"total_ventas"`
	TotalPagos          string                 `json:"total_pagos"`
	TotalConsignaciones string                 `json:"total_consignaciones"`
	ValorTurno          string                 `json:"valor_turno"`
	SaldoCalculado      string                 `json:"saldo_calculado"`
	Observaciones       *string                `json:"observaciones"`
	Pagos               []PagoCajaResponse     `json:"pagos"`
	Consignaciones      []ConsignacionResponse `json:"consignaciones"`
	CerradoPor          *UsuarioBriefResponse  `json:"cerrado_por"`
	CerradoAt           *time.Time             `json:"cerrado_at"`
	CreatedAt           time.Time              `json:"created_at"`
}

func toCuadreResponse(c domaincuadres.Cuadre) CuadreResponse {
	pagos := make([]PagoCajaResponse, len(c.Pagos))
	for i, p := range c.Pagos {
		pagos[i] = toPagoCajaResponse(p)
	}
	consignaciones := make([]ConsignacionResponse, len(c.Consignaciones))
	for i, cn := range c.Consignaciones {
		consignaciones[i] = toConsignacionResponse(cn)
	}
	totalVentas := c.TotalEfectivo.Add(c.TotalNequi).Add(c.TotalDaviplata).Add(c.TotalTransferencia).Add(c.TotalOtros)

	return CuadreResponse{
		ID:                  c.ID,
		Fecha:               c.Fecha.Format("2006-01-02"),
		Estado:              string(c.Estado),
		FondoBase:           c.FondoBase.String(),
		TotalEfectivo:       c.TotalEfectivo.String(),
		TotalNequi:          c.TotalNequi.String(),
		TotalDaviplata:      c.TotalDaviplata.String(),
		TotalTransferencia:  c.TotalTransferencia.String(),
		TotalOtros:          c.TotalOtros.String(),
		TotalVentas:         totalVentas.String(),
		TotalPagos:          c.TotalPagos.String(),
		TotalConsignaciones: c.TotalConsignaciones.String(),
		ValorTurno:          c.ValorTurno.String(),
		SaldoCalculado:      c.SaldoCalculado.String(),
		Observaciones:       c.Observaciones,
		Pagos:               pagos,
		Consignaciones:      consignaciones,
		CerradoPor:          toUsuarioBriefResponse(c.CerradoPor),
		CerradoAt:           c.CerradoAt,
		CreatedAt:           c.CreatedAt,
	}
}

// CuadreListItemResponse is one row of GET /cuadres-caja — no pagos/
// consignaciones, same convention as ventas' list response.
type CuadreListItemResponse struct {
	ID                  int64                 `json:"id"`
	Fecha               string                `json:"fecha"`
	Estado              string                `json:"estado"`
	FondoBase           string                `json:"fondo_base"`
	TotalVentas         string                `json:"total_ventas"`
	TotalPagos          string                `json:"total_pagos"`
	TotalConsignaciones string                `json:"total_consignaciones"`
	SaldoCalculado      string                `json:"saldo_calculado"`
	CerradoPor          *UsuarioBriefResponse `json:"cerrado_por"`
	CreatedAt           time.Time             `json:"created_at"`
}

func toCuadreListItemResponse(c domaincuadres.Cuadre) CuadreListItemResponse {
	totalVentas := c.TotalEfectivo.Add(c.TotalNequi).Add(c.TotalDaviplata).Add(c.TotalTransferencia).Add(c.TotalOtros)
	return CuadreListItemResponse{
		ID:                  c.ID,
		Fecha:               c.Fecha.Format("2006-01-02"),
		Estado:              string(c.Estado),
		FondoBase:           c.FondoBase.String(),
		TotalVentas:         totalVentas.String(),
		TotalPagos:          c.TotalPagos.String(),
		TotalConsignaciones: c.TotalConsignaciones.String(),
		SaldoCalculado:      c.SaldoCalculado.String(),
		CerradoPor:          toUsuarioBriefResponse(c.CerradoPor),
		CreatedAt:           c.CreatedAt,
	}
}

// WarningResponse is a soft, non-blocking notice.
type WarningResponse struct {
	Codigo  string `json:"codigo"`
	Mensaje string `json:"mensaje"`
}

func toWarningResponses(warnings []usecase.Warning) []WarningResponse {
	out := make([]WarningResponse, len(warnings))
	for i, w := range warnings {
		out[i] = WarningResponse{Codigo: w.Codigo, Mensaje: w.Mensaje}
	}
	return out
}

// AbrirCuadreRequest is the payload for POST /cuadres-caja/abrir.
type AbrirCuadreRequest struct {
	FondoBase *string `json:"fondo_base,omitempty" validate:"omitempty,numeric"`
}

// CerrarCuadreRequest is the payload for POST /cuadres-caja/:id/cerrar.
type CerrarCuadreRequest struct {
	ValorTurno    *string `json:"valor_turno,omitempty" validate:"omitempty,numeric"`
	Observaciones *string `json:"observaciones,omitempty" validate:"omitempty,max=1000"`
}

// CreatePagoRequest is the payload for POST /cuadres-caja/:id/pagos.
type CreatePagoRequest struct {
	Concepto string `json:"concepto" validate:"required,min=2,max=200"`
	Monto    string `json:"monto" validate:"required,numeric"`
}

// CreateConsignacionRequest is the payload for POST
// /cuadres-caja/:id/consignaciones.
type CreateConsignacionRequest struct {
	Monto      string  `json:"monto" validate:"required,numeric"`
	Banco      *string `json:"banco,omitempty" validate:"omitempty,max=100"`
	Referencia *string `json:"referencia,omitempty" validate:"omitempty,max=100"`
}
