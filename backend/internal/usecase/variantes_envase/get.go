package variantesenvase

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

// Get loads a single variante_envase (with its stock snapshot) by id.
func (s *Service) Get(ctx context.Context, id int64) (generated.GetVarianteEnvaseByIDRow, error) {
	v, err := s.VariantesEnvase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetVarianteEnvaseByIDRow{}, notFoundErr()
		}
		return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
	}
	return v, nil
}
