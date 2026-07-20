package sesiones

import (
	"context"
	"time"

	domainsesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/sesiones"
	sesionesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/sesiones"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// ListInput mirrors sesionesrepo.ListFilter — role-scoping (forcing
// UsuarioID for a vendedora) is the HTTP handler's job, not this usecase's.
type ListInput struct {
	Page       int
	PageSize   int
	SedeID     int64
	UsuarioID  int64
	FechaDesde *time.Time
	FechaHasta *time.Time
	Abiertas   bool
}

// ListResult is a page of sesiones plus pagination metadata.
type ListResult struct {
	Sesiones []domainsesiones.Sesion
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

	rows, total, err := s.Sesiones.ListPaginated(ctx, sesionesrepo.ListFilter{
		Page:       page,
		PageSize:   pageSize,
		SedeID:     in.SedeID,
		UsuarioID:  in.UsuarioID,
		FechaDesde: in.FechaDesde,
		FechaHasta: in.FechaHasta,
		Abiertas:   in.Abiertas,
	})
	if err != nil {
		return nil, internalErr(err)
	}

	sesionesList := make([]domainsesiones.Sesion, len(rows))
	for i, r := range rows {
		sesionesList[i] = sesionFromListRow(r)
	}

	return &ListResult{Sesiones: sesionesList, Total: total, Page: page, PageSize: pageSize}, nil
}
