// Package reportes holds the pure domain types for report generation: the
// period selector and the logic that resolves a period into an effective
// [desde, hasta] date range. It has no external dependencies beyond the
// project's domain/errors package.
package reportes

import (
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
)

// ReportePeriodo selects how the reporting window is derived.
type ReportePeriodo string

const (
	// PeriodoDia covers a single calendar day (uses Fecha).
	PeriodoDia ReportePeriodo = "dia"
	// PeriodoSemana covers Monday..Sunday of the week containing Fecha.
	PeriodoSemana ReportePeriodo = "semana"
	// PeriodoMes covers the calendar month containing Fecha.
	PeriodoMes ReportePeriodo = "mes"
	// PeriodoRango covers an explicit [FechaDesde, FechaHasta] range.
	PeriodoRango ReportePeriodo = "rango"
)

// Valido reports whether p is one of the known periods.
func (p ReportePeriodo) Valido() bool {
	switch p {
	case PeriodoDia, PeriodoSemana, PeriodoMes, PeriodoRango:
		return true
	default:
		return false
	}
}

// ReporteParams is the resolved, validated input to a report. The time.Time
// values must already be expressed in the report timezone (America/Bogota):
// ResolverRango derives week/month boundaries from their Location.
type ReporteParams struct {
	Periodo    ReportePeriodo
	Fecha      *time.Time // dia/semana/mes
	FechaDesde *time.Time // rango
	FechaHasta *time.Time // rango
	UsuarioID  *int64     // optional filter (ventas, sesiones)
}

// ResolverRango interprets Periodo and returns the inclusive [desde, hasta]
// range in the timezone of the provided dates. hasta is the last nanosecond of
// its day, so queries can filter with `created_at >= desde AND created_at <= hasta`.
func (p ReporteParams) ResolverRango() (desde, hasta time.Time, err error) {
	switch p.Periodo {
	case PeriodoDia:
		if p.Fecha == nil {
			return time.Time{}, time.Time{}, errFechaRequerida()
		}
		return startOfDay(*p.Fecha), endOfDay(*p.Fecha), nil

	case PeriodoSemana:
		if p.Fecha == nil {
			return time.Time{}, time.Time{}, errFechaRequerida()
		}
		// Monday-based week: Go's Weekday has Sunday=0..Saturday=6.
		offset := (int(p.Fecha.Weekday()) + 6) % 7
		lunes := p.Fecha.AddDate(0, 0, -offset)
		domingo := lunes.AddDate(0, 0, 6)
		return startOfDay(lunes), endOfDay(domingo), nil

	case PeriodoMes:
		if p.Fecha == nil {
			return time.Time{}, time.Time{}, errFechaRequerida()
		}
		f := *p.Fecha
		primero := time.Date(f.Year(), f.Month(), 1, 0, 0, 0, 0, f.Location())
		ultimo := primero.AddDate(0, 1, -1)
		return startOfDay(primero), endOfDay(ultimo), nil

	case PeriodoRango:
		if p.FechaDesde == nil || p.FechaHasta == nil {
			return time.Time{}, time.Time{}, errRangoRequerido()
		}
		if p.FechaHasta.Before(*p.FechaDesde) {
			return time.Time{}, time.Time{}, errRangoInvalido()
		}
		return startOfDay(*p.FechaDesde), endOfDay(*p.FechaHasta), nil

	default:
		return time.Time{}, time.Time{}, errPeriodoInvalido()
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func errFechaRequerida() *domainerrors.DomainError {
	return domainerrors.NewValidation(
		"Fecha requerida",
		"Debes indicar la fecha del reporte.",
		map[string][]string{"fecha": {"Este campo es obligatorio para el periodo elegido."}},
	)
}

func errRangoRequerido() *domainerrors.DomainError {
	return domainerrors.NewValidation(
		"Rango requerido",
		"Debes indicar la fecha de inicio y de fin del rango.",
		map[string][]string{
			"fecha_desde": {"Este campo es obligatorio para el periodo rango."},
			"fecha_hasta": {"Este campo es obligatorio para el periodo rango."},
		},
	)
}

func errRangoInvalido() *domainerrors.DomainError {
	return domainerrors.NewValidation(
		"Rango inválido",
		"La fecha de fin no puede ser anterior a la fecha de inicio.",
		map[string][]string{"fecha_hasta": {"Debe ser igual o posterior a la fecha de inicio."}},
	)
}

func errPeriodoInvalido() *domainerrors.DomainError {
	return domainerrors.NewValidation(
		"Periodo inválido",
		"El periodo debe ser uno de: dia, semana, mes, rango.",
		map[string][]string{"periodo": {"Valor no reconocido."}},
	)
}
