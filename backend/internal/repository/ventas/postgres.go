package ventas

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

const defaultPageSize = 20

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx. db may be
// a *pgxpool.Pool or a pgx.Tx, so Insert can run inside CreateVenta's own
// transaction.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.ListVentasPaginatedRow, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListVentasPaginated(ctx, generated.ListVentasPaginatedParams{
		Limit:        int32(pageSize),
		Offset:       int32(offset),
		SedeID:       filter.SedeID,
		UsuarioID:    filter.UsuarioID,
		MetodoPagoID: filter.MetodoPagoID,
		FechaDesde:   repo.TimestamptzPtr(filter.FechaDesde),
		FechaHasta:   repo.TimestamptzPtr(filter.FechaHasta),
		TotalMin:     filter.TotalMin,
		TotalMax:     filter.TotalMax,
		ConDescuento: filter.ConDescuento,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountVentas(ctx, generated.CountVentasParams{
		SedeID:       filter.SedeID,
		UsuarioID:    filter.UsuarioID,
		MetodoPagoID: filter.MetodoPagoID,
		FechaDesde:   repo.TimestamptzPtr(filter.FechaDesde),
		FechaHasta:   repo.TimestamptzPtr(filter.FechaHasta),
		TotalMin:     filter.TotalMin,
		TotalMax:     filter.TotalMax,
		ConDescuento: filter.ConDescuento,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.GetVentaByIDRow, error) {
	row, err := r.q.GetVentaByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.GetVentaByIDRow{}, ErrNotFound
		}
		return generated.GetVentaByIDRow{}, err
	}
	return row, nil
}

func (r *postgresRepository) Insert(ctx context.Context, params InsertParams) (generated.Venta, error) {
	return r.q.InsertVenta(ctx, generated.InsertVentaParams{
		SedeID:         params.SedeID,
		UsuarioID:      params.UsuarioID,
		MetodoPagoID:   params.MetodoPagoID,
		Subtotal:       params.Subtotal,
		DescuentoPct:   params.DescuentoPct,
		DescuentoMonto: params.DescuentoMonto,
		Total:          params.Total,
		Observaciones:  repo.TextPtr(params.Observaciones),
	})
}

func (r *postgresRepository) UpdateObservaciones(ctx context.Context, id int64, observaciones *string) (generated.Venta, error) {
	v, err := r.q.UpdateVentaObservaciones(ctx, generated.UpdateVentaObservacionesParams{
		ID:            id,
		Observaciones: repo.TextPtr(observaciones),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Venta{}, ErrNotFound
		}
		return generated.Venta{}, err
	}
	return v, nil
}

func (r *postgresRepository) ResumenHoy(ctx context.Context, sedeID int64, diaInicio, diaFin time.Time) (generated.GetResumenVentasHoyRow, error) {
	return r.q.GetResumenVentasHoy(ctx, generated.GetResumenVentasHoyParams{
		SedeID:    sedeID,
		DiaInicio: repo.Timestamptz(diaInicio),
		DiaFin:    repo.Timestamptz(diaFin),
	})
}

func (r *postgresRepository) VentasPorVendedoraHoy(ctx context.Context, sedeID int64, diaInicio, diaFin time.Time) ([]generated.GetVentasPorVendedoraHoyRow, error) {
	return r.q.GetVentasPorVendedoraHoy(ctx, generated.GetVentasPorVendedoraHoyParams{
		SedeID:    sedeID,
		DiaInicio: repo.Timestamptz(diaInicio),
		DiaFin:    repo.Timestamptz(diaFin),
	})
}

func (r *postgresRepository) TopFraganciasHoy(ctx context.Context, sedeID int64, diaInicio, diaFin time.Time) ([]generated.GetTopFraganciasHoyRow, error) {
	return r.q.GetTopFraganciasHoy(ctx, generated.GetTopFraganciasHoyParams{
		SedeID:    sedeID,
		DiaInicio: repo.Timestamptz(diaInicio),
		DiaFin:    repo.Timestamptz(diaFin),
	})
}
