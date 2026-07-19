// Package fragancias is the persistence port for the fragancias table.
package fragancias

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no fragancia matches the lookup.
var ErrNotFound = errors.New("fragancia not found")

// ListFilter narrows and orders ListPaginated results.
type ListFilter struct {
	Page           int
	PageSize       int
	SortCol        string
	SortDir        string
	Q              string
	SedeID         int64
	Genero         string
	Activo         string
	StockBajo      bool
	IncludeDeleted bool
}

// UpdateFields carries the optional fields for a partial update: a nil
// pointer means "leave unchanged".
type UpdateFields struct {
	NombreComercial   *string
	NombreAlternativo *string
	Genero            *string
	GramosMinimo      *string // decimal string; parsed by the caller
}

// Repository is the persistence port for fragancias, consumed by usecases.
type Repository interface {
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListFraganciasPaginatedRow, int64, error)
	GetByID(ctx context.Context, id int64) (generated.GetFraganciaByIDRow, error)
	GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.Fragancia, error)
	ExistsNombreComercial(ctx context.Context, sedeID int64, nombre string, excludeID int64) (bool, error)
	Insert(ctx context.Context, sedeID int64, nombreComercial string, nombreAlternativo *string, genero string, gramosMinimo string) (generated.Fragancia, error)
	Update(ctx context.Context, id int64, fields UpdateFields) (generated.Fragancia, error)
	SoftDelete(ctx context.Context, id int64) error
	Restore(ctx context.Context, id int64) (generated.Fragancia, error)
}
