package metodospago

import (
	"context"
	"errors"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
)

// DeleteInput is the request context: who is deleting what. Delete is
// admin-only, enforced at the router.
type DeleteInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Delete removes a metodo_pago. If no ventas reference it, it's deleted
// physically (nothing keeps its row meaningful); otherwise ventas.metodo_pago_id
// has an ON DELETE RESTRICT FK, so it's soft-deleted instead — hidden from
// new sales but still resolvable for historical ventas.
func (s *Service) Delete(ctx context.Context, in DeleteInput) error {
	m, err := s.MetodosPago.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}
	if m.DeletedAt.Valid {
		return notFoundErr()
	}

	ventasAsociadas, err := s.MetodosPago.CountVentas(ctx, in.TargetID)
	if err != nil {
		return internalErr(err)
	}

	if ventasAsociadas == 0 {
		if err := s.MetodosPago.HardDelete(ctx, in.TargetID); err != nil {
			return internalErr(err)
		}
		s.audit(ctx, &in.RequesterID, "metodo_pago_eliminado_fisico", in.IP, in.UserAgent, &in.TargetID, snapshot(m), nil)
		return nil
	}

	if err := s.MetodosPago.SoftDelete(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}
	s.audit(ctx, &in.RequesterID, "metodo_pago_eliminado", in.IP, in.UserAgent, &in.TargetID, snapshot(m), nil)
	return nil
}
