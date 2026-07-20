package ventas

import "github.com/shopspring/decimal"

// DiscountResult is the outcome of applying the automatic, non-cumulative
// discount rule to a venta's subtotal.
type DiscountResult struct {
	Pct   decimal.Decimal
	Monto decimal.Decimal
	Total decimal.Decimal
}

// DiscountService applies the automatic discount rule to a venta's
// subtotal. Pure: no I/O. The result is persisted as-is on the venta and
// never recalculated later — changing the rule tomorrow doesn't alter
// historical ventas.
type DiscountService interface {
	Apply(subtotal decimal.Decimal) DiscountResult
}
