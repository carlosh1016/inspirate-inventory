package ventas_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/envases"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/productos"
	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
	usecaseventas "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/ventas"
)

func testModelo() *envases.ModeloEnvase {
	return &envases.ModeloEnvase{
		ID:                 1,
		PrecioSolo:         decimal.RequireFromString("10000"),
		PrecioConFragancia: decimal.RequireFromString("25000"),
		PrecioRecarga:      decimal.RequireFromString("15000"),
	}
}

func testProducto(precio string) *productos.Producto {
	return &productos.Producto{ID: 1, Precio: decimal.RequireFromString(precio)}
}

func TestPricingEnvaseConFragancia(t *testing.T) {
	svc := usecaseventas.NewPricingService()
	result, err := svc.Calculate(domainventas.PricingInput{
		TipoLinea: domainventas.TipoLineaEnvaseConFragancia, ModeloEnvase: testModelo(), Cantidad: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PrecioUnitario.String() != "25000" || result.SubtotalLinea.String() != "25000" {
		t.Fatalf("expected precio_unitario=25000 subtotal=25000, got %+v", result)
	}
}

func TestPricingRecarga(t *testing.T) {
	svc := usecaseventas.NewPricingService()
	result, err := svc.Calculate(domainventas.PricingInput{
		TipoLinea: domainventas.TipoLineaRecarga, ModeloEnvase: testModelo(), Cantidad: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PrecioUnitario.String() != "15000" {
		t.Fatalf("expected precio_recarga=15000, got %s", result.PrecioUnitario)
	}
}

func TestPricingEnvaseSolo(t *testing.T) {
	svc := usecaseventas.NewPricingService()
	result, err := svc.Calculate(domainventas.PricingInput{
		TipoLinea: domainventas.TipoLineaEnvaseSolo, ModeloEnvase: testModelo(), Cantidad: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PrecioUnitario.String() != "10000" {
		t.Fatalf("expected precio_solo=10000, got %s", result.PrecioUnitario)
	}
}

func TestPricingProductoOtro(t *testing.T) {
	svc := usecaseventas.NewPricingService()
	result, err := svc.Calculate(domainventas.PricingInput{
		TipoLinea: domainventas.TipoLineaProductoOtro, Producto: testProducto("18000"), Cantidad: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PrecioUnitario.String() != "18000" || result.SubtotalLinea.String() != "18000" {
		t.Fatalf("expected precio=18000, got %+v", result)
	}
}

func TestPricingConFeromona(t *testing.T) {
	svc := usecaseventas.NewPricingService()
	result, err := svc.Calculate(domainventas.PricingInput{
		TipoLinea: domainventas.TipoLineaEnvaseConFragancia, ModeloEnvase: testModelo(),
		FeromonaProducto: testProducto("1000"), Cantidad: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// subtotal = precio_unitario*cantidad + precio_feromona*cantidad = 25000 + 1000 = 26000
	if result.PrecioFeromona.String() != "1000" || result.SubtotalLinea.String() != "26000" {
		t.Fatalf("expected precio_feromona=1000 subtotal=26000, got %+v", result)
	}
}

func TestPricingSinFeromonaNoSuma(t *testing.T) {
	svc := usecaseventas.NewPricingService()
	result, err := svc.Calculate(domainventas.PricingInput{
		TipoLinea: domainventas.TipoLineaEnvaseConFragancia, ModeloEnvase: testModelo(), Cantidad: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.PrecioFeromona.IsZero() || result.SubtotalLinea.String() != "25000" {
		t.Fatalf("expected no feromona addition, got %+v", result)
	}
}

func TestPricingCantidadMultiple(t *testing.T) {
	svc := usecaseventas.NewPricingService()
	result, err := svc.Calculate(domainventas.PricingInput{
		TipoLinea: domainventas.TipoLineaEnvaseConFragancia, ModeloEnvase: testModelo(),
		FeromonaProducto: testProducto("1000"), Cantidad: 2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// subtotal = precio_unitario*cantidad + precio_feromona*cantidad = 25000*2 + 1000*2 = 52000
	if result.SubtotalLinea.String() != "52000" {
		t.Fatalf("expected subtotal=52000 (precio_unitario*cantidad + feromona*cantidad), got %s", result.SubtotalLinea)
	}
}
