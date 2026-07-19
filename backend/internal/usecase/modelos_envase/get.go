package modelosenvase

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
)

// Get loads a single modelo_envase, including how many child rows currently
// reference it, by id.
func (s *Service) Get(ctx context.Context, id int64) (generated.GetModeloEnvaseByIDRow, error) {
	m, err := s.ModelosEnvase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetModeloEnvaseByIDRow{}, notFoundErr()
		}
		return generated.GetModeloEnvaseByIDRow{}, internalErr(err)
	}
	return m, nil
}
