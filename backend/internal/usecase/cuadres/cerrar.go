package cuadres

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// CerrarInput is the request payload plus the requester's context.
type CerrarInput struct {
	TargetID      int64
	ValorTurno    *decimal.Decimal
	Observaciones *string
	RequesterID   int64
	IP            string
	UserAgent     string
}

// auditCerrarSnapshot is what gets JSON-encoded into datos_despues for
// cuadre_cerrado — the full frozen breakdown, for later audit review.
type auditCerrarSnapshot struct {
	TotalEfectivo       string  `json:"total_efectivo"`
	TotalNequi          string  `json:"total_nequi"`
	TotalDaviplata      string  `json:"total_daviplata"`
	TotalTransferencia  string  `json:"total_transferencia"`
	TotalOtros          string  `json:"total_otros"`
	TotalPagos          string  `json:"total_pagos"`
	TotalConsignaciones string  `json:"total_consignaciones"`
	ValorTurno          string  `json:"valor_turno"`
	SaldoCalculado      string  `json:"saldo_calculado"`
	Observaciones       *string `json:"observaciones"`
}

// Cerrar freezes today's ventas totals + pagos/consignaciones sums into the
// cuadre's columns and marks it cerrado. Once cerrado, a cuadre is
// immutable — there is no reopen endpoint.
func (s *Service) Cerrar(ctx context.Context, in CerrarInput) (*domaincuadres.Cuadre, error) {
	row, err := s.Cuadres.GetByID(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, cuadresrepo.ErrNotFound) {
			return nil, notFoundErr()
		}
		return nil, internalErr(err)
	}
	if row.Estado != generated.EstadoCuadreEnumAbierto {
		return nil, domainerrors.NewConflict("Cuadre ya cerrado", "Este cuadre de caja ya fue cerrado.")
	}

	valorTurno := decimal.Zero
	if in.ValorTurno != nil {
		if in.ValorTurno.IsNegative() {
			return nil, domainerrors.NewValidation("Valor de turno inválido", "El valor del turno no puede ser negativo.", nil)
		}
		valorTurno = *in.ValorTurno
	}

	totales, err := s.Totales.CalcularParaFecha(ctx, row.SedeID, row.Fecha.Time)
	if err != nil {
		return nil, internalErr(err)
	}
	totalPagos, err := s.Pagos.GetTotalByCuadre(ctx, row.ID)
	if err != nil {
		return nil, internalErr(err)
	}
	totalConsignaciones, err := s.Consignaciones.GetTotalByCuadre(ctx, row.ID)
	if err != nil {
		return nil, internalErr(err)
	}

	saldoCalculado := row.FondoBase.
		Add(totales.Efectivo).
		Sub(totalPagos).
		Sub(totalConsignaciones).
		Sub(valorTurno)

	cerrado, err := s.Cuadres.Cerrar(ctx, cuadresrepo.CerrarParams{
		ID:                  row.ID,
		TotalEfectivo:       totales.Efectivo,
		TotalNequi:          totales.Nequi,
		TotalDaviplata:      totales.Daviplata,
		TotalTransferencia:  totales.Transferencia,
		TotalOtros:          totales.Otros,
		TotalPagos:          totalPagos,
		TotalConsignaciones: totalConsignaciones,
		ValorTurno:          valorTurno,
		SaldoCalculado:      saldoCalculado,
		Observaciones:       in.Observaciones,
		CerradoPorUsuarioID: in.RequesterID,
	})
	if err != nil {
		if errors.Is(err, cuadresrepo.ErrNotFound) {
			return nil, domainerrors.NewConflict("Cuadre ya cerrado", "Este cuadre de caja ya fue cerrado.")
		}
		return nil, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "cuadre_cerrado", in.IP, in.UserAgent, &cerrado.ID, nil, auditCerrarSnapshot{
		TotalEfectivo:       totales.Efectivo.String(),
		TotalNequi:          totales.Nequi.String(),
		TotalDaviplata:      totales.Daviplata.String(),
		TotalTransferencia:  totales.Transferencia.String(),
		TotalOtros:          totales.Otros.String(),
		TotalPagos:          totalPagos.String(),
		TotalConsignaciones: totalConsignaciones.String(),
		ValorTurno:          valorTurno.String(),
		SaldoCalculado:      saldoCalculado.String(),
		Observaciones:       in.Observaciones,
	})

	return s.GetByID(ctx, cerrado.ID)
}
