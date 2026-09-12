package conteos

import (
	"context"
	"errors"

	domainconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/conteos"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	conteosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/conteos"
)

// CerrarInput is the request payload plus the requester's context.
type CerrarInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Cerrar marks a conteo as cerrado. Once cerrado, a conteo is immutable —
// gramos_fisico can no longer be edited, and there is no reopen endpoint.
// Cerrar itself does NOT touch stock_actual: any correction found necessary
// is applied separately via Inventario > Movimientos > Corrección.
func (s *Service) Cerrar(ctx context.Context, in CerrarInput) (*domainconteos.ConteoInventario, error) {
	cerrado, err := s.Conteos.Cerrar(ctx, in.TargetID, in.RequesterID)
	if err != nil {
		if errors.Is(err, conteosrepo.ErrNotFound) {
			return nil, domainerrors.NewConflict("Conteo ya cerrado", "Este conteo de inventario ya fue cerrado.")
		}
		return nil, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "conteo_cerrado", in.IP, in.UserAgent, &cerrado.ID, nil, nil)

	return s.GetByID(ctx, cerrado.ID)
}
