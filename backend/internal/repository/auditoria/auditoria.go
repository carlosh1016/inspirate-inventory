// Package auditoria is the persistence port for the auditoria table: Insert
// (write side, used across many usecases since Tanda 1) plus the read side
// (List/Count/GetByID/AccionesDistintas) added in Tanda 6.
package auditoria

import (
	"context"
	"errors"
	"time"

	domainauditoria "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auditoria"
)

// ErrNotFound is returned by GetByID when no evento matches the id.
var ErrNotFound = errors.New("evento de auditoría no encontrado")

// Entry is one row to record in `auditoria`. UsuarioID is nil for events
// like a failed login with an unknown correo.
type Entry struct {
	UsuarioID     *int64
	Accion        string
	TablaAfectada *string
	RegistroID    *int64
	DatosAntes    []byte
	DatosDespues  []byte
	IP            string
	UserAgent     string
}

// ListFiltro carries the optional filters and pagination for List/Count.
// Zero-valued numeric/string filters mean "no filter"; nil dates mean unbounded.
type ListFiltro struct {
	UsuarioID     int64
	Accion        string
	TablaAfectada string
	FechaDesde    *time.Time
	FechaHasta    *time.Time
	Limit         int32
	Offset        int32
}

// Repository is the persistence port for auditoria, consumed by usecases.
type Repository interface {
	Insert(ctx context.Context, e Entry) error
	List(ctx context.Context, f ListFiltro) ([]domainauditoria.Evento, error)
	Count(ctx context.Context, f ListFiltro) (int64, error)
	GetByID(ctx context.Context, id int64) (domainauditoria.Evento, error)
	AccionesDistintas(ctx context.Context) ([]string, error)
}
