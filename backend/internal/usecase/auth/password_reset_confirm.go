package auth

import (
	"context"
	"errors"
	"time"

	domainauth "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	passwordresets "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/password_resets"
)

// PasswordResetConfirmInput is the request payload plus audit context.
type PasswordResetConfirmInput struct {
	Token         string
	PasswordNueva string
	IP            string
	UserAgent     string
}

// PasswordResetConfirm consumes a reset token, sets the new password, and
// revokes every existing session for that user.
func (s *Service) PasswordResetConfirm(ctx context.Context, in PasswordResetConfirmInput) error {
	if err := domainauth.ValidatePassword(in.PasswordNueva); err != nil {
		return err
	}

	reset, err := s.PasswordResets.GetByHash(ctx, hashToken(in.Token))
	if err != nil {
		if errors.Is(err, passwordresets.ErrNotFound) {
			return invalidResetTokenErr()
		}
		return internalErr(err)
	}

	domainReset := domainauth.PasswordReset{
		ID:        reset.ID,
		UsuarioID: reset.UsuarioID,
		ExpiresAt: reset.ExpiresAt.Time,
		UsedAt:    repo.TimePtr(reset.UsedAt),
	}

	if domainReset.Used() || domainReset.Expired(time.Now()) {
		return invalidResetTokenErr()
	}

	newHash, err := domainauth.HashPassword(in.PasswordNueva)
	if err != nil {
		return internalErr(err)
	}

	if err := s.Usuarios.UpdatePassword(ctx, reset.UsuarioID, newHash); err != nil {
		return internalErr(err)
	}

	if err := s.PasswordResets.MarkUsed(ctx, reset.ID); err != nil {
		return internalErr(err)
	}

	if err := s.RefreshTokens.RevokeAllByUser(ctx, reset.UsuarioID); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &reset.UsuarioID, "password_reset_completed", in.IP, in.UserAgent)
	return nil
}

func invalidResetTokenErr() error {
	return domainerrors.NewUnauthorized(
		"Enlace inválido",
		"El enlace para restablecer tu contraseña no es válido o ya expiró.",
	)
}
