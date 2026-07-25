package reportes

import (
	"context"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) VentasResumen(ctx context.Context, f RangoFiltro) (VentasResumen, error) {
	row, err := r.q.ReporteVentasResumen(ctx, generated.ReporteVentasResumenParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Timestamptz(f.Desde),
		FechaHasta: repo.Timestamptz(f.Hasta),
		UsuarioID:  f.UsuarioID,
	})
	if err != nil {
		return VentasResumen{}, err
	}
	return VentasResumen{
		TotalEfectivo:      row.TotalEfectivo,
		TotalNequi:         row.TotalNequi,
		TotalDaviplata:     row.TotalDaviplata,
		TotalTransferencia: row.TotalTransferencia,
		TotalOtros:         row.TotalOtros,
		VentasCount:        row.VentasCount,
		TotalVentas:        row.TotalVentas,
		DescuentoTotal:     row.DescuentoTotal,
	}, nil
}

func (r *postgresRepository) VentasPorVendedora(ctx context.Context, f RangoFiltro) ([]VentasPorVendedora, error) {
	rows, err := r.q.ReporteVentasPorVendedora(ctx, generated.ReporteVentasPorVendedoraParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Timestamptz(f.Desde),
		FechaHasta: repo.Timestamptz(f.Hasta),
		UsuarioID:  f.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]VentasPorVendedora, 0, len(rows))
	for _, row := range rows {
		out = append(out, VentasPorVendedora{
			UsuarioID:      row.UsuarioID,
			NombreCompleto: row.NombreCompleto,
			VentasCount:    row.VentasCount,
			Total:          row.Total,
		})
	}
	return out, nil
}

func (r *postgresRepository) VentasDetalle(ctx context.Context, f RangoFiltro) ([]VentaDetalle, error) {
	rows, err := r.q.ReporteVentasDetalle(ctx, generated.ReporteVentasDetalleParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Timestamptz(f.Desde),
		FechaHasta: repo.Timestamptz(f.Hasta),
		UsuarioID:  f.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]VentaDetalle, 0, len(rows))
	for _, row := range rows {
		out = append(out, VentaDetalle{
			ID:               row.ID,
			CreatedAt:        row.CreatedAt.Time,
			UsuarioNombre:    row.UsuarioNombre,
			MetodoPagoNombre: row.MetodoPagoNombre,
			Subtotal:         row.Subtotal,
			DescuentoPct:     row.DescuentoPct,
			DescuentoMonto:   row.DescuentoMonto,
			Total:            row.Total,
			Observaciones:    repo.StringPtr(row.Observaciones),
		})
	}
	return out, nil
}

func (r *postgresRepository) VentasItems(ctx context.Context, f RangoFiltro) ([]VentaItem, error) {
	rows, err := r.q.ReporteVentasItems(ctx, generated.ReporteVentasItemsParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Timestamptz(f.Desde),
		FechaHasta: repo.Timestamptz(f.Hasta),
		UsuarioID:  f.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]VentaItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, VentaItem{
			VentaID:         row.VentaID,
			CreatedAt:       row.CreatedAt.Time,
			TipoLinea:       row.TipoLinea,
			FraganciaNombre: repo.StringPtr(row.FraganciaNombre),
			EnvaseNombre:    row.EnvaseNombre,
			ProductoNombre:  repo.StringPtr(row.ProductoNombre),
			FeromonaNombre:  repo.StringPtr(row.FeromonaNombre),
			Gramos:          repo.NullDecimalPtr(row.GramosFragancia),
			Cantidad:        row.Cantidad,
			PrecioUnitario:  row.PrecioUnitario,
			Subtotal:        row.Subtotal,
		})
	}
	return out, nil
}

func (r *postgresRepository) StockFragancias(ctx context.Context, f StockFiltro) ([]StockFragancia, error) {
	rows, err := r.q.ReporteStockFragancias(ctx, generated.ReporteStockFraganciasParams{
		SedeID:           f.SedeID,
		IncludeInactivos: f.IncluirInactivos,
	})
	if err != nil {
		return nil, err
	}
	out := make([]StockFragancia, 0, len(rows))
	for _, row := range rows {
		out = append(out, StockFragancia{
			NombreComercial:   row.NombreComercial,
			NombreAlternativo: repo.StringPtr(row.NombreAlternativo),
			Genero:            row.Genero,
			StockVitrina:      row.StockVitrina,
			StockBodega:       row.StockBodega,
			StockTotal:        row.StockTotal,
			Minimo:            row.Minimo,
			BajoMinimo:        row.BajoMinimo,
		})
	}
	return out, nil
}

