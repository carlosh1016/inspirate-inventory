package variantesenvase

import (
	"context"
	"strings"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

var allowedSortCols = map[string]bool{
	"color":      true,
	"created_at": true,
}

// ListInput is the raw, unvalidated query params for GET /variantes-envase.
type ListInput struct {
	Page           int
	PageSize       int
	Sort           string
	Q              string
	SedeID         int64
	ModeloEnvaseID int64
	Activo         string // "true" | "false" | "all"; default "true"
	StockBajo      bool
	IncludeDeleted bool
}

// ListResult is a page of variantes_envase plus pagination metadata.
type ListResult struct {
	Items    []generated.ListVariantesEnvasePaginatedRow
	Total    int64
	Page     int
	PageSize int
}

// List returns a filtered, sorted, paginated page of variantes_envase.
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

	items, total, err := s.VariantesEnvase.ListPaginated(ctx, repo.ListFilter{
		Page:           page,
		PageSize:       pageSize,
		SortCol:        sortCol,
		SortDir:        sortDir,
		Q:              in.Q,
		SedeID:         in.SedeID,
		ModeloEnvaseID: in.ModeloEnvaseID,
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
		defaultCol = "color"
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
