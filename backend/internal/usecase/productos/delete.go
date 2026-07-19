package productos

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

// DeleteInput is the request context: who is deleting what. Delete is
// admin-only, enforced at the router.
type DeleteInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Delete soft-deletes a producto. Requires zero stock everywhere (vitrina +
// bodega) — the caller must adjust it to zero first (Tanda 3 movimientos).
func (s *Service) Delete(ctx context.Context, in DeleteInput) error {
	p, err := s.Productos.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}
	if p.DeletedAt.Valid {
		return notFoundErr()
	}

	vitrina, bodega, err := s.StockActual.GetStockTotal(ctx, p.SedeID, stockactual.TipoItemProducto, p.ID)
	if err != nil {
		return internalErr(err)
	}
	if !vitrina.Add(bodega).IsZero() {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"No se puede eliminar un producto con stock. Ajusta el stock a 0 primero.",
		)
	}

	if err := s.Productos.SoftDelete(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "producto_eliminado", in.IP, in.UserAgent, &in.TargetID, snapshot(p), nil)
	return nil
}
