package ventas

import (
	"context"
	"errors"

	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
	ventasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/ventas"
)

// UpdateInput is the request payload plus the requester's context. Ventas
// are immutable except for Observaciones — admin-only, enforced at the
// router.
type UpdateInput struct {
	TargetID      int64
	Observaciones *string
	RequesterID   int64
	IP            string
	UserAgent     string
}

// Update sets a venta's observaciones. Every other field is immutable —
// there's nothing else this usecase can change.
func (s *Service) Update(ctx context.Context, in UpdateInput) (domainventas.Venta, error) {
	before, err := s.Ventas.GetByID(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, ventasrepo.ErrNotFound) {
			return domainventas.Venta{}, notFoundErr()
		}
		return domainventas.Venta{}, internalErr(err)
	}

	updated, err := s.Ventas.UpdateObservaciones(ctx, in.TargetID, in.Observaciones)
	if err != nil {
		if errors.Is(err, ventasrepo.ErrNotFound) {
			return domainventas.Venta{}, notFoundErr()
		}
		return domainventas.Venta{}, internalErr(err)
	}

	var antesObs, despuesObs *string
	if before.Observaciones.Valid {
		antesObs = &before.Observaciones.String
	}
	if updated.Observaciones.Valid {
		despuesObs = &updated.Observaciones.String
	}
	s.audit(ctx, &in.RequesterID, "venta_observaciones_editadas", in.IP, in.UserAgent, &updated.ID,
		auditSnapshot{Observaciones: antesObs}, auditSnapshot{Observaciones: despuesObs})

	venta, err := s.loadVentaCompleta(ctx, updated.ID)
	if err != nil {
		return domainventas.Venta{}, err
	}
	return *venta, nil
}
