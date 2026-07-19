package metodospago

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
)

// Get loads a single metodo_pago by id.
func (s *Service) Get(ctx context.Context, id int64) (generated.MetodosPago, error) {
	m, err := s.MetodosPago.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.MetodosPago{}, notFoundErr()
		}
		return generated.MetodosPago{}, internalErr(err)
	}
	return m, nil
}
