package reportes

import (
	"bytes"
	"context"

	reporterepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
)

// alertaTipoParaFiltro maps a tipo_item filter value to the label used in the
// alertas sheet's "Tipo" column.
var alertaTipoParaFiltro = map[string]string{
	"fragancia":       "Fragancia",
	"variante_envase": "Envase",
	"producto":        "Producto",
}

// GenerarStock builds the stock snapshot report (Fragancias, Envases,
// Productos, Alertas sheets). It has no date range. TipoItem, when set,
// restricts which sheets carry data (the others keep only their headers).
func (s *Service) GenerarStock(ctx context.Context, sedeID int64, params StockParams) ([]byte, error) {
	filtro := reporterepo.StockFiltro{SedeID: sedeID, IncluirInactivos: params.IncluirInactivos}

	data, err := s.fetchStock(ctx, filtro, params.TipoItem)
	if err != nil {
		return nil, err
	}
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}
	fragancias, envases, productos, alertas := data.fragancias, data.envases, data.productos, data.alertas

	b := NewXLSXBuilder(s.loc)

	// --- Hoja "Fragancias" ---
	fh := b.NewSheet("Fragancias")
	fh.WriteHeaders([]string{
		"Nombre comercial", "Nombre alternativo", "Género",
		"Stock vitrina (g)", "Stock bodega (g)", "Total (g)", "Mínimo (g)", "Bajo mínimo",
	})
	for _, f := range fragancias {
		fh.WriteRow(
			f.NombreComercial, f.NombreAlternativo, f.Genero,
			Gramos(f.StockVitrina), Gramos(f.StockBodega), Gramos(f.StockTotal), Gramos(f.Minimo),
			f.BajoMinimo,
		)
	}
	fh.AutoWidth()

	// --- Hoja "Envases" ---
	eh := b.NewSheet("Envases")
	eh.WriteHeaders([]string{
		"Tipo", "Tamaño (oz)", "Color", "Precio solo", "Precio con fragancia", "Precio recarga",
		"Stock vitrina", "Stock bodega", "Total", "Mínimo", "Bajo mínimo",
	})
	for _, e := range envases {
		eh.WriteRow(
			e.Tipo, Gramos(e.TamanoOz), e.Color,
			e.PrecioSolo, e.PrecioConFragancia, e.PrecioRecarga,
			Numero(e.StockVitrina), Numero(e.StockBodega), Numero(e.StockTotal), Numero(e.Minimo),
			e.BajoMinimo,
		)
	}
	eh.AutoWidth()

	// --- Hoja "Productos" ---
	ph := b.NewSheet("Productos")
	ph.WriteHeaders([]string{
		"Nombre", "Categoría", "Precio",
		"Stock vitrina", "Stock bodega", "Total", "Mínimo", "Bajo mínimo",
	})
	for _, p := range productos {
		ph.WriteRow(
			p.Nombre, p.Categoria, p.Precio,
			Numero(p.StockVitrina), Numero(p.StockBodega), Numero(p.StockTotal), Numero(p.Minimo),
			p.BajoMinimo,
		)
	}
	ph.AutoWidth()

	// --- Hoja "Alertas" ---
	writeAlertasSheet(b, alertas, alertaTipoParaFiltro[params.TipoItem])

	b.DeleteDefaultSheet()

	var buf bytes.Buffer
	if err := b.Render(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// stockDatasets bundles the four stock report datasets.
type stockDatasets struct {
	fragancias []reporterepo.StockFragancia
	envases    []reporterepo.StockEnvase
	productos  []reporterepo.StockProducto
	alertas    []reporterepo.StockAlerta
}

// fetchStock loads the stock datasets, skipping the per-tipo query when a
// tipo_item filter excludes it (that sheet then keeps only its headers).
func (s *Service) fetchStock(ctx context.Context, filtro reporterepo.StockFiltro, tipoItem string) (stockDatasets, error) {
	var d stockDatasets
	var err error
	if tipoItem == "" || tipoItem == "fragancia" {
		if d.fragancias, err = s.repo.StockFragancias(ctx, filtro); err != nil {
			return d, wrapErr(err)
		}
	}
	if tipoItem == "" || tipoItem == "variante_envase" {
		if d.envases, err = s.repo.StockEnvases(ctx, filtro); err != nil {
			return d, wrapErr(err)
		}
	}
	if tipoItem == "" || tipoItem == "producto" {
		if d.productos, err = s.repo.StockProductos(ctx, filtro); err != nil {
			return d, wrapErr(err)
		}
	}
	if d.alertas, err = s.repo.StockAlertas(ctx, filtro); err != nil {
		return d, wrapErr(err)
	}
	return d, nil
}

// writeAlertasSheet builds the "Alertas" sheet. wantTipo (mapped from the
// tipo_item filter) restricts the rows; "" keeps all. Fragancia quantities are
// grams (2 decimals), envase/producto quantities are unit counts.
func writeAlertasSheet(b *XLSXBuilder, alertas []reporterepo.StockAlerta, wantTipo string) {
	ah := b.NewSheet("Alertas")
	ah.WriteHeaders([]string{"Tipo", "Nombre", "Ubicación", "Stock actual", "Mínimo", "Faltante"})
	for _, a := range alertas {
		if wantTipo != "" && a.Tipo != wantTipo {
			continue
		}
		var stock, minimo, faltante interface{}
		if a.Tipo == "Fragancia" {
			stock, minimo, faltante = Gramos(a.StockActual), Gramos(a.Minimo), Gramos(a.Faltante)
		} else {
			stock, minimo, faltante = Numero(a.StockActual), Numero(a.Minimo), Numero(a.Faltante)
		}
		ah.WriteRow(a.Tipo, a.Nombre, "Total", stock, minimo, faltante)
	}
	ah.AutoWidth()
}
