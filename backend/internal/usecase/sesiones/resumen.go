package sesiones

import (
	"context"
	"time"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
)

// ResumenInput is the request payload for GET /sesiones-laborales/resumen.
// FechaDesde/FechaHasta are required — the HTTP handler rejects the
// request before calling this if either is missing.
type ResumenInput struct {
	FechaDesde time.Time
	FechaHasta time.Time
	UsuarioID  int64
}

// ResumenItem is one vendedora's aggregated worked time in the range.
// Sesiones still abiertas are excluded (only closed sesiones count towards
// horas_trabajadas).
type ResumenItem struct {
	UsuarioID      int64
	NombreCompleto string
	SesionesCount  int64
	DiasTrabajados int64
	TotalHoras     time.Duration
}

func (s *Service) Resumen(ctx context.Context, in ResumenInput) ([]ResumenItem, error) {
	rows, err := s.Sesiones.GetResumen(ctx, in.FechaDesde, in.FechaHasta, in.UsuarioID)
	if err != nil {
		return nil, internalErr(err)
	}

	items := make([]ResumenItem, len(rows))
	for i, r := range rows {
		var total time.Duration
		if d := repo.IntervalToDuration(r.TotalHoras); d != nil {
			total = *d
		}
		items[i] = ResumenItem{
			UsuarioID:      r.UsuarioID,
			NombreCompleto: r.NombreCompleto,
			SesionesCount:  r.SesionesCount,
			DiasTrabajados: r.DiasTrabajados,
			TotalHoras:     total,
		}
	}
	return items, nil
}
