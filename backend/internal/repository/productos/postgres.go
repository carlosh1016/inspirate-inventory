package productos

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

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListProductosPaginatedRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListProductosPaginated(ctx, generated.ListProductosPaginatedParams{
		Limit:          int32(pageSize),
		Offset:         int32(offset),
		IncludeDeleted: filter.IncludeDeleted,
		SedeID:         filter.SedeID,
		Categoria:      filter.Categoria,
		Activo:         filter.Activo,
		Q:              filter.Q,
		StockBajo:      filter.StockBajo,
		SortCol:        filter.SortCol,
		SortDir:        filter.SortDir,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountProductos(ctx, generated.CountProductosParams{
		IncludeDeleted: filter.IncludeDeleted,
		SedeID:         filter.SedeID,
		Categoria:      filter.Categoria,
		Activo:         filter.Activo,
		Q:              filter.Q,
		StockBajo:      filter.StockBajo,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.GetProductoByIDRow, error) {
	row, err := r.q.GetProductoByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetProductoByIDRow{}, ErrNotFound
		}
		return generated.GetProductoByIDRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) GetByIDIncludingDeleted(ctx context.Context, id int64) (generated.Producto, error) {
	p, err := r.q.GetProductoByIDIncludingDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Producto{}, ErrNotFound
		}
		return generated.Producto{}, err
	}
	return p, nil
}

func (r *postgresRepository) ExistsNombreCategoria(ctx context.Context, sedeID int64, nombre, categoria string, excludeID int64) (bool, error) {
	return r.q.ExistsProductoNombreCategoria(ctx, generated.ExistsProductoNombreCategoriaParams{
		SedeID:    sedeID,
		Nombre:    nombre,
		Categoria: generated.CategoriaProductoEnum(categoria),
		ExcludeID: excludeID,
	})
}

func (r *postgresRepository) Insert(ctx context.Context, sedeID int64, nombre, categoria, precio string, stockMinimo int32) (generated.Producto, error) {
	precioDec, err := decimal.NewFromString(precio)
	if err != nil {
		return generated.Producto{}, err
	}

	return r.q.InsertProducto(ctx, generated.InsertProductoParams{
		SedeID:      sedeID,
		Nombre:      nombre,
		Categoria:   generated.CategoriaProductoEnum(categoria),
		Precio:      precioDec,
		StockMinimo: stockMinimo,
	})
}

func (r *postgresRepository) Update(ctx context.Context, id int64, fields UpdateFields) (generated.Producto, error) {
	params := generated.UpdateProductoParams{ID: id}

	if fields.Nombre != nil {
		params.Nombre = repo.Text(*fields.Nombre)
	}
	if fields.Categoria != nil {
		params.Categoria = generated.NullCategoriaProductoEnum{CategoriaProductoEnum: generated.CategoriaProductoEnum(*fields.Categoria), Valid: true}
	}
	if fields.Precio != nil {
		precioDec, err := decimal.NewFromString(*fields.Precio)
		if err != nil {
			return generated.Producto{}, err
		}
		params.Precio = decimal.NullDecimal{Decimal: precioDec, Valid: true}
	}
	if fields.StockMinimo != nil {
		params.StockMinimo = repo.Int4(fields.StockMinimo)
	}

	p, err := r.q.UpdateProducto(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Producto{}, ErrNotFound
		}
		return generated.Producto{}, err
	}
	return p, nil
}

func (r *postgresRepository) SoftDelete(ctx context.Context, id int64) error {
	return r.q.SoftDeleteProducto(ctx, id)
}
