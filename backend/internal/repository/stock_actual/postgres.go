package stockactual

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

const defaultPageSize = 20

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx. db may be
// a *pgxpool.Pool or a pgx.Tx, so callers can run within a transaction —
// InitializeStock is always called alongside the item's own insert.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) InitializeStock(ctx context.Context, sedeID int64, tipoItem string, itemID int64) error {
	for _, ubicacion := range [2]generated.UbicacionEnum{generated.UbicacionEnumVitrina, generated.UbicacionEnumBodega} {
		if err := r.q.InsertStockActual(ctx, generated.InsertStockActualParams{
			SedeID:    sedeID,
			TipoItem:  generated.TipoItemEnum(tipoItem),
			ItemID:    itemID,
			Ubicacion: ubicacion,
			Cantidad:  decimal.Zero,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *postgresRepository) GetStockTotal(ctx context.Context, sedeID int64, tipoItem string, itemID int64) (decimal.Decimal, decimal.Decimal, error) {
	row, err := r.q.GetStockTotalByItem(ctx, generated.GetStockTotalByItemParams{
		SedeID:   sedeID,
		TipoItem: generated.TipoItemEnum(tipoItem),
		ItemID:   itemID,
	})
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, err
	}
	return row.Vitrina, row.Bodega, nil
}

func (r *postgresRepository) ListUnificado(ctx context.Context, filter ListUnificadoFilter) ([]generated.ListStockUnificadoRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListStockUnificado(ctx, generated.ListStockUnificadoParams{
		Limit:            int32(pageSize),
		Offset:           int32(offset),
		SedeID:           filter.SedeID,
		TipoItemFilter:   filter.TipoItemFilter,
		StockBajo:        filter.StockBajo,
		StockCero:        filter.StockCero,
		IncludeInactivos: filter.IncludeInactivos,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountStockUnificado(ctx, generated.CountStockUnificadoParams{
		SedeID:           filter.SedeID,
		TipoItemFilter:   filter.TipoItemFilter,
		StockBajo:        filter.StockBajo,
		StockCero:        filter.StockCero,
		IncludeInactivos: filter.IncludeInactivos,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) GetForUpdate(ctx context.Context, sedeID int64, tipoItem string, itemID int64, ubicacion string) (generated.StockActual, bool, error) {
	row, err := r.q.GetStockActualForUpdate(ctx, generated.GetStockActualForUpdateParams{
		SedeID:    sedeID,
		TipoItem:  generated.TipoItemEnum(tipoItem),
		ItemID:    itemID,
		Ubicacion: generated.UbicacionEnum(ubicacion),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.StockActual{}, false, nil
		}
		return generated.StockActual{}, false, err
	}
	return row, true, nil
}

func (r *postgresRepository) SetCantidad(ctx context.Context, sedeID int64, tipoItem string, itemID int64, ubicacion string, cantidad decimal.Decimal) error {
	_, err := r.q.UpsertStockActual(ctx, generated.UpsertStockActualParams{
		SedeID:    sedeID,
		TipoItem:  generated.TipoItemEnum(tipoItem),
		ItemID:    itemID,
		Ubicacion: generated.UbicacionEnum(ubicacion),
		Cantidad:  cantidad,
	})
	return err
}
