package reportes

import (
	"bytes"
	"context"

	"github.com/shopspring/decimal"

	domainreportes "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/reportes"
	reporterepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
)

// GenerarVentas builds the ventas report (Resumen, Ventas, Items sheets) for
// the resolved range and returns it as XLSX bytes.
func (s *Service) GenerarVentas(ctx context.Context, sedeID int64, params domainreportes.ReporteParams) ([]byte, error) {
	desde, hasta, err := s.resolverRango(params)
	if err != nil {
		return nil, err
	}

	usuarioID := int64(0)
	if params.UsuarioID != nil {
		usuarioID = *params.UsuarioID
	}
	filtro := reporterepo.RangoFiltro{SedeID: sedeID, Desde: desde, Hasta: hasta, UsuarioID: usuarioID}

	resumen, err := s.repo.VentasResumen(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	porVendedora, err := s.repo.VentasPorVendedora(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	detalle, err := s.repo.VentasDetalle(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	items, err := s.repo.VentasItems(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}

	b := NewXLSXBuilder(s.loc)

	// --- Hoja "Resumen" ---
	res := b.NewSheet("Resumen")
	res.WriteTitle(tituloRango("Reporte de ventas", desde, hasta), 3)
	res.SkipRow()

	res.WriteHeaders([]string{"Métrica", "Valor"})
	ticket := decimal.Zero
	if resumen.VentasCount > 0 {
		ticket = resumen.TotalVentas.DivRound(decimal.NewFromInt(resumen.VentasCount), 0)
	}
	res.WriteRow("Total de ventas", resumen.TotalVentas)
	res.WriteRow("Número de ventas", resumen.VentasCount)
	res.WriteRow("Ticket promedio", ticket)
	res.WriteRow("Descuentos aplicados", resumen.DescuentoTotal)
	res.SkipRow()

	res.WriteHeaders([]string{"Método de pago", "Total"})
	res.WriteRow("Efectivo", resumen.TotalEfectivo)
	res.WriteRow("Nequi", resumen.TotalNequi)
	res.WriteRow("Daviplata", resumen.TotalDaviplata)
	res.WriteRow("Transferencia", resumen.TotalTransferencia)
	res.WriteRow("Otros", resumen.TotalOtros)
	sumaMetodos := resumen.TotalEfectivo.Add(resumen.TotalNequi).Add(resumen.TotalDaviplata).
		Add(resumen.TotalTransferencia).Add(resumen.TotalOtros)
	res.WriteTotalRow("Total", sumaMetodos)
	res.SkipRow()

	res.WriteHeaders([]string{"Vendedora", "Ventas", "Total"})
	for _, v := range porVendedora {
		res.WriteRow(v.NombreCompleto, v.VentasCount, v.Total)
	}
	res.AutoWidth()

	// --- Hoja "Ventas" ---
	vh := b.NewSheet("Ventas")
	vh.WriteHeaders([]string{
		"Consecutivo", "Fecha", "Hora", "Vendedora", "Método de pago",
		"Subtotal", "Descuento %", "Descuento", "Total", "Observaciones",
	})
	for _, v := range detalle {
		vh.WriteRow(
			FormatConsecutivoVenta(v.ID),
			Fecha(v.CreatedAt),
			Hora(v.CreatedAt),
			v.UsuarioNombre,
			v.MetodoPagoNombre,
			v.Subtotal,
			Porcentaje(v.DescuentoPct),
			v.DescuentoMonto,
			v.Total,
			v.Observaciones,
		)
	}
	vh.AutoWidth()

	// --- Hoja "Items" ---
	ih := b.NewSheet("Items")
	ih.WriteHeaders([]string{
		"Consecutivo venta", "Fecha", "Tipo línea", "Fragancia", "Envase",
		"Producto", "Feromona", "Gramos", "Cantidad", "Precio unitario", "Subtotal",
	})
	for _, it := range items {
		ih.WriteRow(
			FormatConsecutivoVenta(it.VentaID),
			Fecha(it.CreatedAt),
			it.TipoLinea,
			it.FraganciaNombre,
			emptyOrNil(it.EnvaseNombre),
			it.ProductoNombre,
			it.FeromonaNombre,
			gramosCell(it.Gramos),
			int(it.Cantidad),
			it.PrecioUnitario,
			it.Subtotal,
		)
	}
	ih.AutoWidth()

	b.DeleteDefaultSheet()

	var buf bytes.Buffer
	if err := b.Render(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
