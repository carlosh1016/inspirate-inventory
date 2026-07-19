package fragancias

import (
	"context"
	"errors"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// Get loads a single fragancia (with its stock snapshot) by id.
func (s *Service) Get(ctx context.Context, id int64) (generated.GetFraganciaByIDRow, error) {
	f, err := s.Fragancias.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetFraganciaByIDRow{}, notFoundErr()
		}
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}
	return f, nil
}
