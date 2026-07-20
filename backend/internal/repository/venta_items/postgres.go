package ventaitems

import (
	"context"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx. db may be
// a *pgxpool.Pool or a pgx.Tx, so Insert can run inside CreateVenta's own
// transaction.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) Insert(ctx context.Context, params InsertParams) (generated.VentaItem, error) {
	return r.q.InsertVentaItem(ctx, generated.InsertVentaItemParams{
		VentaID:            params.VentaID,
		TipoLinea:          generated.TipoLineaEnum(params.TipoLinea),
		FraganciaID:        repo.Int8(params.FraganciaID),
		VarianteEnvaseID:   repo.Int8(params.VarianteEnvaseID),
		ProductoID:         repo.Int8(params.ProductoID),
		FeromonaProductoID: repo.Int8(params.FeromonaProductoID),
		GramosFragancia:    repo.NullDecimal(params.GramosFragancia),
		Cantidad:           params.Cantidad,
		PrecioUnitario:     params.PrecioUnitario,
		Subtotal:           params.Subtotal,
	})
}

func (r *postgresRepository) GetByVentaID(ctx context.Context, ventaID int64) ([]generated.GetVentaItemsByVentaIDRow, error) {
	return r.q.GetVentaItemsByVentaID(ctx, ventaID)
}
