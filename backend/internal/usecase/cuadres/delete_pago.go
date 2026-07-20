package cuadres

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	pagoscajarepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/pagos_caja"
)

// DeletePago removes a pago from an abierto cuadre — admin only (enforced
// by the HTTP router), and only while the cuadre is still abierto.
func (s *Service) DeletePago(ctx context.Context, cuadreID, pagoID int64) error {
	pago, err := s.Pagos.GetByID(ctx, pagoID)
	if err != nil {
		if errors.Is(err, pagoscajarepo.ErrNotFound) {
			return pagoNotFoundErr()
		}
		return internalErr(err)
	}
	if pago.CuadreCajaID != cuadreID {
		return pagoNotFoundErr()
	}

	cuadre, err := s.Cuadres.GetByID(ctx, cuadreID)
	if err != nil {
		if errors.Is(err, cuadresrepo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}
	if cuadre.Estado != generated.EstadoCuadreEnumAbierto {
		return domainerrors.NewBusinessRule("Cuadre cerrado", "El cuadre de caja está cerrado, no se pueden eliminar pagos.")
	}

	if err := s.Pagos.Delete(ctx, pagoID); err != nil {
		return internalErr(err)
	}
	return nil
}
