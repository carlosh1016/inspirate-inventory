package fragancias_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/fragancias"
)

func TestStockTotal(t *testing.T) {
	f := fragancias.Fragancia{
		StockVitrina: decimal.NewFromInt(10),
		StockBodega:  decimal.NewFromInt(5),
	}
	if got := f.StockTotal(); !got.Equal(decimal.NewFromInt(15)) {
		t.Fatalf("expected 15, got %s", got)
	}
}

func TestPuedeEliminarse(t *testing.T) {
	sinStock := fragancias.Fragancia{StockVitrina: decimal.Zero, StockBodega: decimal.Zero}
	if !sinStock.PuedeEliminarse() {
		t.Error("expected a fragancia with zero stock to be deletable")
	}

	conStock := fragancias.Fragancia{StockVitrina: decimal.NewFromInt(1), StockBodega: decimal.Zero}
	if conStock.PuedeEliminarse() {
		t.Error("expected a fragancia with stock to not be deletable")
	}
}
