package cuadres

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
)

// GetByID loads one cuadre with its pagos/consignaciones. While abierto,
// totals are recomputed live (ventas totals + running pagos/consignaciones
// sums); once cerrado, the persisted (frozen) totals are returned as-is.
func (s *Service) GetByID(ctx context.Context, id int64) (*domaincuadres.Cuadre, error) {
	row, err := s.Cuadres.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, cuadresrepo.ErrNotFound) {
			return nil, notFoundErr()
		}
		return nil, internalErr(err)
	}
	cuadre := toDomainCuadre(cuadresCajaFromGetByIDRow(row))
	attachCerradoPor(&cuadre, row.CerradoPorNombre.String, row.CerradoPorNombre.Valid)
	return s.enrich(ctx, cuadre)
}

// enrich populates Pagos/Consignaciones and, for an abierto cuadre,
// recomputes the live totals + saldo_calculado.
func (s *Service) enrich(ctx context.Context, cuadre domaincuadres.Cuadre) (*domaincuadres.Cuadre, error) {
	pagosRows, err := s.Pagos.GetByCuadre(ctx, cuadre.ID)
	if err != nil {
		return nil, internalErr(err)
	}
	pagos := make([]domaincuadres.PagoCaja, len(pagosRows))
	for i, p := range pagosRows {
		pagos[i] = toDomainPagoCaja(p)
	}
	cuadre.Pagos = pagos

	consigRows, err := s.Consignaciones.GetByCuadre(ctx, cuadre.ID)
	if err != nil {
		return nil, internalErr(err)
	}
	consignaciones := make([]domaincuadres.Consignacion, len(consigRows))
	for i, c := range consigRows {
		consignaciones[i] = toDomainConsignacion(c)
	}
	cuadre.Consignaciones = consignaciones

	if cuadre.Estado != domaincuadres.EstadoAbierto {
		return &cuadre, nil
	}

	totales, err := s.Totales.CalcularParaFecha(ctx, cuadre.SedeID, cuadre.Fecha)
	if err != nil {
		return nil, internalErr(err)
	}

	totalPagos := decimal.Zero
	for _, p := range pagos {
		totalPagos = totalPagos.Add(p.Monto)
	}
	totalConsignaciones := decimal.Zero
	for _, c := range consignaciones {
		totalConsignaciones = totalConsignaciones.Add(c.Monto)
	}

	cuadre.TotalEfectivo = totales.Efectivo
	cuadre.TotalNequi = totales.Nequi
	cuadre.TotalDaviplata = totales.Daviplata
	cuadre.TotalTransferencia = totales.Transferencia
	cuadre.TotalOtros = totales.Otros
	cuadre.TotalPagos = totalPagos
	cuadre.TotalConsignaciones = totalConsignaciones
	cuadre.SaldoCalculado = cuadre.FondoBase.
		Add(totales.Efectivo).
		Sub(totalPagos).
		Sub(totalConsignaciones).
		Sub(cuadre.ValorTurno)

	return &cuadre, nil
}
