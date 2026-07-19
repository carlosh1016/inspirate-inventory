package usuarios

import (
	"context"
	"errors"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

// ActivateInput is the request context: who is activating whom.
type ActivateInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Activate reactivates a usuario.
func (s *Service) Activate(ctx context.Context, in ActivateInput) error {
	if _, err := s.Usuarios.GetByID(ctx, in.TargetID); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}

	if err := s.Usuarios.Activate(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "usuario_activado", in.IP, in.UserAgent, &in.TargetID, nil, nil)
	return nil
}
