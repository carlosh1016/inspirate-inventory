package reportes

import (
	"bytes"
	"context"

	domainreportes "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/reportes"
	reporterepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
)

// GenerarMovimientos builds the movimientos report (single sheet) for the
// resolved range, plus the movimientos-specific filters.
func (s *Service) GenerarMovimientos(ctx context.Context, sedeID int64, params domainreportes.ReporteParams, extra MovimientosFiltros) ([]byte, error) {
	desde, hasta, err := s.resolverRango(params)
	if err != nil {
		return nil, err
	}

	usuarioID := int64(0)
	if params.UsuarioID != nil {
		usuarioID = *params.UsuarioID
	}
	rows, err := s.repo.Movimientos(ctx, reporterepo.MovimientosFiltro{
		SedeID:    sedeID,
		Desde:     desde,
		Hasta:     hasta,
		Tipo:      extra.Tipo,
		TipoItem:  extra.TipoItem,
		ItemID:    extra.ItemID,
		UsuarioID: usuarioID,
	})
	if err != nil {
		return nil, wrapErr(err)
	}
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}

	b := NewXLSXBuilder(s.loc)
	sh := b.NewSheet("Movimientos")
	sh.WriteHeaders([]string{
		"Fecha", "Hora", "Tipo movimiento", "Tipo item", "Item", "Ubicación",
		"Cantidad", "Stock anterior", "Stock posterior", "Usuario", "Motivo", "ID venta",
	})
	for _, m := range rows {
		cantidad, anterior, posterior := cantidadCells(m)
		var ventaCell interface{}
		if m.VentaID != nil {
			ventaCell = FormatConsecutivoVenta(*m.VentaID)
		}
		sh.WriteRow(
			Fecha(m.CreatedAt),
			Hora(m.CreatedAt),
			m.Tipo,
			m.TipoItem,
			m.ItemNombre,
			m.Ubicacion,
			cantidad,
			anterior,
			posterior,
			m.UsuarioNombre,
			m.Motivo,
			ventaCell,
		)
	}
	sh.AutoWidth()

	b.DeleteDefaultSheet()

	var buf bytes.Buffer
	if err := b.Render(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// cantidadCells formats a movimiento's quantity columns: grams (2 decimals) for
// fragancia, unit counts otherwise.
func cantidadCells(m reporterepo.Movimiento) (cantidad, anterior, posterior interface{}) {
	if m.TipoItem == "fragancia" {
		return Gramos(m.Cantidad), Gramos(m.StockAnterior), Gramos(m.StockPosterior)
	}
	return Numero(m.Cantidad), Numero(m.StockAnterior), Numero(m.StockPosterior)
}
