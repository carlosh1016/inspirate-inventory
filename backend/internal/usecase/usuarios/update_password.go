package usuarios

import (
	"context"
	"errors"

	domainauth "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

// UpdatePasswordInput is the request payload plus the requester's context.
// PasswordActual is required (and checked) only when the requester is
// changing their own password; RequesterIsAdmin lets an admin change
// someone else's without it.
type UpdatePasswordInput struct {
	TargetID         int64
	PasswordActual   string
	PasswordNueva    string
	RequesterID      int64
	RequesterIsAdmin bool
	IP               string
	UserAgent        string
}

// UpdatePassword changes a usuario's password and revokes all of their
// existing sessions.
func (s *Service) UpdatePassword(ctx context.Context, in UpdatePasswordInput) error {
	if err := domainauth.ValidatePassword(in.PasswordNueva); err != nil {
		return err
	}

	target, err := s.Usuarios.GetByID(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}

	isSelf := in.RequesterID == in.TargetID
	switch {
	case isSelf:
		if !domainauth.CheckPassword(target.PasswordHash, in.PasswordActual) {
			return domainerrors.NewValidation(
				"Contraseña actual incorrecta",
				"La contraseña actual no es correcta.",
				map[string][]string{"password_actual": {"No es correcta."}},
			)
		}
	case !in.RequesterIsAdmin:
		return domainerrors.NewForbidden("Acceso denegado", "No tienes permiso para realizar esta acción.")
	}

	hash, err := domainauth.HashPassword(in.PasswordNueva)
	if err != nil {
		return internalErr(err)
	}

	if err := s.Usuarios.UpdatePassword(ctx, in.TargetID, hash); err != nil {
		return internalErr(err)
	}

	if err := s.RefreshTokens.RevokeAllByUser(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "password_changed", in.IP, in.UserAgent, &in.TargetID, nil, nil)
	return nil
}
