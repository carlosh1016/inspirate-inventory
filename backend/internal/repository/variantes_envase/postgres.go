package variantesenvase

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

const defaultPageSize = 20

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx. db may be
// a *pgxpool.Pool or a pgx.Tx, so Insert can run inside a transaction
// alongside the initial stock rows.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListVariantesEnvasePaginatedRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListVariantesEnvasePaginated(ctx, generated.ListVariantesEnvasePaginatedParams{
		Limit:          int32(pageSize),
		Offset:         int32(offset),
		IncludeDeleted: filter.IncludeDeleted,
		SedeID:         filter.SedeID,
		ModeloEnvaseID: filter.ModeloEnvaseID,
		Activo:         filter.Activo,
		Q:              filter.Q,
		StockBajo:      filter.StockBajo,
		SortCol:        filter.SortCol,
		SortDir:        filter.SortDir,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountVariantesEnvase(ctx, generated.CountVariantesEnvaseParams{
		IncludeDeleted: filter.IncludeDeleted,
		SedeID:         filter.SedeID,
		ModeloEnvaseID: filter.ModeloEnvaseID,
		Activo:         filter.Activo,
		Q:              filter.Q,
		StockBajo:      filter.StockBajo,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.GetVarianteEnvaseByIDRow, error) {
	row, err := r.q.GetVarianteEnvaseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetVarianteEnvaseByIDRow{}, ErrNotFound
		}
		return generated.GetVarianteEnvaseByIDRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.VariantesEnvase, error) {
	v, err := r.q.GetVarianteEnvaseByIDIncludingDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.VariantesEnvase{}, ErrNotFound
		}
		return generated.VariantesEnvase{}, err
	}
	return v, nil
}

func (r *postgresRepository) ExistsColor(ctx context.Context, modeloEnvaseID int64, color string, excludeID int64) (bool, error) {
	return r.q.ExistsVarianteEnvaseColor(ctx, generated.ExistsVarianteEnvaseColorParams{
		ModeloEnvaseID: modeloEnvaseID,
		Color:          color,
		ExcludeID:      excludeID,
	})
}

func (r *postgresRepository) Insert(ctx context.Context, modeloEnvaseID, sedeID int64, color string, stockMinimo int32) (generated.VariantesEnvase, error) {
	return r.q.InsertVarianteEnvase(ctx, generated.InsertVarianteEnvaseParams{
		ModeloEnvaseID: modeloEnvaseID,
		SedeID:         sedeID,
		Color:          color,
		StockMinimo:    stockMinimo,
	})
}

func (r *postgresRepository) Update(ctx context.Context, id int64, fields UpdateFields) (generated.VariantesEnvase, error) {
	params := generated.UpdateVarianteEnvaseParams{ID: id}

	if fields.Color != nil {
		params.Color = repo.Text(*fields.Color)
	}
	if fields.StockMinimo != nil {
		params.StockMinimo = repo.Int4(fields.StockMinimo)
	}

	v, err := r.q.UpdateVarianteEnvase(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.VariantesEnvase{}, ErrNotFound
		}
		return generated.VariantesEnvase{}, err
	}
	return v, nil
}

func (r *postgresRepository) SoftDelete(ctx context.Context, id int64) error {
	return r.q.SoftDeleteVarianteEnvase(ctx, id)
}
