package cuadres

import (
	"context"
	"errors"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
)

// GetHoy loads today's cuadre for sedeID (today = midnight in Colombia).
// Returns (nil, nil) — not an error — when no cuadre exists yet for today;
// the HTTP handler renders that as {"data": null}, meaning "open the day".
func (s *Service) GetHoy(ctx context.Context, sedeID int64) (*domaincuadres.Cuadre, error) {
	hoy := s.hoy()
	row, err := s.Cuadres.GetBySedeFecha(ctx, sedeID, hoy)
	if err != nil {
		if errors.Is(err, cuadresrepo.ErrNotFound) {
			return nil, nil
		}
		return nil, internalErr(err)
	}
	cuadre := toDomainCuadre(cuadresCajaFromGetBySedeFechaRow(row))
	attachCerradoPor(&cuadre, row.CerradoPorNombre.String, row.CerradoPorNombre.Valid)
	return s.enrich(ctx, cuadre)
}
