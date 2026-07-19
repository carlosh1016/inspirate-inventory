package usuarios

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

// Get loads a single usuario by id.
func (s *Service) Get(ctx context.Context, id int64) (generated.Usuario, error) {
	user, err := s.Usuarios.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.Usuario{}, notFoundErr()
		}
		return generated.Usuario{}, internalErr(err)
	}
	return user, nil
}
