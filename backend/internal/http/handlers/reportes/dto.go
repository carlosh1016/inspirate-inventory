package reportes

import (
	"net/http"
	"strconv"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainreportes "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/reportes"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/reportes"
)

const fechaLayout = "2006-01-02"

// parseReporteParams reads the common report query params (periodo, fecha,
// fecha_desde, fecha_hasta, usuario_id) interpreting dates in h.loc. Field-level
// requirements per periodo are enforced later by ReporteParams.ResolverRango.
func (h *Handler) parseReporteParams(r *http.Request) (domainreportes.ReporteParams, error) {
	q := r.URL.Query()

	periodo := domainreportes.ReportePeriodo(q.Get("periodo"))
	if periodo == "" {
		periodo = domainreportes.PeriodoDia
	}
	if !periodo.Valido() {
		return domainreportes.ReporteParams{}, domainerrors.NewValidation(
			"Periodo inválido", "El periodo debe ser uno de: dia, semana, mes, rango.",
			map[string][]string{"periodo": {"Valor no reconocido."}},
		)
	}

	fecha, err := h.parseFecha(q.Get("fecha"), "fecha")
	if err != nil {
		return domainreportes.ReporteParams{}, err
	}
	desde, err := h.parseFecha(q.Get("fecha_desde"), "fecha_desde")
	if err != nil {
		return domainreportes.ReporteParams{}, err
	}
	hasta, err := h.parseFecha(q.Get("fecha_hasta"), "fecha_hasta")
	if err != nil {
		return domainreportes.ReporteParams{}, err
	}

	var usuarioID *int64
	if s := q.Get("usuario_id"); s != "" {
		id, convErr := strconv.ParseInt(s, 10, 64)
		if convErr != nil {
			return domainreportes.ReporteParams{}, domainerrors.NewValidation(
				"Solicitud inválida", "usuario_id debe ser numérico.",
				map[string][]string{"usuario_id": {"Debe ser un número."}},
			)
		}
		usuarioID = &id
	}

	return domainreportes.ReporteParams{
		Periodo:    periodo,
		Fecha:      fecha,
		FechaDesde: desde,
		FechaHasta: hasta,
		UsuarioID:  usuarioID,
	}, nil
}

// parseFecha parses an optional YYYY-MM-DD date in the report timezone. An empty
// string yields nil; a malformed one yields a 422.
func (h *Handler) parseFecha(s, field string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.ParseInLocation(fechaLayout, s, h.loc)
	if err != nil {
		return nil, domainerrors.NewValidation(
			"Fecha inválida", "El formato de fecha debe ser AAAA-MM-DD.",
			map[string][]string{field: {"Formato inválido, usa AAAA-MM-DD."}},
		)
	}
	return &t, nil
}

// parseStockParams reads the stock report's query params.
func (h *Handler) parseStockParams(r *http.Request) (usecase.StockParams, error) {
	q := r.URL.Query()
	incluir, _ := strconv.ParseBool(q.Get("incluir_inactivos"))

	tipoItem := q.Get("tipo_item")
	switch tipoItem {
	case "", "fragancia", "variante_envase", "producto":
	default:
		return usecase.StockParams{}, domainerrors.NewValidation(
			"Solicitud inválida", "tipo_item debe ser fragancia, variante_envase o producto.",
			map[string][]string{"tipo_item": {"Valor no reconocido."}},
		)
	}
	return usecase.StockParams{IncluirInactivos: incluir, TipoItem: tipoItem}, nil
}

// parseMovimientosFiltros reads the movimientos-specific filters.
func parseMovimientosFiltros(r *http.Request) usecase.MovimientosFiltros {
	q := r.URL.Query()
	itemID, _ := strconv.ParseInt(q.Get("item_id"), 10, 64)
	return usecase.MovimientosFiltros{
		Tipo:     q.Get("tipo"),
		TipoItem: q.Get("tipo_item"),
		ItemID:   itemID,
	}
}

// filenameRango builds "<tipo>-<desde>-al-<hasta>.xlsx", collapsing to
// "<tipo>-<fecha>.xlsx" when the range is a single day. Falls back to the tipo
// alone if the range cannot be resolved (the generation would have failed first).
func filenameRango(tipo string, params domainreportes.ReporteParams) string {
	desde, hasta, err := params.ResolverRango()
	if err != nil {
		return tipo + ".xlsx"
	}
	d := desde.Format(fechaLayout)
	h := hasta.Format(fechaLayout)
	if d == h {
		return tipo + "-" + d + ".xlsx"
	}
	return tipo + "-" + d + "-al-" + h + ".xlsx"
}
