package modelosenvase

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
)

// DeleteInput is the request context: who is deleting what.
type DeleteInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Delete soft-deletes a modelo_envase. Requires no child envase rows
// currently pointing at it.
func (s *Service) Delete(ctx context.Context, in DeleteInput) error {
	m, err := s.ModelosEnvase.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}
	if m.DeletedAt.Valid {
		return notFoundErr()
	}

	variantesActivas, err := s.ModelosEnvase.CountVariantesActivas(ctx, in.TargetID)
	if err != nil {
		return internalErr(err)
	}
	if variantesActivas > 0 {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"No se puede eliminar un modelo de envase con variantes asociadas.",
		)
	}

	if err := s.ModelosEnvase.SoftDelete(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "modelo_envase_eliminado", in.IP, in.UserAgent, &in.TargetID, snapshot(m), nil)
	return nil
}
