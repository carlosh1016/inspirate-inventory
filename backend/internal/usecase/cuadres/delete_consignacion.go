package cuadres

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/consignaciones"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// DeleteConsignacion removes a consignacion from an abierto cuadre — admin
// only (enforced by the HTTP router), and only while the cuadre is still
// abierto.
func (s *Service) DeleteConsignacion(ctx context.Context, cuadreID, consignacionID int64) error {
	consig, err := s.Consignaciones.GetByID(ctx, consignacionID)
	if err != nil {
		if errors.Is(err, consignaciones.ErrNotFound) {
			return consignacionNotFoundErr()
		}
		return internalErr(err)
	}
	if consig.CuadreCajaID != cuadreID {
		return consignacionNotFoundErr()
	}

	cuadre, err := s.Cuadres.GetByID(ctx, cuadreID)
	if err != nil {
		if errors.Is(err, cuadresrepo.ErrNotFound) {
			return notFoundErr()
		}
		return internalErr(err)
	}
	if cuadre.Estado != generated.EstadoCuadreEnumAbierto {
		return domainerrors.NewBusinessRule("Cuadre cerrado", "El cuadre de caja está cerrado, no se pueden eliminar consignaciones.")
	}

	if err := s.Consignaciones.Delete(ctx, consignacionID); err != nil {
		return internalErr(err)
	}
	return nil
}
