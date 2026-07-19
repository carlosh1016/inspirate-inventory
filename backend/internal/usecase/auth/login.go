package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainauth "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

// LoginInput is the login request payload plus audit context.
type LoginInput struct {
	Correo    string
	Password  string
	IP        string
	UserAgent string
}

// LoginResult is a fresh Session plus the authenticated usuario, so the
// handler can return both without a second lookup.
type LoginResult struct {
	Session domainauth.Session
	Usuario generated.Usuario
}

// Login authenticates a user and, on success, issues a new Session.
func (s *Service) Login(ctx context.Context, in LoginInput) (LoginResult, error) {
	allowed, retryAfter, err := s.LoginLimiter.Allow(ctx, "login:"+in.IP+":"+in.Correo)
	if err != nil {
		return LoginResult{}, internalErr(err)
	}
	if !allowed {
		return LoginResult{}, domainerrors.NewRateLimit(
			"Demasiados intentos",
			fmt.Sprintf("Demasiados intentos de inicio de sesión. Intenta de nuevo en %s.", retryAfter.Round(time.Second)),
		)
	}

	user, err := s.Usuarios.GetByCorreo(ctx, in.Correo)
	if err != nil {
		if errors.Is(err, usuarios.ErrNotFound) {
			s.audit(ctx, nil, "login_failed", in.IP, in.UserAgent)
			return LoginResult{}, invalidCredentialsErr()
		}
		return LoginResult{}, internalErr(err)
	}

	if !user.IsActive {
		s.audit(ctx, &user.ID, "login_failed", in.IP, in.UserAgent)
		return LoginResult{}, domainerrors.NewForbidden(
			"Cuenta inactiva",
			"Tu cuenta está inactiva. Contacta a un administrador.",
		)
	}

	if !domainauth.CheckPassword(user.PasswordHash, in.Password) {
		s.audit(ctx, &user.ID, "login_failed", in.IP, in.UserAgent)
		return LoginResult{}, invalidCredentialsErr()
	}

	session, err := s.issueSession(ctx, user, in.IP, in.UserAgent)
	if err != nil {
		return LoginResult{}, err
	}

	if err := s.Usuarios.UpdateLastLogin(ctx, user.ID); err != nil {
		return LoginResult{}, internalErr(err)
	}

	s.audit(ctx, &user.ID, "login_success", in.IP, in.UserAgent)

	return LoginResult{Session: session, Usuario: user}, nil
}

func invalidCredentialsErr() error {
	return domainerrors.NewUnauthorized("Credenciales inválidas", "El correo o la contraseña son incorrectos.")
}
