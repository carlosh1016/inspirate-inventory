package productos

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
)

// Get loads a single producto (with its stock snapshot) by id.
func (s *Service) Get(ctx context.Context, id int64) (generated.GetProductoByIDRow, error) {
	p, err := s.Productos.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetProductoByIDRow{}, notFoundErr()
		}
		return generated.GetProductoByIDRow{}, internalErr(err)
	}
	return p, nil
}
