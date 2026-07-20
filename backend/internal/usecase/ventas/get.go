package ventas

import (
	"context"

	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
	movimientosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
)

// maxMovimientosPorVenta is a safe upper bound for how many movimientos one
// venta can generate (well above the 50-item cap, since even the largest
// possible venta produces at most 3 movimientos per line).
const maxMovimientosPorVenta = 200

// GetResult is a venta plus the ids of every movimiento it generated —
// composed here rather than carried on the domain Venta type, since
// "which movimientos this venta produced" is a query-time concern, not an
// intrinsic property of a venta.
type GetResult struct {
	Venta                domainventas.Venta
	MovimientosGenerados []int64
}

// Get loads a single venta with its items and the ids of the movimientos
// it generated. Authorization (a vendedora may only see her own ventas) is
// the HTTP handler's responsibility — this usecase doesn't know about roles.
func (s *Service) Get(ctx context.Context, id int64) (GetResult, error) {
	venta, err := s.loadVentaCompleta(ctx, id)
	if err != nil {
		return GetResult{}, err
	}

	rows, _, err := s.MovimientosRepo.ListPaginated(ctx, movimientosrepo.ListFilter{
		VentaID:  id,
		PageSize: maxMovimientosPorVenta,
	})
	if err != nil {
		return GetResult{}, internalErr(err)
	}

	ids := make([]int64, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}

	return GetResult{Venta: *venta, MovimientosGenerados: ids}, nil
}
