package movimientos

import (
	"context"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

const defaultPageSize = 20

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx. db may be
// a *pgxpool.Pool or a pgx.Tx, so Insert can run inside InventoryService's
// transaction alongside the stock_actual update it belongs with.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListMovimientosPaginatedRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListMovimientosPaginated(ctx, generated.ListMovimientosPaginatedParams{
		Limit:      int32(pageSize),
		Offset:     int32(offset),
		SedeID:     filter.SedeID,
		TipoItem:   filter.TipoItem,
		ItemID:     filter.ItemID,
		Tipo:       filter.Tipo,
		UsuarioID:  filter.UsuarioID,
		Ubicacion:  filter.Ubicacion,
		VentaID:    filter.VentaID,
		FechaDesde: repo.TimestamptzPtr(filter.FechaDesde),
		FechaHasta: repo.TimestamptzPtr(filter.FechaHasta),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountMovimientos(ctx, generated.CountMovimientosParams{
		SedeID:     filter.SedeID,
		TipoItem:   filter.TipoItem,
		ItemID:     filter.ItemID,
		Tipo:       filter.Tipo,
		UsuarioID:  filter.UsuarioID,
		Ubicacion:  filter.Ubicacion,
		VentaID:    filter.VentaID,
		FechaDesde: repo.TimestamptzPtr(filter.FechaDesde),
		FechaHasta: repo.TimestamptzPtr(filter.FechaHasta),
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) Insert(ctx context.Context, params InsertParams) (generated.MovimientosInventario, error) {
	return r.q.InsertMovimiento(ctx, generated.InsertMovimientoParams{
		SedeID:         params.SedeID,
		UsuarioID:      params.UsuarioID,
		TipoItem:       generated.TipoItemEnum(params.TipoItem),
		ItemID:         params.ItemID,
		Tipo:           generated.TipoMovimientoEnum(params.Tipo),
		Ubicacion:      generated.UbicacionEnum(params.Ubicacion),
		Cantidad:       params.Cantidad,
		StockAnterior:  params.StockAnterior,
		StockPosterior: params.StockPosterior,
		Motivo:         repo.TextPtr(params.Motivo),
		VentaID:        repo.Int8(params.VentaID),
	})
}
