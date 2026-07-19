package variantesenvase

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

// DeleteInput is the request context: who is deleting what.
type DeleteInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Delete soft-deletes a variante_envase. Requires zero stock everywhere
// (vitrina + bodega) — the caller must adjust it to zero first (Tanda 3
// movimientos).
func (s *Service) Delete(ctx context.Context, in DeleteInput) error {
	v, err := s.VariantesEnvase.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}
	if v.DeletedAt.Valid {
		return notFoundErr()
	}

	vitrina, bodega, err := s.StockActual.GetStockTotal(ctx, v.SedeID, stockactual.TipoItemVarianteEnvase, v.ID)
	if err != nil {
		return internalErr(err)
	}
	if !vitrina.Add(bodega).IsZero() {
		return domainerrors.NewBusinessRule(
			"Operación no permitida",
			"No se puede eliminar una variante de envase con stock. Ajusta el stock a 0 primero.",
		)
	}

	if err := s.VariantesEnvase.SoftDelete(ctx, in.TargetID); err != nil {
		return internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "variante_envase_eliminada", in.IP, in.UserAgent, &in.TargetID, snapshot(v), nil)
	return nil
}
