package sesiones

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

const defaultPageSize = 20

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) Insert(ctx context.Context, sedeID, usuarioID int64, entradaAt time.Time) (generated.SesionesLaborale, error) {
	return r.q.InsertSesion(ctx, generated.InsertSesionParams{
		SedeID:    sedeID,
		UsuarioID: usuarioID,
		EntradaAt: repo.Timestamptz(entradaAt),
	})
}

func (r *postgresRepository) GetAbiertaPorUsuario(ctx context.Context, usuarioID int64) (generated.SesionesLaborale, error) {
	row, err := r.q.GetSesionAbiertaPorUsuario(ctx, usuarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.SesionesLaborale{}, ErrNotFound
		}
		return generated.SesionesLaborale{}, err
	}
	return row, nil
}

func (r *postgresRepository) Cerrar(ctx context.Context, id int64, salidaAt time.Time) (generated.SesionesLaborale, error) {
	row, err := r.q.CerrarSesion(ctx, generated.CerrarSesionParams{
		ID:       id,
		SalidaAt: repo.Timestamptz(salidaAt),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.SesionesLaborale{}, ErrNotFound
		}
		return generated.SesionesLaborale{}, err
	}
	return row, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.GetSesionByIDRow, error) {
	row, err := r.q.GetSesionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetSesionByIDRow{}, ErrNotFound
		}
		return generated.GetSesionByIDRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListSesionesRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListSesiones(ctx, generated.ListSesionesParams{
		Limit:      int32(pageSize),
		Offset:     int32(offset),
		SedeID:     filter.SedeID,
		UsuarioID:  filter.UsuarioID,
		FechaDesde: repo.TimestamptzPtr(filter.FechaDesde),
		FechaHasta: repo.TimestamptzPtr(filter.FechaHasta),
		Abiertas:   filter.Abiertas,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountSesiones(ctx, generated.CountSesionesParams{
		SedeID:     filter.SedeID,
		UsuarioID:  filter.UsuarioID,
		FechaDesde: repo.TimestamptzPtr(filter.FechaDesde),
		FechaHasta: repo.TimestamptzPtr(filter.FechaHasta),
		Abiertas:   filter.Abiertas,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) UpdateManual(ctx context.Context, params UpdateManualParams) (generated.SesionesLaborale, error) {
	return r.q.UpdateSesionManual(ctx, generated.UpdateSesionManualParams{
		ID:        params.ID,
		EntradaAt: repo.TimestamptzPtr(params.EntradaAt),
		SalidaAt:  repo.TimestamptzPtr(params.SalidaAt),
	})
}

func (r *postgresRepository) GetResumen(ctx context.Context, fechaDesde, fechaHasta time.Time, usuarioID int64) ([]generated.GetResumenSesionesRow, error) {
	return r.q.GetResumenSesiones(ctx, generated.GetResumenSesionesParams{
		FechaDesde: repo.Timestamptz(fechaDesde),
		FechaHasta: repo.Timestamptz(fechaHasta),
		UsuarioID:  usuarioID,
	})
}
