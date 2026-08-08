package fragancias

import (
	"context"
	"strings"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

var allowedSortCols = map[string]bool{
	"nombre_comercial": true,
	"created_at":       true,
}

// ListInput is the raw, unvalidated query params for GET /fragancias.
// Activo filters on the activo column (a manual on/off flag); IncludeDeleted
// is the separate, independent toggle for soft-deleted rows — an admin
// needs it to find something to send to POST /:id/restaurar.
type ListInput struct {
	Page           int
	PageSize       int
	Sort           string
	Q              string
	SedeID         int64
	Genero         string
	NumeroGenero   int32
	Activo         string // "true" | "false" | "all"; default "true"
	StockBajo      bool
	IncludeDeleted bool
}

// ListResult is a page of fragancias plus pagination metadata.
type ListResult struct {
	Items    []generated.ListFraganciasPaginatedRow
	Total    int64
	Page     int
	PageSize int
}

// List returns a filtered, sorted, paginated page of fragancias.
func (s *Service) List(ctx context.Context, in ListInput) (ListResult, error) {
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

	activo := in.Activo
	if activo == "" {
		activo = "true"
	}

	sortCol, sortDir := parseSort(in.Sort)

	items, total, err := s.Fragancias.ListPaginated(ctx, repo.ListFilter{
		Page:           page,
		PageSize:       pageSize,
		SortCol:        sortCol,
		SortDir:        sortDir,
		Q:              in.Q,
		SedeID:         in.SedeID,
		Genero:         in.Genero,
		NumeroGenero:   in.NumeroGenero,
		Activo:         activo,
		StockBajo:      in.StockBajo,
		IncludeDeleted: in.IncludeDeleted,
	})
	if err != nil {
		return ListResult{}, internalErr(err)
	}

	return ListResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

// parseSort validates "campo:direccion" against a column whitelist, falling
// back to a sane default for anything malformed or disallowed.
func parseSort(sort string) (col, dir string) {
	const (
		defaultCol = "nombre_comercial"
		defaultDir = "asc"
	)

	parts := strings.SplitN(sort, ":", 2)
	if len(parts) != 2 {
		return defaultCol, defaultDir
	}

	col, dir = parts[0], parts[1]
	if !allowedSortCols[col] || (dir != "asc" && dir != "desc") {
		return defaultCol, defaultDir
	}
	return col, dir
}
