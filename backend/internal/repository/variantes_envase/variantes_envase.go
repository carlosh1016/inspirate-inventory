// Package variantesenvase is the persistence port for the variantes_envase
// table.
package variantesenvase

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no variante_envase matches the lookup.
var ErrNotFound = errors.New("variante_envase not found")

// ListFilter narrows and orders ListPaginated results.
type ListFilter struct {
	Page           int
	PageSize       int
	SortCol        string
	SortDir        string
	Q              string
	SedeID         int64
	ModeloEnvaseID int64
	Activo         string
	StockBajo      bool
	IncludeDeleted bool
}

// UpdateFields carries the optional fields for a partial update: a nil
// pointer means "leave unchanged". Both admin and vendedora share this same
// set — variantes_envase has no field that's admin-only on update.
type UpdateFields struct {
	Color       *string
	StockMinimo *int32
}

// Repository is the persistence port for variantes_envase, consumed by
// usecases.
type Repository interface {
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListVariantesEnvasePaginatedRow, int64, error)
	GetByID(ctx context.Context, id int64) (generated.GetVarianteEnvaseByIDRow, error)
	GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.VariantesEnvase, error)
	ExistsColor(ctx context.Context, modeloEnvaseID int64, color string, excludeID int64) (bool, error)
	Insert(ctx context.Context, modeloEnvaseID, sedeID int64, color string, stockMinimo int32) (generated.VariantesEnvase, error)
	Update(ctx context.Context, id int64, fields UpdateFields) (generated.VariantesEnvase, error)
	SoftDelete(ctx context.Context, id int64) error
}
