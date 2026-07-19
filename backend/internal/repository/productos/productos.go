// Package productos is the persistence port for the productos table.
package productos

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no producto matches the lookup.
var ErrNotFound = errors.New("producto not found")

// ListFilter narrows and orders ListPaginated results.
type ListFilter struct {
	Page           int
	PageSize       int
	SortCol        string
	SortDir        string
	Q              string
	SedeID         int64
	Categoria      string
	Activo         string
	StockBajo      bool
	IncludeDeleted bool
}

// UpdateFields carries the optional fields for a partial update: a nil
// pointer means "leave unchanged".
type UpdateFields struct {
	Nombre      *string
	Categoria   *string
	Precio      *string // decimal string; parsed by the caller
	StockMinimo *int32
}

// Repository is the persistence port for productos, consumed by usecases.
type Repository interface {
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListProductosPaginatedRow, int64, error)
	GetByID(ctx context.Context, id int64) (generated.GetProductoByIDRow, error)
	GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.Producto, error)
	ExistsNombreCategoria(ctx context.Context, sedeID int64, nombre, categoria string, excludeID int64) (bool, error)
	Insert(ctx context.Context, sedeID int64, nombre, categoria, precio string, stockMinimo int32) (generated.Producto, error)
	Update(ctx context.Context, id int64, fields UpdateFields) (generated.Producto, error)
	SoftDelete(ctx context.Context, id int64) error
}
