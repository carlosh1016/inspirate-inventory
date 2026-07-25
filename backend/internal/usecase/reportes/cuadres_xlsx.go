package reportes

import (
	"bytes"
	"context"

	domainreportes "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/reportes"
	reporterepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
)

// GenerarCuadres builds the cuadres de caja report (Cuadres, Pagos,
// Consignaciones sheets) for the resolved range. Only cerrados are included.
func (s *Service) GenerarCuadres(ctx context.Context, sedeID int64, params domainreportes.ReporteParams) ([]byte, error) {
	desde, hasta, err := s.resolverRango(params)
	if err != nil {
		return nil, err
	}
	filtro := reporterepo.RangoFiltro{SedeID: sedeID, Desde: desde, Hasta: hasta}

	cuadres, err := s.repo.CuadresCerrados(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	pagos, err := s.repo.CuadresPagos(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	consignaciones, err := s.repo.CuadresConsignaciones(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}

	b := NewXLSXBuilder(s.loc)

	// --- Hoja "Cuadres" ---
	ch := b.NewSheet("Cuadres")
	ch.WriteHeaders([]string{
		"Fecha", "Estado", "Fondo base", "Total efectivo", "Total nequi", "Total daviplata",
		"Total transferencia", "Total otros", "Total ventas", "Pagos", "Consignaciones",
		"Valor turno", "Saldo", "Cerrado por", "Cerrado el", "Observaciones",
	})
	for _, c := range cuadres {
		totalVentas := c.TotalEfectivo.Add(c.TotalNequi).Add(c.TotalDaviplata).
			Add(c.TotalTransferencia).Add(c.TotalOtros)
		var cerradoEl interface{}
		if c.CerradoAt != nil {
			cerradoEl = *c.CerradoAt
		}
		ch.WriteRow(
			FechaCal(c.Fecha),
			c.Estado,
			c.FondoBase,
			c.TotalEfectivo,
			c.TotalNequi,
			c.TotalDaviplata,
			c.TotalTransferencia,
			c.TotalOtros,
			totalVentas,
			c.TotalPagos,
			c.TotalConsignaciones,
			c.ValorTurno,
			c.SaldoCalculado,
			c.CerradoPor,
			cerradoEl,
			c.Observaciones,
		)
	}
	ch.AutoWidth()

	// --- Hoja "Pagos" ---
	ph := b.NewSheet("Pagos")
	ph.WriteHeaders([]string{"Fecha cuadre", "Concepto", "Monto", "Usuario", "Fecha registro"})
	for _, p := range pagos {
		ph.WriteRow(FechaCal(p.FechaCuadre), p.Concepto, p.Monto, p.UsuarioNombre, p.CreatedAt)
	}
	ph.AutoWidth()

	// --- Hoja "Consignaciones" ---
	csh := b.NewSheet("Consignaciones")
	csh.WriteHeaders([]string{"Fecha cuadre", "Monto", "Banco", "Referencia", "Usuario", "Fecha registro"})
	for _, c := range consignaciones {
		csh.WriteRow(FechaCal(c.FechaCuadre), c.Monto, c.Banco, c.Referencia, c.UsuarioNombre, c.CreatedAt)
	}
	csh.AutoWidth()

	b.DeleteDefaultSheet()

	var buf bytes.Buffer
	if err := b.Render(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
