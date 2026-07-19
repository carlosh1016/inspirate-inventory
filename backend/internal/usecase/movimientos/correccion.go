package movimientos

import (
	"context"

	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
)

// CorreccionInput has the same shape as AjusteInput — corrección is used
// when no other movimiento type applies, but follows identical rules
// (absolute target, motivo required, admin-only).
type CorreccionInput = AjusteInput

// Correccion sets the absolute stock quantity for one item/ubicación, same
// semantics as Ajuste (see ajustarCantidad), logged under its own
// auditoria acción.
func (s *Service) Correccion(ctx context.Context, in CorreccionInput) (*domainmovimientos.Movimiento, error) {
	return s.ajustarCantidad(ctx, in, domainmovimientos.TipoCorreccion, "correccion_inventario")
}
