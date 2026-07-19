// Package metodospago is the persistence port for the metodos_pago table.
package metodospago

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no metodo_pago matches the lookup.
var ErrNotFound = errors.New("metodo_pago not found")

// ListFilter narrows and orders ListPaginated results.
type ListFilter struct {
	Page           int
	PageSize       int
	SortCol        string
	SortDir        string
	Q              string
	Activo         string
	IncludeDeleted bool
}

// UpdateFields carries the optional fields for a partial update: a nil
// pointer means "leave unchanged".
type UpdateFields struct {
	Nombre *string
	Codigo *string
}

// Repository is the persistence port for metodos_pago, consumed by
// usecases.
type Repository interface {
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.MetodosPago, int64, error)
	GetByID(ctx context.Context, id int64) (generated.MetodosPago, error)
	GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.MetodosPago, error)
	ExistsCodigo(ctx context.Context, codigo string, excludeID int64) (bool, error)
	ExistsNombre(ctx context.Context, nombre string, excludeID int64) (bool, error)
	CountVentas(ctx context.Context, id int64) (int64, error)
	Insert(ctx context.Context, nombre, codigo string) (generated.MetodosPago, error)
	Update(ctx context.Context, id int64, fields UpdateFields) (generated.MetodosPago, error)
	SoftDelete(ctx context.Context, id int64) error
	HardDelete(ctx context.Context, id int64) error
}
