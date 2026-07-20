package ventas_test

import (
	"testing"

	"github.com/shopspring/decimal"

	usecaseventas "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/ventas"
)

func TestDiscountSinDescuentoBajoUmbral(t *testing.T) {
	svc := usecaseventas.NewDiscountService()
	result := svc.Apply(decimal.RequireFromString("49999"))
	if !result.Pct.IsZero() || !result.Monto.IsZero() || result.Total.String() != "49999" {
		t.Fatalf("expected no discount below 50000, got %+v", result)
	}
}

func TestDiscount5PorcientoEnUmbral(t *testing.T) {
	svc := usecaseventas.NewDiscountService()
	result := svc.Apply(decimal.RequireFromString("50000"))
	if result.Pct.String() != "5" || result.Monto.String() != "2500" || result.Total.String() != "47500" {
		t.Fatalf("expected 5%% (2500) at 50000, got %+v", result)
	}
}

func TestDiscount5PorcientoJustoBajo100000(t *testing.T) {
	svc := usecaseventas.NewDiscountService()
	result := svc.Apply(decimal.RequireFromString("99999"))
	if result.Pct.String() != "5" {
		t.Fatalf("expected 5%% at 99999, got %+v", result)
	}
}

func TestDiscount7PorcientoEnUmbral(t *testing.T) {
	svc := usecaseventas.NewDiscountService()
	result := svc.Apply(decimal.RequireFromString("100000"))
	if result.Pct.String() != "7" || result.Monto.String() != "7000" || result.Total.String() != "93000" {
		t.Fatalf("expected 7%% (7000) at 100000, got %+v", result)
	}
}

func TestDiscount7PorcientoMuySobreUmbral(t *testing.T) {
	svc := usecaseventas.NewDiscountService()
	result := svc.Apply(decimal.RequireFromString("250000"))
	if result.Pct.String() != "7" || result.Monto.String() != "17500" {
		t.Fatalf("expected 7%% (17500) at 250000, got %+v", result)
	}
}

func TestDiscountRedondeoHalfAwayFromZero(t *testing.T) {
	svc := usecaseventas.NewDiscountService()
	// 55555 * 5% = 2777.75 -> half-away-from-zero rounds to 2778.
	result := svc.Apply(decimal.RequireFromString("55555"))
	if result.Monto.String() != "2778" {
		t.Fatalf("expected monto=2778 (round-half-away-from-zero of 2777.75), got %s", result.Monto)
	}
}
