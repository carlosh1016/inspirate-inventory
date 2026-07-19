package stockactual

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

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
