package cuadres

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

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx. db may be
// a *pgxpool.Pool or a pgx.Tx.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListCuadresPaginatedRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListCuadresPaginated(ctx, generated.ListCuadresPaginatedParams{
		Limit:      int32(pageSize),
		Offset:     int32(offset),
		SedeID:     filter.SedeID,
		Estado:     filter.Estado,
		FechaDesde: repo.DatePtr(filter.FechaDesde),
		FechaHasta: repo.DatePtr(filter.FechaHasta),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountCuadres(ctx, generated.CountCuadresParams{
		SedeID:     filter.SedeID,
		Estado:     filter.Estado,
		FechaDesde: repo.DatePtr(filter.FechaDesde),
		FechaHasta: repo.DatePtr(filter.FechaHasta),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.GetCuadreByIDRow, error) {
	row, err := r.q.GetCuadreByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetCuadreByIDRow{}, ErrNotFound
		}
		return generated.GetCuadreByIDRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) GetBySedeFecha(ctx context.Context, sedeID int64, fecha time.Time) (generated.GetCuadreBySedeFechaRow, error) {
	row, err := r.q.GetCuadreBySedeFecha(ctx, generated.GetCuadreBySedeFechaParams{
		SedeID: sedeID,
		Fecha:  repo.Date(fecha),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetCuadreBySedeFechaRow{}, ErrNotFound
		}
		return generated.GetCuadreBySedeFechaRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) GetAbiertoAnterior(ctx context.Context, sedeID int64, fecha time.Time) (generated.CuadresCaja, error) {
	row, err := r.q.GetCuadreAbiertoAnterior(ctx, generated.GetCuadreAbiertoAnteriorParams{
		SedeID: sedeID,
		Fecha:  repo.Date(fecha),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.CuadresCaja{}, ErrNotFound
		}
		return generated.CuadresCaja{}, err
	}
	return row, nil
}

func (r *postgresRepository) Insert(ctx context.Context, params InsertParams) (generated.CuadresCaja, error) {
	return r.q.InsertCuadre(ctx, generated.InsertCuadreParams{
		SedeID:    params.SedeID,
		Fecha:     repo.Date(params.Fecha),
		FondoBase: params.FondoBase,
	})
}

func (r *postgresRepository) Cerrar(ctx context.Context, params CerrarParams) (generated.CuadresCaja, error) {
	row, err := r.q.UpdateCuadreCerrar(ctx, generated.UpdateCuadreCerrarParams{
		ID:                  params.ID,
		TotalEfectivo:       params.TotalEfectivo,
		TotalNequi:          params.TotalNequi,
		TotalDaviplata:      params.TotalDaviplata,
		TotalTransferencia:  params.TotalTransferencia,
		TotalOtros:          params.TotalOtros,
		TotalPagos:          params.TotalPagos,
		TotalConsignaciones: params.TotalConsignaciones,
		ValorTurno:          params.ValorTurno,
		SaldoCalculado:      params.SaldoCalculado,
		Observaciones:       repo.TextPtr(params.Observaciones),
		CerradoPorUsuarioID: repo.Int8(&params.CerradoPorUsuarioID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.CuadresCaja{}, ErrNotFound
		}
		return generated.CuadresCaja{}, err
	}
	return row, nil
}

func (r *postgresRepository) ExistsCerradoBySedeFecha(ctx context.Context, sedeID int64, fecha time.Time) (bool, error) {
	return r.q.ExistsCuadreCerradoBySedeFecha(ctx, generated.ExistsCuadreCerradoBySedeFechaParams{
		SedeID: sedeID,
		Fecha:  repo.Date(fecha),
	})
}

func (r *postgresRepository) TotalesPorMetodoEnFecha(ctx context.Context, sedeID int64, diaInicio, diaFin time.Time) (generated.GetTotalesPorMetodoEnFechaRow, error) {
	return r.q.GetTotalesPorMetodoEnFecha(ctx, generated.GetTotalesPorMetodoEnFechaParams{
		SedeID:    sedeID,
		DiaInicio: repo.Timestamptz(diaInicio),
		DiaFin:    repo.Timestamptz(diaFin),
	})
}
