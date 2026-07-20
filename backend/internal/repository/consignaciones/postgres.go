package consignaciones

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) Insert(ctx context.Context, params InsertParams) (generated.Consignacione, error) {
	return r.q.InsertConsignacion(ctx, generated.InsertConsignacionParams{
		CuadreCajaID: params.CuadreCajaID,
		UsuarioID:    params.UsuarioID,
		Monto:        params.Monto,
		Banco:        repo.TextPtr(params.Banco),
		Referencia:   repo.TextPtr(params.Referencia),
	})
}

func (r *postgresRepository) GetByCuadre(ctx context.Context, cuadreCajaID int64) ([]generated.GetConsignacionesByCuadreRow, error) {
	return r.q.GetConsignacionesByCuadre(ctx, cuadreCajaID)
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.Consignacione, error) {
	row, err := r.q.GetConsignacionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Consignacione{}, ErrNotFound
		}
		return generated.Consignacione{}, err
	}
	return row, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	return r.q.DeleteConsignacion(ctx, id)
}

func (r *postgresRepository) GetTotalByCuadre(ctx context.Context, cuadreCajaID int64) (decimal.Decimal, error) {
	return r.q.GetTotalConsignacionesByCuadre(ctx, cuadreCajaID)
}
