package metodospago

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
)

// UpdateInput is the request payload plus the requester's context. A nil
// field means "leave unchanged". Update is admin-only, enforced at the
// router.
type UpdateInput struct {
	TargetID    int64
	Nombre      *string
	Codigo      *string
	RequesterID int64
	IP          string
	UserAgent   string
}

// Update applies a partial update to a metodo_pago, enforcing nombre and
// codigo uniqueness when either changes.
func (s *Service) Update(ctx context.Context, in UpdateInput) (generated.MetodosPago, error) {
	before, err := s.MetodosPago.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.MetodosPago{}, notFoundErr()
		}
		return generated.MetodosPago{}, internalErr(err)
	}
	if before.DeletedAt.Valid {
		return generated.MetodosPago{}, notFoundErr()
	}

	if in.Codigo != nil && *in.Codigo != before.Codigo {
		exists, err := s.MetodosPago.ExistsCodigo(ctx, *in.Codigo, in.TargetID)
		if err != nil {
			return generated.MetodosPago{}, internalErr(err)
		}
		if exists {
			return generated.MetodosPago{}, domainerrors.NewConflict(
				"Código en uso",
				"Ya existe un método de pago con ese código.",
			)
		}
	}
	if in.Nombre != nil && *in.Nombre != before.Nombre {
		exists, err := s.MetodosPago.ExistsNombre(ctx, *in.Nombre, in.TargetID)
		if err != nil {
			return generated.MetodosPago{}, internalErr(err)
		}
		if exists {
			return generated.MetodosPago{}, domainerrors.NewConflict(
				"Nombre en uso",
				"Ya existe un método de pago con ese nombre.",
			)
		}
	}

	updated, err := s.MetodosPago.Update(ctx, in.TargetID, repo.UpdateFields{
		Nombre: in.Nombre,
		Codigo: in.Codigo,
	})
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.MetodosPago{}, notFoundErr()
		}
		return generated.MetodosPago{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "metodo_pago_editado", in.IP, in.UserAgent, &updated.ID, snapshot(before), snapshot(updated))

	return updated, nil
}
