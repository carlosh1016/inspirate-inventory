package ventas

import (
	"context"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ResumenHoyResult is the raw data behind GET /ventas/hoy/resumen — the
// HTTP handler shapes this into the response DTO.
type ResumenHoyResult struct {
	Fecha         string
	Resumen       generated.GetResumenVentasHoyRow
	PorVendedora  []generated.GetVentasPorVendedoraHoyRow
	TopFragancias []generated.GetTopFraganciasHoyRow
}

// ResumenHoy summarizes today's ventas for sedeID, where "today" is
// computed in s.Location (America/Bogota, or its fixed-offset fallback —
// resolved once at startup, see cmd/api/main.go).
func (s *Service) ResumenHoy(ctx context.Context, sedeID int64) (ResumenHoyResult, error) {
	now := time.Now().In(s.Location)
	diaInicio := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.Location)
	diaFin := diaInicio.Add(24 * time.Hour)

	resumen, err := s.Ventas.ResumenHoy(ctx, sedeID, diaInicio, diaFin)
	if err != nil {
		return ResumenHoyResult{}, internalErr(err)
	}
	porVendedora, err := s.Ventas.VentasPorVendedoraHoy(ctx, sedeID, diaInicio, diaFin)
	if err != nil {
		return ResumenHoyResult{}, internalErr(err)
	}
	topFragancias, err := s.Ventas.TopFraganciasHoy(ctx, sedeID, diaInicio, diaFin)
	if err != nil {
		return ResumenHoyResult{}, internalErr(err)
	}

	return ResumenHoyResult{
		Fecha:         diaInicio.Format("2006-01-02"),
		Resumen:       resumen,
		PorVendedora:  porVendedora,
		TopFragancias: topFragancias,
	}, nil
}
