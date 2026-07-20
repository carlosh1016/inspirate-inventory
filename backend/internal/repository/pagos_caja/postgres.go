package pagoscaja

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) Insert(ctx context.Context, params InsertParams) (generated.PagosCaja, error) {
	return r.q.InsertPagoCaja(ctx, generated.InsertPagoCajaParams{
		CuadreCajaID: params.CuadreCajaID,
		UsuarioID:    params.UsuarioID,
		Concepto:     params.Concepto,
		Monto:        params.Monto,
	})
}

func (r *postgresRepository) GetByCuadre(ctx context.Context, cuadreCajaID int64) ([]generated.GetPagosByCuadreRow, error) {
	return r.q.GetPagosByCuadre(ctx, cuadreCajaID)
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.PagosCaja, error) {
	row, err := r.q.GetPagoByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.PagosCaja{}, ErrNotFound
		}
		return generated.PagosCaja{}, err
	}
	return row, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	return r.q.DeletePagoCaja(ctx, id)
}

func (r *postgresRepository) GetTotalByCuadre(ctx context.Context, cuadreCajaID int64) (decimal.Decimal, error) {
	return r.q.GetTotalPagosByCuadre(ctx, cuadreCajaID)
}
