package auth

import (
	"context"
	"errors"
	"time"

	domainauth "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
)

// RefreshInput is the plaintext refresh token read from the cookie plus
// audit context.
type RefreshInput struct {
	RefreshToken string
	IP           string
	UserAgent    string
}

// Refresh rotates a valid refresh token into a new Session. Reuse of an
// already-revoked token revokes every session for that user (the session is
// assumed compromised) and returns 401.
func (s *Service) Refresh(ctx context.Context, in RefreshInput) (domainauth.Session, error) {
	if in.RefreshToken == "" {
		return domainauth.Session{}, invalidSessionErr()
	}

	token, err := s.RefreshTokens.GetByHash(ctx, hashToken(in.RefreshToken))
	if err != nil {
		if errors.Is(err, refreshtokens.ErrNotFound) {
			return domainauth.Session{}, invalidSessionErr()
		}
		return domainauth.Session{}, internalErr(err)
	}

	domainToken := domainauth.RefreshToken{
		ID:        token.ID,
		UsuarioID: token.UsuarioID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt.Time,
		RevokedAt: repo.TimePtr(token.RevokedAt),
	}

	if domainToken.Revoked() {
		if err := s.RefreshTokens.RevokeAllByUser(ctx, token.UsuarioID); err != nil {
			return domainauth.Session{}, internalErr(err)
		}
		return domainauth.Session{}, invalidSessionErr()
	}

	if domainToken.Expired(time.Now()) {
		return domainauth.Session{}, invalidSessionErr()
	}

	user, err := s.Usuarios.GetByID(ctx, token.UsuarioID)
	if err != nil {
		return domainauth.Session{}, internalErr(err)
	}
	if !user.IsActive {
		return domainauth.Session{}, invalidSessionErr()
	}

	if err := s.RefreshTokens.Revoke(ctx, token.ID); err != nil {
		return domainauth.Session{}, internalErr(err)
	}

	return s.issueSession(ctx, user, in.IP, in.UserAgent)
}

func invalidSessionErr() error {
	return domainerrors.NewUnauthorized("Sesión inválida", "Tu sesión expiró o no es válida. Inicia sesión de nuevo.")
}
