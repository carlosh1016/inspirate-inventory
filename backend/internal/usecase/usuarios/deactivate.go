package usuarios

import (
	"context"
	"errors"

	domainusuarios "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/usuarios"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

// DeactivateInput is the request context: who is deactivating whom.
type DeactivateInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Deactivate disables a usuario's account and revokes all of their
// sessions. Blocked for self-deactivation and for the last active admin.
func (s *Service) Deactivate(ctx context.Context, in DeactivateInput) error {
	target, err := s.Usuarios.GetByID(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}

	var activeAdmins int64
	if string(target.Rol) == "admin" {
		activeAdmins, err = s.Usuarios.CountActiveAdmins(ctx)
		if err != nil {
			return internalErr(err)
		}
	}

	domainUser := domainusuarios.Usuario{ID: target.ID, Rol: string(target.Rol), IsActive: target.IsActive}
	if err := domainUser.CanBeDeactivatedBy(in.RequesterID, activeAdmins); err != nil {
		return err
	}

	if err := s.Usuarios.Deactivate(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}
	if err := s.RefreshTokens.RevokeAllByUser(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "usuario_desactivado", in.IP, in.UserAgent, &in.TargetID, nil, nil)
	return nil
}
