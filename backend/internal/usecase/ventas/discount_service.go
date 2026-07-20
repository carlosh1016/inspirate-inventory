package ventas

import (
	"github.com/shopspring/decimal"

	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
)

var (
	umbralDescuento7 = decimal.NewFromInt(100000)
	umbralDescuento5 = decimal.NewFromInt(50000)
	cien             = decimal.NewFromInt(100)
)

type discountService struct{}

// NewDiscountService builds the default DiscountService: automatic,
// non-cumulative discount — 7% at or above 100.000, 5% at or above 50.000,
// none below that. The amount is rounded to 0 decimal places (COP has no
// cents) using round-half-away-from-zero, shopspring/decimal's Round()
// default (e.g. 2777.75 -> 2778) — confirmed with the user, not banker's
// rounding.
func NewDiscountService() domainventas.DiscountService {
	return discountService{}
}

func (discountService) Apply(subtotal decimal.Decimal) domainventas.DiscountResult {
	var pct decimal.Decimal
	switch {
	case subtotal.GreaterThanOrEqual(umbralDescuento7):
		pct = decimal.NewFromInt(7)
	case subtotal.GreaterThanOrEqual(umbralDescuento5):
		pct = decimal.NewFromInt(5)
	default:
		pct = decimal.Zero
	}

	monto := subtotal.Mul(pct).Div(cien).Round(0)
	total := subtotal.Sub(monto)

	return domainventas.DiscountResult{Pct: pct, Monto: monto, Total: total}
}
