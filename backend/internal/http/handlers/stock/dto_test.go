package stock

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

func TestToStockItemResponseUbicacionTransform(t *testing.T) {
	item := generated.ListStockUnificadoRow{
		TipoItem:     "fragancia",
		ItemID:       1,
		Nombre:       "Bleu de Chanel",
		StockVitrina: decimal.RequireFromString("30.00"),
		StockBodega:  decimal.RequireFromString("70.00"),
		StockTotal:   decimal.RequireFromString("100.00"),
		Minimo:       decimal.RequireFromString("10.00"),
		BajoMinimo:   false,
		Unidad:       "gramos",
	}

	t.Run("sin filtro devuelve vitrina y bodega reales", func(t *testing.T) {
		got := toStockItemResponse(item, "")
		if got.StockVitrina != "30" || got.StockBodega != "70" || got.StockTotal != "100" {
			t.Fatalf("expected real stock unchanged, got %+v", got)
		}
	})

	t.Run("ubicacion=vitrina oculta bodega y recalcula total", func(t *testing.T) {
		got := toStockItemResponse(item, "vitrina")
		if got.StockVitrina != "30" || got.StockBodega != "0" || got.StockTotal != "30" {
			t.Fatalf("expected bodega zeroed and total=vitrina, got %+v", got)
		}
	})

	t.Run("ubicacion=bodega oculta vitrina y recalcula total", func(t *testing.T) {
		got := toStockItemResponse(item, "bodega")
		if got.StockBodega != "70" || got.StockVitrina != "0" || got.StockTotal != "70" {
			t.Fatalf("expected vitrina zeroed and total=bodega, got %+v", got)
		}
	})

	t.Run("bajo_minimo no se recalcula con el filtro de ubicacion", func(t *testing.T) {
		bajo := item
		bajo.BajoMinimo = true
		got := toStockItemResponse(bajo, "vitrina")
		if !got.BajoMinimo {
			t.Fatalf("expected bajo_minimo to keep reflecting the real total, got %+v", got)
		}
	})
}
