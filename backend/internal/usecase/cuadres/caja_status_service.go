package cuadres

import (
	"context"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
)

// CajaStatusService answers "can a venta be registered for sede/fecha
// right now?" — consumed by usecase/ventas.CreateVenta to replace the
// M10 TODO(cuadre) placeholder. It intentionally does not participate in
// CreateVenta's own transaction (see usecase/ventas/create.go): it's a
// plain read against the pool, same as the rest of CreateVenta's Fase 1
// pre-validations.
type CajaStatusService interface {
	// VerificarPuedeRegistrarVenta returns nil if fecha has no cuadre yet,
	// or its cuadre is still abierto. Returns a *DomainError (business rule,
	// 422) if fecha's cuadre is already cerrado.
	VerificarPuedeRegistrarVenta(ctx context.Context, sedeID int64, fecha time.Time) error
}

type cajaStatusService struct {
	cuadres cuadresrepo.Repository
}

// NewCajaStatusService builds a CajaStatusService.
func NewCajaStatusService(cuadres cuadresrepo.Repository) CajaStatusService {
	return &cajaStatusService{cuadres: cuadres}
}

func (s *cajaStatusService) VerificarPuedeRegistrarVenta(ctx context.Context, sedeID int64, fecha time.Time) error {
	cerrado, err := s.cuadres.ExistsCerradoBySedeFecha(ctx, sedeID, fecha)
	if err != nil {
		return internalErr(err)
	}
	if cerrado {
		return domainerrors.NewBusinessRule(
			"Cuadre de caja cerrado",
			"El cuadre de caja del día está cerrado. No se pueden registrar más ventas para hoy.",
		)
	}
	return nil
}
