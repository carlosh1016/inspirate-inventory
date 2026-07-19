package usuarios

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainusuarios "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/usuarios"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

// UpdateInput is the request payload plus the requesting admin's context.
// A nil field means "leave unchanged".
type UpdateInput struct {
	TargetID       int64
	NombreCompleto *string
	Correo         *string
	Rol            *string
	RequesterID    int64
	IP             string
	UserAgent      string
}

// Update applies a partial update to a usuario, enforcing: an admin can't
// change their own rol, demoting the last active admin is blocked, and
// changing correo requires it be free and revokes existing sessions.
func (s *Service) Update(ctx context.Context, in UpdateInput) (generated.Usuario, error) {
	before, err := s.Usuarios.GetByID(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.Usuario{}, notFoundErr()
		}
		return generated.Usuario{}, internalErr(err)
	}

	changingRol := in.Rol != nil && *in.Rol != string(before.Rol)
	if changingRol {
		activeAdmins, err := s.Usuarios.CountActiveAdmins(ctx)
		if err != nil {
			return generated.Usuario{}, internalErr(err)
		}
		domainUser := domainusuarios.Usuario{ID: before.ID, Rol: string(before.Rol), IsActive: before.IsActive}
		if err := domainUser.CanChangeRoleBy(in.RequesterID, *in.Rol, activeAdmins); err != nil {
			return generated.Usuario{}, err
		}
	}

	changingCorreo := in.Correo != nil && *in.Correo != before.Correo
	if changingCorreo {
		exists, err := s.Usuarios.ExistsCorreo(ctx, *in.Correo)
		if err != nil {
			return generated.Usuario{}, internalErr(err)
		}
		if exists {
			return generated.Usuario{}, domainerrors.NewConflict("Correo en uso", "Ya existe un usuario con ese correo.")
		}
	}

	updated, err := s.Usuarios.Update(ctx, in.TargetID, repo.UpdateFields{
		NombreCompleto: in.NombreCompleto,
		Correo:         in.Correo,
		Rol:            in.Rol,
	})
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.Usuario{}, notFoundErr()
		}
		return generated.Usuario{}, internalErr(err)
	}

	if changingCorreo {
		if err := s.RefreshTokens.RevokeAllByUser(ctx, in.TargetID); err != nil {
			return generated.Usuario{}, internalErr(err)
		}
	}

	s.audit(ctx, &in.RequesterID, "usuario_editado", in.IP, in.UserAgent, &updated.ID, snapshot(before), snapshot(updated))

	return updated, nil
}
