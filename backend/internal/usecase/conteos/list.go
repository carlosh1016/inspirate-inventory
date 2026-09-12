package conteos

import (
	"context"

	domainconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/conteos"
	conteosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/conteos"
)

// ListInput mirrors conteosrepo.ListFilter.
type ListInput struct {
	Page     int
	PageSize int
	SedeID   int64
}

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// ListResult is a page of conteos (without Items — List is a summary view,
// same convention as cuadres' List) plus pagination metadata.
type ListResult struct {
	Conteos  []domainconteos.ConteoInventario
	Total    int64
	Page     int
	PageSize int
}

func (s *Service) List(ctx context.Context, in ListInput) (*ListResult, error) {
	page := in.Page
	if page < 1 {
		page = defaultPage
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	rows, total, err := s.Conteos.ListPaginated(ctx, conteosrepo.ListFilter{
		Page:     page,
		PageSize: pageSize,
		SedeID:   in.SedeID,
	})
	if err != nil {
		return nil, internalErr(err)
	}

	conteos := make([]domainconteos.ConteoInventario, len(rows))
	for i, r := range rows {
		conteos[i] = conteoFromListRow(r)
	}

	return &ListResult{Conteos: conteos, Total: total, Page: page, PageSize: pageSize}, nil
}
