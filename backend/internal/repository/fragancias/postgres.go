package fragancias

import (
	"context"
	"errors"

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
// a *pgxpool.Pool or a pgx.Tx, so Insert can run inside a transaction
// alongside the initial stock rows.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListFraganciasPaginatedRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListFraganciasPaginated(ctx, generated.ListFraganciasPaginatedParams{
		Limit:          int32(pageSize),
		Offset:         int32(offset),
		IncludeDeleted: filter.IncludeDeleted,
		SedeID:         filter.SedeID,
		Genero:         filter.Genero,
		NumeroGenero:   filter.NumeroGenero,
		Activo:         filter.Activo,
		Q:              filter.Q,
		StockBajo:      filter.StockBajo,
		SortCol:        filter.SortCol,
		SortDir:        filter.SortDir,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountFragancias(ctx, generated.CountFraganciasParams{
		IncludeDeleted: filter.IncludeDeleted,
		SedeID:         filter.SedeID,
		Genero:         filter.Genero,
		NumeroGenero:   filter.NumeroGenero,
		Activo:         filter.Activo,
		Q:              filter.Q,
		StockBajo:      filter.StockBajo,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.GetFraganciaByIDRow, error) {
	row, err := r.q.GetFraganciaByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetFraganciaByIDRow{}, ErrNotFound
		}
		return generated.GetFraganciaByIDRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.Fragancia, error) {
	f, err := r.q.GetFraganciaByIDIncludingDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Fragancia{}, ErrNotFound
		}
		return generated.Fragancia{}, err
	}
	return f, nil
}

func (r *postgresRepository) ExistsNombreComercial(ctx context.Context, sedeID int64, nombre string, excludeID int64) (bool, error) {
	return r.q.ExistsFraganciaNombreComercial(ctx, generated.ExistsFraganciaNombreComercialParams{
		SedeID:          sedeID,
		NombreComercial: nombre,
		ExcludeID:       excludeID,
	})
}

func (r *postgresRepository) ExistsNumeroGenero(ctx context.Context, sedeID int64, genero string, numeroGenero int32, excludeID int64) (bool, error) {
	return r.q.ExistsFraganciaNumeroGenero(ctx, generated.ExistsFraganciaNumeroGeneroParams{
		SedeID:       sedeID,
		Genero:       generated.GeneroEnum(genero),
		NumeroGenero: numeroGenero,
		ExcludeID:    excludeID,
	})
}

func (r *postgresRepository) NextNumeroGenero(ctx context.Context, sedeID int64, genero string) (int32, error) {
	return r.q.NextNumeroGeneroFragancia(ctx, generated.NextNumeroGeneroFraganciaParams{
		SedeID: sedeID,
		Genero: generated.GeneroEnum(genero),
	})
}

func (r *postgresRepository) Insert(ctx context.Context, sedeID int64, nombreComercial string, nombreAlternativo *string, genero string, gramosMinimo string, numeroGenero int32) (generated.Fragancia, error) {
	gramos, err := decimal.NewFromString(gramosMinimo)
	if err != nil {
		return generated.Fragancia{}, err
	}

	return r.q.InsertFragancia(ctx, generated.InsertFraganciaParams{
		SedeID:            sedeID,
		NombreComercial:   nombreComercial,
		NombreAlternativo: repo.TextPtr(nombreAlternativo),
		Genero:            generated.GeneroEnum(genero),
		GramosMinimo:      gramos,
		NumeroGenero:      numeroGenero,
	})
}

func (r *postgresRepository) Update(ctx context.Context, id int64, fields UpdateFields) (generated.Fragancia, error) {
	params := generated.UpdateFraganciaParams{ID: id}

	if fields.NombreComercial != nil {
		params.NombreComercial = repo.Text(*fields.NombreComercial)
	}
	if fields.NombreAlternativo != nil {
		params.NombreAlternativo = repo.Text(*fields.NombreAlternativo)
	}
	if fields.Genero != nil {
		params.Genero = generated.NullGeneroEnum{GeneroEnum: generated.GeneroEnum(*fields.Genero), Valid: true}
	}
	if fields.GramosMinimo != nil {
		gramos, err := decimal.NewFromString(*fields.GramosMinimo)
		if err != nil {
			return generated.Fragancia{}, err
		}
		params.GramosMinimo = decimal.NullDecimal{Decimal: gramos, Valid: true}
	}
	if fields.NumeroGenero != nil {
		params.NumeroGenero = repo.Int4(fields.NumeroGenero)
	}

	f, err := r.q.UpdateFragancia(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Fragancia{}, ErrNotFound
		}
		return generated.Fragancia{}, err
	}
	return f, nil
}

func (r *postgresRepository) SoftDelete(ctx context.Context, id int64) error {
	return r.q.SoftDeleteFragancia(ctx, id)
}

func (r *postgresRepository) Restore(ctx context.Context, id int64) (generated.Fragancia, error) {
	f, err := r.q.RestoreFragancia(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Fragancia{}, ErrNotFound
		}
		return generated.Fragancia{}, err
	}
	return f, nil
}
