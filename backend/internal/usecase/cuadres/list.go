package cuadres

import (
	"context"
	"time"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
)

// ListInput mirrors cuadresrepo.ListFilter — role-scoping (if any is ever
// needed) is the HTTP handler's job, not this usecase's.
type ListInput struct {
	Page       int
	PageSize   int
	SedeID     int64
	Estado     string
	FechaDesde *time.Time
	FechaHasta *time.Time
}

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// ListResult is a page of cuadres (without Pagos/Consignaciones — List is a
// summary view, same convention as ventas' List) plus pagination metadata.
type ListResult struct {
	Cuadres  []domaincuadres.Cuadre
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

	rows, total, err := s.Cuadres.ListPaginated(ctx, cuadresrepo.ListFilter{
		Page:       page,
		PageSize:   pageSize,
		SedeID:     in.SedeID,
		Estado:     in.Estado,
		FechaDesde: in.FechaDesde,
		FechaHasta: in.FechaHasta,
	})
	if err != nil {
		return nil, internalErr(err)
	}

	cuadres := make([]domaincuadres.Cuadre, len(rows))
	for i, r := range rows {
		c := toDomainCuadre(cuadresCajaFromListRow(r))
		attachCerradoPor(&c, r.CerradoPorNombre.String, r.CerradoPorNombre.Valid)
		cuadres[i] = c
	}

	return &ListResult{Cuadres: cuadres, Total: total, Page: page, PageSize: pageSize}, nil
}
