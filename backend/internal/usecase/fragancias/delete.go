package fragancias

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

// DeleteInput is the request context: who is deleting what.
type DeleteInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Delete soft-deletes a fragancia. Requires zero stock everywhere (vitrina +
// bodega) — the caller must adjust it to zero first (Tanda 3 movimientos).
func (s *Service) Delete(ctx context.Context, in DeleteInput) error {
	f, err := s.Fragancias.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}
	if f.DeletedAt.Valid {
		return notFoundErr()
	}

	vitrina, bodega, err := s.StockActual.GetStockTotal(ctx, f.SedeID, stockactual.TipoItemFragancia, f.ID)
	if err != nil {
		return internalErr(err)
	}
	if !vitrina.Add(bodega).IsZero() {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"No se puede eliminar una fragancia con stock. Ajusta el stock a 0 primero.",
		)
	}

	if err := s.Fragancias.SoftDelete(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "fragancia_eliminada", in.IP, in.UserAgent, &in.TargetID, snapshot(f), nil)
	return nil
}
