package modelosenvase

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

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListModelosEnvasePaginatedRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListModelosEnvasePaginated(ctx, generated.ListModelosEnvasePaginatedParams{
		Limit:          int32(pageSize),
		Offset:         int32(offset),
		IncludeDeleted: filter.IncludeDeleted,
		Tipo:           filter.Tipo,
		Activo:         filter.Activo,
		Q:              filter.Q,
		SortCol:        filter.SortCol,
		SortDir:        filter.SortDir,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountModelosEnvase(ctx, generated.CountModelosEnvaseParams{
		IncludeDeleted: filter.IncludeDeleted,
		Tipo:           filter.Tipo,
		Activo:         filter.Activo,
		Q:              filter.Q,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.GetModeloEnvaseByIDRow, error) {
	row, err := r.q.GetModeloEnvaseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetModeloEnvaseByIDRow{}, ErrNotFound
		}
		return generated.GetModeloEnvaseByIDRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.ModelosEnvase, error) {
	m, err := r.q.GetModeloEnvaseByIDIncludingDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.ModelosEnvase{}, ErrNotFound
		}
		return generated.ModelosEnvase{}, err
	}
	return m, nil
}

func (r *postgresRepository) ExistsTipoTamano(ctx context.Context, tipo, tamanoOz string, excludeID int64) (bool, error) {
	tamano, err := decimal.NewFromString(tamanoOz)
	if err != nil {
		return false, err
	}
	return r.q.ExistsModeloEnvaseTipoTamano(ctx, generated.ExistsModeloEnvaseTipoTamanoParams{
		Tipo:      tipo,
		TamanoOz:  tamano,
		ExcludeID: excludeID,
	})
}

func (r *postgresRepository) CountVariantesActivas(ctx context.Context, modeloEnvaseID int64) (int64, error) {
	return r.q.CountVariantesActivasByModelo(ctx, modeloEnvaseID)
}

func (r *postgresRepository) Insert(ctx context.Context, tipo, tamanoOz, equivGramos, precioSolo, precioConFragancia, precioRecarga string, tieneVariantes bool) (generated.ModelosEnvase, error) {
	tamano, err := decimal.NewFromString(tamanoOz)
	if err != nil {
		return generated.ModelosEnvase{}, err
	}
	equiv, err := decimal.NewFromString(equivGramos)
	if err != nil {
		return generated.ModelosEnvase{}, err
	}
	solo, err := decimal.NewFromString(precioSolo)
	if err != nil {
		return generated.ModelosEnvase{}, err
	}
	conFragancia, err := decimal.NewFromString(precioConFragancia)
	if err != nil {
		return generated.ModelosEnvase{}, err
	}
	recarga, err := decimal.NewFromString(precioRecarga)
	if err != nil {
		return generated.ModelosEnvase{}, err
	}

	return r.q.InsertModeloEnvase(ctx, generated.InsertModeloEnvaseParams{
		Tipo:               tipo,
		TamanoOz:           tamano,
		EquivGramos:        equiv,
		PrecioSolo:         solo,
		PrecioConFragancia: conFragancia,
		PrecioRecarga:      recarga,
		TieneVariantes:     tieneVariantes,
	})
}

func (r *postgresRepository) Update(ctx context.Context, id int64, fields UpdateFields) (generated.ModelosEnvase, error) {
	params := generated.UpdateModeloEnvaseParams{ID: id}

	if fields.Tipo != nil {
		params.Tipo = repo.Text(*fields.Tipo)
	}
	if fields.TamanoOz != nil {
		d, err := decimal.NewFromString(*fields.TamanoOz)
		if err != nil {
			return generated.ModelosEnvase{}, err
		}
		params.TamanoOz = decimal.NullDecimal{Decimal: d, Valid: true}
	}
	if fields.EquivGramos != nil {
		d, err := decimal.NewFromString(*fields.EquivGramos)
		if err != nil {
			return generated.ModelosEnvase{}, err
		}
		params.EquivGramos = decimal.NullDecimal{Decimal: d, Valid: true}
	}
	if fields.PrecioSolo != nil {
		d, err := decimal.NewFromString(*fields.PrecioSolo)
		if err != nil {
			return generated.ModelosEnvase{}, err
		}
		params.PrecioSolo = decimal.NullDecimal{Decimal: d, Valid: true}
	}
	if fields.PrecioConFragancia != nil {
		d, err := decimal.NewFromString(*fields.PrecioConFragancia)
		if err != nil {
			return generated.ModelosEnvase{}, err
		}
		params.PrecioConFragancia = decimal.NullDecimal{Decimal: d, Valid: true}
	}
	if fields.PrecioRecarga != nil {
		d, err := decimal.NewFromString(*fields.PrecioRecarga)
		if err != nil {
			return generated.ModelosEnvase{}, err
		}
		params.PrecioRecarga = decimal.NullDecimal{Decimal: d, Valid: true}
	}

	m, err := r.q.UpdateModeloEnvase(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.ModelosEnvase{}, ErrNotFound
		}
		return generated.ModelosEnvase{}, err
	}
	return m, nil
}

func (r *postgresRepository) SoftDelete(ctx context.Context, id int64) error {
	return r.q.SoftDeleteModeloEnvase(ctx, id)
}
