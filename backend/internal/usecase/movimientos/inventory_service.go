package movimientos

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
	fraganciasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	movimientosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
	productosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	stockactualrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	variantesenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

// InventoryService is the transactional engine every movimiento goes
// through — Tanda 4's ventas usecase will import this package and call
// RegisterBatchTx directly to apply stock changes as part of registering a
// sale, inside its own transaction.
type InventoryService interface {
	// RegisterBatch registers one or more movimientos in a single new
	// transaction (RepeatableRead + row-level locking on stock_actual), all
	// or nothing.
	RegisterBatch(ctx context.Context, sedeID, usuarioID int64, inputs []domainmovimientos.MovimientoInput) ([]domainmovimientos.Movimiento, error)
	// RegisterBatchTx is the same engine running inside a transaction the
	// caller already owns (used by ventas in Tanda 4, and internally by
	// RegisterBatch).
	RegisterBatchTx(ctx context.Context, tx pgx.Tx, sedeID, usuarioID int64, inputs []domainmovimientos.MovimientoInput) ([]domainmovimientos.Movimiento, error)
}

// RegisterBatch opens its own RepeatableRead transaction and delegates to
// RegisterBatchTx.
func (s *Service) RegisterBatch(ctx context.Context, sedeID, usuarioID int64, inputs []domainmovimientos.MovimientoInput) ([]domainmovimientos.Movimiento, error) {
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return nil, internalErr(err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // no-op once committed
	}()

	result, err := s.RegisterBatchTx(ctx, tx, sedeID, usuarioID, inputs)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, internalErr(err)
	}
	return result, nil
}

// RegisterBatchTx applies every input sequentially within tx: for each one
// it locks the target stock_actual row (SELECT ... FOR UPDATE), computes
// the resulting quantity, and — if it would go negative — records the
// shortfall and keeps checking the rest of the batch instead of stopping,
// so the caller gets every offending item in one response. If any item is
// insufficient, the whole batch is reported as failed (the caller's
// transaction rollback, in RegisterBatch or in the caller for
// RegisterBatchTx, discards whatever was already applied in this loop —
// the batch is all-or-nothing). Repeated (tipo_item, item, ubicacion) pairs
// within the same batch chain correctly: FOR UPDATE + an immediate write
// per input means a later input targeting the same row sees the earlier
// input's effect, not the pre-batch snapshot.
func (s *Service) RegisterBatchTx(ctx context.Context, tx pgx.Tx, sedeID, usuarioID int64, inputs []domainmovimientos.MovimientoInput) ([]domainmovimientos.Movimiento, error) {
	if len(inputs) == 0 {
		return nil, domainerrors.NewValidation("Solicitud inválida", "Debes incluir al menos un ítem.", nil)
	}

	txStockActual := stockactualrepo.NewPostgres(tx)
	txMovimientos := movimientosrepo.NewPostgres(tx)

	var insuficientes []domainstock.StockInsuficienteItem
	result := make([]domainmovimientos.Movimiento, 0, len(inputs))

	for _, in := range inputs {
		nombre, err := s.resolveItemNombre(ctx, tx, in.TipoItem, in.ItemID)
		if err != nil {
			return nil, err
		}

		row, _, err := txStockActual.GetForUpdate(ctx, sedeID, string(in.TipoItem), in.ItemID, string(in.Ubicacion))
		if err != nil {
			return nil, internalErr(err)
		}

		stockAnterior := row.Cantidad
		stockPosterior := stockAnterior.Add(in.Cantidad)
		if stockPosterior.IsNegative() {
			insuficientes = append(insuficientes, domainstock.StockInsuficienteItem{
				TipoItem:   string(in.TipoItem),
				ItemID:     in.ItemID,
				Nombre:     nombre,
				Requerido:  in.Cantidad.Abs().String(),
				Disponible: stockAnterior.String(),
			})
			continue
		}

		if err := txStockActual.SetCantidad(ctx, sedeID, string(in.TipoItem), in.ItemID, string(in.Ubicacion), stockPosterior); err != nil {
			return nil, internalErr(err)
		}

		movRow, err := txMovimientos.Insert(ctx, movimientosrepo.InsertParams{
			SedeID:         sedeID,
			UsuarioID:      usuarioID,
			TipoItem:       string(in.TipoItem),
			ItemID:         in.ItemID,
			Tipo:           string(in.Tipo),
			Ubicacion:      string(in.Ubicacion),
			Cantidad:       in.Cantidad,
			StockAnterior:  stockAnterior,
			StockPosterior: stockPosterior,
			Motivo:         in.Motivo,
			VentaID:        in.VentaID,
		})
		if err != nil {
			return nil, internalErr(err)
		}

		result = append(result, toDomainMovimiento(movRow, nombre))
	}

	if len(insuficientes) > 0 {
		err := domainerrors.NewBusinessRule(
			"Stock insuficiente",
			"No hay suficiente stock disponible para uno o más ítems.",
		)
		err.Extra = insuficientes
		return nil, err
	}

	return result, nil
}

// resolveItemNombre validates that tipoItem/itemID exists and isn't
// soft-deleted, returning a short display name for it. Runs against tx so
// it sees the same snapshot as the rest of the batch.
func (s *Service) resolveItemNombre(ctx context.Context, tx pgx.Tx, tipoItem domainstock.TipoItem, itemID int64) (string, error) {
	switch tipoItem {
	case domainstock.TipoItemFragancia:
		f, err := fraganciasrepo.NewPostgres(tx).GetByIDIncludingDeleted(ctx, itemID)
		if err != nil {
			if errors.Is(err, fraganciasrepo.ErrNotFound) {
				return "", itemNoEncontradoErr(tipoItem, itemID)
			}
			return "", internalErr(err)
		}
		if f.DeletedAt.Valid {
			return "", itemNoEncontradoErr(tipoItem, itemID)
		}
		return f.NombreComercial, nil

	case domainstock.TipoItemVarianteEnvase:
		v, err := variantesenvaserepo.NewPostgres(tx).GetByIDIncludingDeleted(ctx, itemID)
		if err != nil {
			if errors.Is(err, variantesenvaserepo.ErrNotFound) {
				return "", itemNoEncontradoErr(tipoItem, itemID)
			}
			return "", internalErr(err)
		}
		if v.DeletedAt.Valid {
			return "", itemNoEncontradoErr(tipoItem, itemID)
		}
		return v.Color, nil

	case domainstock.TipoItemProducto:
		p, err := productosrepo.NewPostgres(tx).GetByIDIncludingDeleted(ctx, itemID)
		if err != nil {
			if errors.Is(err, productosrepo.ErrNotFound) {
				return "", itemNoEncontradoErr(tipoItem, itemID)
			}
			return "", internalErr(err)
		}
		if p.DeletedAt.Valid {
			return "", itemNoEncontradoErr(tipoItem, itemID)
		}
		return p.Nombre, nil

	default:
		return "", domainerrors.NewValidation("Solicitud inválida", "tipo_item debe ser fragancia, variante_envase o producto.", nil)
	}
}