func (r *postgresRepository) StockEnvases(ctx context.Context, f StockFiltro) ([]StockEnvase, error) {
	rows, err := r.q.ReporteStockEnvases(ctx, generated.ReporteStockEnvasesParams{
		SedeID:           f.SedeID,
		IncludeInactivos: f.IncluirInactivos,
	})
	if err != nil {
		return nil, err
	}
	out := make([]StockEnvase, 0, len(rows))
	for _, row := range rows {
		out = append(out, StockEnvase{
			Tipo:               row.Tipo,
			TamanoOz:           row.TamanoOz,
			Color:              row.Color,
			PrecioSolo:         row.PrecioSolo,
			PrecioConFragancia: row.PrecioConFragancia,
			PrecioRecarga:      row.PrecioRecarga,
			StockVitrina:       row.StockVitrina,
			StockBodega:        row.StockBodega,
			StockTotal:         row.StockTotal,
			Minimo:             row.Minimo,
			BajoMinimo:         row.BajoMinimo,
		})
	}
	return out, nil
}

func (r *postgresRepository) StockProductos(ctx context.Context, f StockFiltro) ([]StockProducto, error) {
	rows, err := r.q.ReporteStockProductos(ctx, generated.ReporteStockProductosParams{
		SedeID:           f.SedeID,
		IncludeInactivos: f.IncluirInactivos,
	})
	if err != nil {
		return nil, err
	}
	out := make([]StockProducto, 0, len(rows))
	for _, row := range rows {
		out = append(out, StockProducto{
			Nombre:       row.Nombre,
			Categoria:    row.Categoria,
			Precio:       row.Precio,
			StockVitrina: row.StockVitrina,
			StockBodega:  row.StockBodega,
			StockTotal:   row.StockTotal,
			Minimo:       row.Minimo,
			BajoMinimo:   row.BajoMinimo,
		})
	}
	return out, nil
}

func (r *postgresRepository) StockAlertas(ctx context.Context, f StockFiltro) ([]StockAlerta, error) {
	rows, err := r.q.ReporteStockAlertas(ctx, generated.ReporteStockAlertasParams{
		SedeID:           f.SedeID,
		IncludeInactivos: f.IncluirInactivos,
	})
	if err != nil {
		return nil, err
	}
	out := make([]StockAlerta, 0, len(rows))
	for _, row := range rows {
		out = append(out, StockAlerta{
			Tipo:        row.Tipo,
			Nombre:      row.Nombre,
			StockActual: row.StockActual,
			Minimo:      row.Minimo,
			Faltante:    row.Faltante,
		})
	}
	return out, nil
}

func (r *postgresRepository) Movimientos(ctx context.Context, f MovimientosFiltro) ([]Movimiento, error) {
	rows, err := r.q.ReporteMovimientos(ctx, generated.ReporteMovimientosParams{
		SedeID:     f.SedeID,
		TipoItem:   f.TipoItem,
		ItemID:     f.ItemID,
		Tipo:       f.Tipo,
		UsuarioID:  f.UsuarioID,
		FechaDesde: repo.Timestamptz(f.Desde),
		FechaHasta: repo.Timestamptz(f.Hasta),
	})
	if err != nil {
		return nil, err
	}
	out := make([]Movimiento, 0, len(rows))
	for _, row := range rows {
		out = append(out, Movimiento{
			CreatedAt:      row.CreatedAt.Time,
			Tipo:           string(row.Tipo),
			TipoItem:       string(row.TipoItem),
			ItemNombre:     row.ItemNombre,
			Ubicacion:      string(row.Ubicacion),
			Cantidad:       row.Cantidad,
			StockAnterior:  row.StockAnterior,
			StockPosterior: row.StockPosterior,
			UsuarioNombre:  row.UsuarioNombre,
			Motivo:         repo.StringPtr(row.Motivo),
			VentaID:        repo.Int8Ptr(row.VentaID),
		})
	}
	return out, nil
}

