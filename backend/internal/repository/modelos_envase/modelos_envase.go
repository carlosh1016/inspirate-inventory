// Package modelosenvase is the persistence port for the modelos_envase
// table.
package modelosenvase

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no modelo_envase matches the lookup.
var ErrNotFound = errors.New("modelo_envase not found")

// ListFilter narrows and orders ListPaginated results.
type ListFilter struct {
	Page           int
	PageSize       int
	SortCol        string
	SortDir        string
	Q              string
	Tipo           string
	Activo         string
	IncludeDeleted bool
}

// UpdateFields carries the optional fields for a partial update: a nil
// pointer means "leave unchanged".
type UpdateFields struct {
	Tipo               *string
	TamanoOz           *string // decimal string; parsed by the caller
	EquivGramos        *string
	PrecioSolo         *string
	PrecioConFragancia *string
	PrecioRecarga      *string
}

// Repository is the persistence port for modelos_envase, consumed by
// usecases.
type Repository interface {
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListModelosEnvasePaginatedRow, int64, error)
	GetByID(ctx context.Context, id int64) (generated.GetModeloEnvaseByIDRow, error)
	GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.ModelosEnvase, error)
	ExistsTipoTamano(ctx context.Context, tipo, tamanoOz string, excludeID int64) (bool, error)
	CountVariantesActivas(ctx context.Context, modeloEnvaseID int64) (int64, error)
	Insert(ctx context.Context, tipo, tamanoOz, equivGramos, precioSolo, precioConFragancia, precioRecarga string) (generated.ModelosEnvase, error)
	Update(ctx context.Context, id int64, fields UpdateFields) (generated.ModelosEnvase, error)
	SoftDelete(ctx context.Context, id int64) error
}
