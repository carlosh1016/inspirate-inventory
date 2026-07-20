package cuadres

import (
	"context"
	"time"

	"github.com/shopspring/decimal"

	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
)

// TotalesPorMetodo is ventas totals for one día, broken down by payment
// method code (efectivo/nequi/daviplata/transferencia/otros) plus the
// grand total and count.
type TotalesPorMetodo struct {
	Efectivo      decimal.Decimal
	Nequi         decimal.Decimal
	Daviplata     decimal.Decimal
	Transferencia decimal.Decimal
	Otros         decimal.Decimal
	Total         decimal.Decimal
	VentasCount   int64
}

// TotalesService computes live ventas totals for a given día, used both by
// GetHoy (while the cuadre is abierto) and by Cerrar (to freeze the
// snapshot).
type TotalesService interface {
	CalcularParaFecha(ctx context.Context, sedeID int64, fecha time.Time) (TotalesPorMetodo, error)
}

type totalesService struct {
	cuadres  cuadresrepo.Repository
	location *time.Location
}

// NewTotalesService builds a TotalesService. loc is the timezone used to
// bound the [diaInicio, diaFin) window for fecha.
func NewTotalesService(cuadres cuadresrepo.Repository, loc *time.Location) TotalesService {
	return &totalesService{cuadres: cuadres, location: loc}
}

func (s *totalesService) CalcularParaFecha(ctx context.Context, sedeID int64, fecha time.Time) (TotalesPorMetodo, error) {
	diaInicio := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, s.location)
	diaFin := diaInicio.Add(24 * time.Hour)

	row, err := s.cuadres.TotalesPorMetodoEnFecha(ctx, sedeID, diaInicio, diaFin)
	if err != nil {
		return TotalesPorMetodo{}, err
	}

	return TotalesPorMetodo{
		Efectivo:      row.TotalEfectivo,
		Nequi:         row.TotalNequi,
		Daviplata:     row.TotalDaviplata,
		Transferencia: row.TotalTransferencia,
		Otros:         row.TotalOtros,
		Total:         row.TotalDia,
		VentasCount:   row.VentasCount,
	}, nil
}