func (r *postgresRepository) CuadresCerrados(ctx context.Context, f RangoFiltro) ([]CuadreCerrado, error) {
	rows, err := r.q.ReporteCuadresCerrados(ctx, generated.ReporteCuadresCerradosParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Date(f.Desde),
		FechaHasta: repo.Date(f.Hasta),
	})
	if err != nil {
		return nil, err
	}
	out := make([]CuadreCerrado, 0, len(rows))
	for _, row := range rows {
		out = append(out, CuadreCerrado{
			Fecha:               row.Fecha.Time,
			Estado:              row.Estado,
			FondoBase:           row.FondoBase,
			TotalEfectivo:       row.TotalEfectivo,
			TotalNequi:          row.TotalNequi,
			TotalDaviplata:      row.TotalDaviplata,
			TotalTransferencia:  row.TotalTransferencia,
			TotalOtros:          row.TotalOtros,
			TotalPagos:          row.TotalPagos,
			TotalConsignaciones: row.TotalConsignaciones,
			ValorTurno:          row.ValorTurno,
			SaldoCalculado:      row.SaldoCalculado,
			CerradoPor:          repo.StringPtr(row.CerradoPor),
			CerradoAt:           repo.TimePtr(row.CerradoAt),
			Observaciones:       repo.StringPtr(row.Observaciones),
		})
	}
	return out, nil
}

func (r *postgresRepository) CuadresPagos(ctx context.Context, f RangoFiltro) ([]CuadrePago, error) {
	rows, err := r.q.ReporteCuadresPagos(ctx, generated.ReporteCuadresPagosParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Date(f.Desde),
		FechaHasta: repo.Date(f.Hasta),
	})
	if err != nil {
		return nil, err
	}
	out := make([]CuadrePago, 0, len(rows))
	for _, row := range rows {
		out = append(out, CuadrePago{
			FechaCuadre:   row.FechaCuadre.Time,
			Concepto:      row.Concepto,
			Monto:         row.Monto,
			UsuarioNombre: row.UsuarioNombre,
			CreatedAt:     row.CreatedAt.Time,
		})
	}
	return out, nil
}

func (r *postgresRepository) CuadresConsignaciones(ctx context.Context, f RangoFiltro) ([]CuadreConsignacion, error) {
	rows, err := r.q.ReporteCuadresConsignaciones(ctx, generated.ReporteCuadresConsignacionesParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Date(f.Desde),
		FechaHasta: repo.Date(f.Hasta),
	})
	if err != nil {
		return nil, err
	}
	out := make([]CuadreConsignacion, 0, len(rows))
	for _, row := range rows {
		out = append(out, CuadreConsignacion{
			FechaCuadre:   row.FechaCuadre.Time,
			Monto:         row.Monto,
			Banco:         repo.StringPtr(row.Banco),
			Referencia:    repo.StringPtr(row.Referencia),
			UsuarioNombre: row.UsuarioNombre,
			CreatedAt:     row.CreatedAt.Time,
		})
	}
	return out, nil
}

func (r *postgresRepository) SesionesResumen(ctx context.Context, f RangoFiltro) ([]SesionResumen, error) {
	rows, err := r.q.ReporteSesionesResumen(ctx, generated.ReporteSesionesResumenParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Timestamptz(f.Desde),
		FechaHasta: repo.Timestamptz(f.Hasta),
		UsuarioID:  f.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]SesionResumen, 0, len(rows))
	for _, row := range rows {
		out = append(out, SesionResumen{
			UsuarioID:      row.UsuarioID,
			NombreCompleto: row.NombreCompleto,
			TotalHoras:     repo.IntervalToDuration(row.TotalHoras),
			DiasTrabajados: row.DiasTrabajados,
			SesionesCount:  row.SesionesCount,
		})
	}
	return out, nil
}

func (r *postgresRepository) SesionesDetalle(ctx context.Context, f RangoFiltro) ([]SesionDetalle, error) {
	rows, err := r.q.ReporteSesionesDetalle(ctx, generated.ReporteSesionesDetalleParams{
		SedeID:     f.SedeID,
		FechaDesde: repo.Timestamptz(f.Desde),
		FechaHasta: repo.Timestamptz(f.Hasta),
		UsuarioID:  f.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	out := make([]SesionDetalle, 0, len(rows))
	for _, row := range rows {
		out = append(out, SesionDetalle{
			NombreCompleto:  row.NombreCompleto,
			EntradaAt:       row.EntradaAt.Time,
			SalidaAt:        repo.TimePtr(row.SalidaAt),
			HorasTrabajadas: repo.IntervalToDuration(row.HorasTrabajadas),
		})
	}
	return out, nil
}
