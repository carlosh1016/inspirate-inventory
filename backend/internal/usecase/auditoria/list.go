package auditoria

import (
	"context"
	"time"

	domainauditoria "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auditoria"
	auditoriarepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
)

// ListInput are the filters and pagination for List. Zero-valued filters mean
// "no filter".
type ListInput struct {
	Page          int
	PageSize      int
	UsuarioID     int64
	Accion        string
	TablaAfectada string
	FechaDesde    *time.Time
	FechaHasta    *time.Time
}

// ListResult is a page of eventos plus pagination metadata.
type ListResult struct {
	Eventos  []domainauditoria.Evento
	Total    int64
	Page     int
	PageSize int
}

// List returns audit events matching the filters, newest first.
func (s *Service) List(ctx context.Context, in ListInput) (ListResult, error) {
	limit, offset, page, pageSize := normalizePaging(in.Page, in.PageSize)

	filtro := auditoriarepo.ListFiltro{
		UsuarioID:     in.UsuarioID,
		Accion:        in.Accion,
		TablaAfectada: in.TablaAfectada,
		FechaDesde:    in.FechaDesde,
		FechaHasta:    in.FechaHasta,
		Limit:         limit,
		Offset:        offset,
	}

	eventos, err := s.repo.List(ctx, filtro)
	if err != nil {
		return ListResult{}, err
	}
	total, err := s.repo.Count(ctx, filtro)
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{Eventos: eventos, Total: total, Page: page, PageSize: pageSize}, nil
}
