// Package reportes builds downloadable XLSX reports from aggregated data. Each
// Generar* method returns the report as an in-memory byte slice (see the memory
// TODO on very large ranges below). Report generation honors context
// cancellation/timeouts, surfaced to callers as a CodeGatewayTimeout error.
package reportes

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainreportes "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/reportes"
	reporterepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
)

// Service builds the XLSX reports.
type Service struct {
	repo reporterepo.Repository
	loc  *time.Location
}

// NewService wires the report repository and the report timezone (America/Bogota).
func NewService(repo reporterepo.Repository, loc *time.Location) *Service {
	return &Service{repo: repo, loc: loc}
}

// StockParams are the (non-ranged) filters for the stock snapshot report.
type StockParams struct {
	IncluirInactivos bool
	TipoItem         string // "", "fragancia", "variante_envase", "producto"
}

// MovimientosFiltros are the movimientos-specific filters layered on top of the
// resolved date range.
type MovimientosFiltros struct {
	Tipo     string // "" = todos
	TipoItem string // "" = todos
	ItemID   int64  // 0 = todos
}

// resolverRango resolves the report's [desde, hasta] window in the report
// timezone, mapping validation failures to DomainError.
func (s *Service) resolverRango(params domainreportes.ReporteParams) (time.Time, time.Time, error) {
	desde, hasta, err := params.ResolverRango()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return desde, hasta, nil
}

// wrapErr maps context cancellation/deadline errors coming from the DB layer to
// a CodeGatewayTimeout DomainError (HTTP 504); other errors pass through.
func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return domainerrors.NewGatewayTimeout(
			"Reporte demorado",
			"La generación del reporte tomó demasiado tiempo. Intenta con un rango más corto.",
		)
	}
	return err
}

// checkCtx returns a timeout DomainError if ctx is already done.
func checkCtx(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return wrapErr(err)
	}
	return nil
}

// emptyOrNil returns nil for an empty string so the cell is left blank.
func emptyOrNil(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// gramosCell wraps an optional grams value as a Gramos cell, or nil when absent.
func gramosCell(d *decimal.Decimal) interface{} {
	if d == nil {
		return nil
	}
	return Gramos(*d)
}

// tituloRango formats a report title like "prefix - Del 01/07/2026 al 31/07/2026".
func tituloRango(prefix string, desde, hasta time.Time) string {
	return prefix + " - Del " + desde.Format("02/01/2006") + " al " + hasta.Format("02/01/2006")
}
