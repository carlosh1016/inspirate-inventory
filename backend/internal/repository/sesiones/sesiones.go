// Package sesiones is the persistence port for sesiones_laborales, consumed
// by usecase/sesiones.
package sesiones

import (
	"context"
	"errors"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no sesion matches the lookup, and also when
// CerrarSesion's WHERE ... AND salida_at IS NULL guard matches zero rows.
var ErrNotFound = errors.New("sesion not found")

// ListFilter narrows and paginates ListPaginated results. Zero values (0,
// nil, false) mean "no filter" for each field.
type ListFilter struct {
	Page       int
	PageSize   int
	SedeID     int64
	UsuarioID  int64
	FechaDesde *time.Time
	FechaHasta *time.Time
	Abiertas   bool
}

// UpdateManualParams is what UpdateManual needs — nil means "keep the
// existing value" for that column (mirrors sqlc.narg's COALESCE).
type UpdateManualParams struct {
	ID        int64
	EntradaAt *time.Time
	SalidaAt  *time.Time
}

// Repository is the persistence port for sesiones_laborales.
type Repository interface {
	Insert(ctx context.Context, sedeID, usuarioID int64, entradaAt time.Time) (generated.SesionesLaborale, error)
	GetAbiertaPorUsuario(ctx context.Context, usuarioID int64) (generated.SesionesLaborale, error)
	Cerrar(ctx context.Context, id int64, salidaAt time.Time) (generated.SesionesLaborale, error)
	GetByID(ctx context.Context, id int64) (generated.GetSesionByIDRow, error)
	ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListSesionesRow, int64, error)
	UpdateManual(ctx context.Context, params UpdateManualParams) (generated.SesionesLaborale, error)
	GetResumen(ctx context.Context, fechaDesde, fechaHasta time.Time, usuarioID int64) ([]generated.GetResumenSesionesRow, error)
}
