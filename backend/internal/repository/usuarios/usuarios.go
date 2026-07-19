// Package usuarios is the persistence port for the usuarios table.
package usuarios

import (
	"context"
	"errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no usuario matches the lookup.
var ErrNotFound = errors.New("usuario not found")

// ListFilter narrows and orders ListPaginated results. Rol == "" and
// Activo == "all" mean "no filter" on that dimension.
type ListFilter struct {
	Page           int
	PageSize       int
	SortCol        string
	SortDir        string
	Q              string
	Rol            string
	Activo         string
	IncludeDeleted bool
}

// UpdateFields carries the optional fields for a partial update: a nil
// pointer means "leave unchanged".
type UpdateFields struct {
	NombreCompleto *string
	Correo         *string
	Rol            *string
}

// Repository is the persistence port for usuarios, consumed by usecases.
type Repository interface {
	GetByCorreo(ctx context.Context, correo string) (generated.Usuario, error)
	GetByID(ctx context.Context, id int64) (generated.Usuario, error)
	UpdateLastLogin(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error

	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.Usuario, int64, error)
	CountActiveAdmins(ctx context.Context) (int64, error)
	ExistsCorreo(ctx context.Context, correo string) (bool, error)
	Insert(ctx context.Context, sedeID int64, nombreCompleto, correo, passwordHash, rol string) (generated.Usuario, error)
	Update(ctx context.Context, id int64, fields UpdateFields) (generated.Usuario, error)
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error
	SoftDelete(ctx context.Context, id int64) error
}
