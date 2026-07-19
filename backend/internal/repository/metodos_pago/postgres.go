package metodospago

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

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.MetodosPago, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListMetodosPagoPaginated(ctx, generated.ListMetodosPagoPaginatedParams{
		Limit:          int32(pageSize),
		Offset:         int32(offset),
		IncludeDeleted: filter.IncludeDeleted,
		Activo:         filter.Activo,
		Q:              filter.Q,
		SortCol:        filter.SortCol,
		SortDir:        filter.SortDir,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountMetodosPago(ctx, generated.CountMetodosPagoParams{
		IncludeDeleted: filter.IncludeDeleted,
		Activo:         filter.Activo,
		Q:              filter.Q,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.MetodosPago, error) {
	m, err := r.q.GetMetodoPagoByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.MetodosPago{}, ErrNotFound
		}
		return generated.MetodosPago{}, err
	}
	return m, nil
}

func (r *postgresRepository) GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.MetodosPago, error) {
	m, err := r.q.GetMetodoPagoByIDIncludingDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.MetodosPago{}, ErrNotFound
		}
		return generated.MetodosPago{}, err
	}
	return m, nil
}

func (r *postgresRepository) ExistsCodigo(ctx context.Context, codigo string, excludeID int64) (bool, error) {
	return r.q.ExistsMetodoPagoCodigo(ctx, generated.ExistsMetodoPagoCodigoParams{
		Codigo:    codigo,
		ExcludeID: excludeID,
	})
}

func (r *postgresRepository) ExistsNombre(ctx context.Context, nombre string, excludeID int64) (bool, error) {
	return r.q.ExistsMetodoPagoNombre(ctx, generated.ExistsMetodoPagoNombreParams{
		Nombre:    nombre,
		ExcludeID: excludeID,
	})
}

func (r *postgresRepository) CountVentas(ctx context.Context, id int64) (int64, error) {
	return r.q.CountVentasByMetodoPago(ctx, id)
}

func (r *postgresRepository) Insert(ctx context.Context, nombre, codigo string) (generated.MetodosPago, error) {
	return r.q.InsertMetodoPago(ctx, generated.InsertMetodoPagoParams{Nombre: nombre, Codigo: codigo})
}

func (r *postgresRepository) Update(ctx context.Context, id int64, fields UpdateFields) (generated.MetodosPago, error) {
	params := generated.UpdateMetodoPagoParams{ID: id}

	if fields.Nombre != nil {
		params.Nombre = repo.Text(*fields.Nombre)
	}
	if fields.Codigo != nil {
		params.Codigo = repo.Text(*fields.Codigo)
	}

	m, err := r.q.UpdateMetodoPago(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.MetodosPago{}, ErrNotFound
		}
		return generated.MetodosPago{}, err
	}
	return m, nil
}

func (r *postgresRepository) SoftDelete(ctx context.Context, id int64) error {
	return r.q.SoftDeleteMetodoPago(ctx, id)
}

func (r *postgresRepository) HardDelete(ctx context.Context, id int64) error {
	return r.q.HardDeleteMetodoPago(ctx, id)
}
