package movimientos

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
	stockactualrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

// AjusteInput is the request payload plus the requester's context.
// CantidadNueva is the absolute target quantity for that ubicación, not a
// delta.
type AjusteInput struct {
	SedeID        int64
	TipoItem      string
	ItemID        int64
	Ubicacion     string
	CantidadNueva string
	Motivo        string
	RequesterID   int64
	IP            string
	UserAgent     string
}

// ajusteAuditSnapshot is what gets JSON-encoded into
// datos_antes/datos_despues for ajuste_inventario and correccion_inventario.
type ajusteAuditSnapshot struct {
	StockAnterior  string `json:"stock_anterior,omitempty"`
	StockPosterior string `json:"stock_posterior,omitempty"`
	Motivo         string `json:"motivo,omitempty"`
}

// Ajuste sets the absolute stock quantity for one item/ubicación, requiring
// a motivo. Only the difference from the current quantity is registered as
// a movimiento — if the requested quantity matches what's already there,
// nothing is created and Ajuste returns (nil, nil).
func (s *Service) Ajuste(ctx context.Context, in AjusteInput) (*domainmovimientos.Movimiento, error) {
	return s.ajustarCantidad(ctx, in, domainmovimientos.TipoAjuste, "ajuste_inventario")
}

// ajustarCantidad is the shared implementation behind Ajuste and Correccion
// (identical rules: absolute target, motivo required, admin-only at the
// router). It opens its own transaction, locks the target stock_actual row,
// computes the target diff against the *locked* current value, and only
// then delegates to RegisterBatchTx on the same transaction — this
// guarantees the resulting stock_posterior exactly equals CantidadNueva
// even under concurrent writers, since the lock is held continuously from
// the read here through RegisterBatchTx's own (idempotent) re-read of the
// same row.
func (s *Service) ajustarCantidad(ctx context.Context, in AjusteInput, tipo domainmovimientos.TipoMovimiento, accionAuditoria string) (*domainmovimientos.Movimiento, error) {
	motivo := strings.TrimSpace(in.Motivo)
	if motivo == "" {
		return nil, domainerrors.NewValidation("Motivo requerido", "Debes indicar un motivo.", nil)
	}

	cantidadNueva, err := decimal.NewFromString(in.CantidadNueva)
	if err != nil || cantidadNueva.IsNegative() {
		return nil, domainerrors.NewValidation(
			"Cantidad inválida",
			"La cantidad debe ser un número mayor o igual a cero.",
			nil,
		)
	}

	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return nil, internalErr(err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // no-op once committed
	}()

	txStockActual := stockactualrepo.NewPostgres(tx)
	row, _, err := txStockActual.GetForUpdate(ctx, in.SedeID, in.TipoItem, in.ItemID, in.Ubicacion)
	if err != nil {
		return nil, internalErr(err)
	}

	diferencia := cantidadNueva.Sub(row.Cantidad)
	if diferencia.IsZero() {
		return nil, nil
	}

	results, err := s.RegisterBatchTx(ctx, tx, in.SedeID, in.RequesterID, []domainmovimientos.MovimientoInput{
		{
			TipoItem:  domainstock.TipoItem(in.TipoItem),
			ItemID:    in.ItemID,
			Tipo:      tipo,
			Ubicacion: domainstock.Ubicacion(in.Ubicacion),
			Cantidad:  diferencia,
			Motivo:    &motivo,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, internalErr(err)
	}

	result := results[0]
	s.audit(ctx, &in.RequesterID, accionAuditoria, in.IP, in.UserAgent, &result.ID,
		ajusteAuditSnapshot{StockAnterior: row.Cantidad.String()},
		ajusteAuditSnapshot{StockPosterior: result.StockPosterior.String(), Motivo: motivo},
	)

	return &result, nil
}
