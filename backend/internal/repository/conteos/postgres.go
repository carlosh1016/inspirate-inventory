package conteos

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

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

func (r *postgresRepository) Insert(ctx context.Context, sedeID int64, periodo time.Time, creadoPorUsuarioID int64) (generated.ConteosInventario, error) {
	return r.q.InsertConteo(ctx, generated.InsertConteoParams{
		SedeID:             sedeID,
		Periodo:            repo.Date(periodo),
		CreadoPorUsuarioID: creadoPorUsuarioID,
	})
}

func (r *postgresRepository) GenerarItems(ctx context.Context, conteoID, sedeID int64, periodo time.Time) error {
	return r.q.GenerarItemsConteo(ctx, generated.GenerarItemsConteoParams{
		ConteoID: conteoID,
		SedeID:   sedeID,
		Periodo:  repo.Date(periodo),
	})
}

func (r *postgresRepository) GetBySedePeriodo(ctx context.Context, sedeID int64, periodo time.Time) (generated.ConteosInventario, error) {
	row, err := r.q.GetConteoBySedePeriodo(ctx, generated.GetConteoBySedePeriodoParams{
		SedeID:  sedeID,
		Periodo: repo.Date(periodo),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.ConteosInventario{}, ErrNotFound
		}
		return generated.ConteosInventario{}, err
	}
	return row, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.GetConteoByIDRow, error) {
	row, err := r.q.GetConteoByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetConteoByIDRow{}, ErrNotFound
		}
		return generated.GetConteoByIDRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) ListItems(ctx context.Context, conteoID int64) ([]generated.ListConteoItemsRow, error) {
	return r.q.ListConteoItems(ctx, conteoID)
}

func (r *postgresRepository) UpdateItemFisico(ctx context.Context, id, conteoID int64, gramosFisico decimal.Decimal) (generated.UpdateConteoItemFisicoRow, error) {
	row, err := r.q.UpdateConteoItemFisico(ctx, generated.UpdateConteoItemFisicoParams{
		ID:           id,
		ConteoID:     conteoID,
		GramosFisico: decimal.NewNullDecimal(gramosFisico),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.UpdateConteoItemFisicoRow{}, ErrNotFound
		}
		return generated.UpdateConteoItemFisicoRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) Cerrar(ctx context.Context, id, cerradoPorUsuarioID int64) (generated.ConteosInventario, error) {
	row, err := r.q.CerrarConteo(ctx, generated.CerrarConteoParams{
		ID:                  id,
		CerradoPorUsuarioID: repo.Int8(&cerradoPorUsuarioID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.ConteosInventario{}, ErrNotFound
		}
		return generated.ConteosInventario{}, err
	}
	return row, nil
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListConteosPaginatedRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListConteosPaginated(ctx, generated.ListConteosPaginatedParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
		SedeID: filter.SedeID,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountConteos(ctx, filter.SedeID)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
