package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

const passwordResetTTL = time.Hour

// PasswordResetRequestInput is the request payload plus audit context.
type PasswordResetRequestInput struct {
	Correo    string
	IP        string
	UserAgent string
}

// PasswordResetRequest emails a reset link if, and only if, correo belongs
// to an active user — but always succeeds from the caller's point of view,
// so the response never reveals whether an email is registered.
func (s *Service) PasswordResetRequest(ctx context.Context, in PasswordResetRequestInput) error {
	allowed, retryAfter, err := s.ResetLimiter.Allow(ctx, "password-reset:"+in.IP+":"+in.Correo)
	if err != nil {
		return internalErr(err)
	}
	if !allowed {
		return domainerrors.NewRateLimit(
			"Demasiados intentos",
			fmt.Sprintf("Demasiadas solicitudes. Intenta de nuevo en %s.", retryAfter.Round(time.Second)),
		)
	}

	user, err := s.Usuarios.GetByCorreo(ctx, in.Correo)
	if err != nil {
		if errors.Is(err, usuarios.ErrNotFound) {
			return nil
		}
		return internalErr(err)
	}
	if !user.IsActive {
		return nil
	}

	plain, hash, err := generateOpaqueToken()
	if err != nil {
		return internalErr(err)
	}

	if _, err := s.PasswordResets.Insert(ctx, user.ID, hash, time.Now().Add(passwordResetTTL)); err != nil {
		return internalErr(err)
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.FrontendURL, plain)
	if err := s.Mailer.SendPasswordReset(ctx, user.Correo, resetURL, user.NombreCompleto); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &user.ID, "password_reset_requested", in.IP, in.UserAgent)
	return nil
}
