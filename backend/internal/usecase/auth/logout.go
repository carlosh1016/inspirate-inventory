package auth

import (
	"context"
	"errors"

	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
)

// LogoutInput is the plaintext refresh token read from the cookie (if any)
// plus audit context.
type LogoutInput struct {
	RefreshToken string
	IP           string
	UserAgent    string
	UsuarioID    *int64 // from the access token, when present
}

// Logout revokes the given refresh token, if any. It's idempotent: calling
// it with a missing or already-revoked token still succeeds.
func (s *Service) Logout(ctx context.Context, in LogoutInput) error {
	if in.RefreshToken != "" {
		token, err := s.RefreshTokens.GetByHash(ctx, hashToken(in.RefreshToken))
		switch {
		case err == nil:
			if revokeErr := s.RefreshTokens.Revoke(ctx, token.ID); revokeErr != nil {
				return internalErr(revokeErr)
			}
		case errors.Is(err, refreshtokens.ErrNotFound):
			// already gone; nothing to revoke.
		default:
			return internalErr(err)
		}
	}

	s.audit(ctx, in.UsuarioID, "logout", in.IP, in.UserAgent)
	return nil
}
