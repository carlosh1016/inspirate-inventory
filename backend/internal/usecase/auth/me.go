package auth

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

// Me loads the current usuario row for userID. The Auth middleware already
// validated the JWT; this re-fetches in case something changed since the
// token was issued (e.g. the account was deactivated).
func (s *Service) Me(ctx context.Context, userID int64) (generated.Usuario, error) {
	user, err := s.Usuarios.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, usuarios.ErrNotFound) {
			return generated.Usuario{}, invalidSessionErr()
		}
		return generated.Usuario{}, internalErr(err)
	}
	if !user.IsActive {
		return generated.Usuario{}, invalidSessionErr()
	}
	return user, nil
}
