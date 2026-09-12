package conteos

import (
	"context"
	"errors"

	domainconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/conteos"
	conteosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/conteos"
)

// GetByID loads one conteo with its fragancia rows.
func (s *Service) GetByID(ctx context.Context, id int64) (*domainconteos.ConteoInventario, error) {
	row, err := s.Conteos.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, conteosrepo.ErrNotFound) {
			return nil, notFoundErr()
		}
		return nil, internalErr(err)
	}
	conteo := conteoFromGetByIDRow(row)

	itemRows, err := s.Conteos.ListItems(ctx, id)
	if err != nil {
		return nil, internalErr(err)
	}
	items := make([]domainconteos.ConteoItem, len(itemRows))
	for i, r := range itemRows {
		items[i] = toDomainConteoItem(r)
	}
	conteo.Items = items

	return &conteo, nil
}
