package movimientos_test

import (
	"context"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
)

func TestEntradaMercanciaBatch3Items(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	productoA := seedProducto(t, env, "Vela de Vainilla")
	productoB := seedProducto(t, env, "Difusor")

	result, err := env.service.EntradaMercancia(ctx, usecase.EntradaMercanciaInput{
		SedeID: env.sedeID,
		Items: []usecase.EntradaMercanciaItem{
			{TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "bodega", Cantidad: "500.00"},
			{TipoItem: "producto", ItemID: productoA, Ubicacion: "vitrina", Cantidad: "20"},
			{TipoItem: "producto", ItemID: productoB, Ubicacion: "bodega", Cantidad: "15"},
		},
		RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 movimientos, got %d", len(result))
	}
}

func TestTrasladoGeneraDosMovimientosConMismoTimestamp(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "bodega", "200.00")

	result, err := env.service.Traslado(ctx, usecase.TrasladoInput{
		SedeID: env.sedeID,
		Items: []usecase.TrasladoItem{
			{TipoItem: "fragancia", ItemID: fraganciaID, Cantidad: "100.00"},
		},
		RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 movimientos (salida bodega + entrada vitrina), got %d", len(result))
	}
	if !result[0].CreatedAt.Equal(result[1].CreatedAt) {
		t.Fatalf("expected both movimientos to share the same created_at, got %v and %v", result[0].CreatedAt, result[1].CreatedAt)
	}

	vitrina, bodega, err := env.service.StockActual.GetStockTotal(ctx, env.sedeID, "fragancia", fraganciaID)
	if err != nil {
		t.Fatalf("unexpected error reading stock: %v", err)
	}
	if vitrina.String() != "100" || bodega.String() != "100" {
		t.Fatalf("expected vitrina=100 bodega=100 after traslado, got vitrina=%s bodega=%s", vitrina, bodega)
	}
}

func TestTrasladoStockInsuficiente(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "bodega", "50.00")

	_, err := env.service.Traslado(ctx, usecase.TrasladoInput{
		SedeID: env.sedeID,
		Items: []usecase.TrasladoItem{
			{TipoItem: "fragancia", ItemID: fraganciaID, Cantidad: "100.00"},
		},
		RequesterID: env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestAjusteSube(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "vitrina", "10.00")

	m, err := env.service.Ajuste(ctx, usecase.AjusteInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "vitrina",
		CantidadNueva: "25.00", Motivo: "Conteo físico", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected a movimiento, got nil")
	}
	if m.Cantidad.String() != "15" || m.StockPosterior.String() != "25" {
		t.Fatalf("expected cantidad=15 stock_posterior=25, got cantidad=%s stock_posterior=%s", m.Cantidad, m.StockPosterior)
	}
}

func TestAjusteBaja(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "vitrina", "100.00")

	m, err := env.service.Ajuste(ctx, usecase.AjusteInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "vitrina",
		CantidadNueva: "90.00", Motivo: "Conteo físico", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected a movimiento, got nil")
	}
	if m.Cantidad.String() != "-10" || m.StockPosterior.String() != "90" {
		t.Fatalf("expected cantidad=-10 stock_posterior=90, got cantidad=%s stock_posterior=%s", m.Cantidad, m.StockPosterior)
	}
}

func TestAjusteNoOpCuandoCantidadIgual(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "vitrina", "50.00")

	m, err := env.service.Ajuste(ctx, usecase.AjusteInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "vitrina",
		CantidadNueva: "50.00", Motivo: "Conteo físico", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m != nil {
		t.Fatalf("expected no movimiento when cantidad_nueva matches current stock, got %+v", m)
	}
}

func TestAjusteSinMotivoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)

	_, err := env.service.Ajuste(ctx, usecase.AjusteInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "vitrina",
		CantidadNueva: "10.00", RequesterID: env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestAjusteMotivoEnBlancoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)

	_, err := env.service.Ajuste(ctx, usecase.AjusteInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "vitrina",
		CantidadNueva: "10.00", Motivo: "   ", RequesterID: env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestAjusteCantidadNegativaFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)

	_, err := env.service.Ajuste(ctx, usecase.AjusteInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "vitrina",
		CantidadNueva: "-5.00", Motivo: "Conteo físico", RequesterID: env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestDanado(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "vitrina", "20.00")

	m, err := env.service.Danado(ctx, usecase.DanadoInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "vitrina",
		Cantidad: "3.00", Motivo: "Frasco roto", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Cantidad.String() != "-3" || m.StockPosterior.String() != "17" {
		t.Fatalf("expected cantidad=-3 stock_posterior=17, got cantidad=%s stock_posterior=%s", m.Cantidad, m.StockPosterior)
	}
}

func TestDanadoSinMotivoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "vitrina", "20.00")

	_, err := env.service.Danado(ctx, usecase.DanadoInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "vitrina",
		Cantidad: "3.00", RequesterID: env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestCorreccion(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "bodega", "40.00")

	m, err := env.service.Correccion(ctx, usecase.CorreccionInput{
		SedeID: env.sedeID, TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "bodega",
		CantidadNueva: "35.00", Motivo: "Error de digitación en entrada anterior", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil || m.Cantidad.String() != "-5" {
		t.Fatalf("expected a movimiento with cantidad=-5, got %+v", m)
	}
}

func TestListConFiltros(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	productoID := seedProducto(t, env, "Vela de Vainilla")

	if _, err := env.service.EntradaMercancia(ctx, usecase.EntradaMercanciaInput{
		SedeID: env.sedeID,
		Items: []usecase.EntradaMercanciaItem{
			{TipoItem: "fragancia", ItemID: fraganciaID, Ubicacion: "bodega", Cantidad: "50.00"},
			{TipoItem: "producto", ItemID: productoID, Ubicacion: "vitrina", Cantidad: "5"},
		},
		RequesterID: env.requesterID,
	}); err != nil {
		t.Fatalf("seeding entrada: %v", err)
	}

	t.Run("filtra_por_tipo_item", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, TipoItem: string(domainstock.TipoItemProducto)})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || string(result.Items[0].TipoItem) != string(domainstock.TipoItemProducto) {
			t.Fatalf("expected only the producto movimiento, got total=%d items=%+v", result.Total, result.Items)
		}
	})

	t.Run("filtra_por_item_id", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, TipoItem: string(domainstock.TipoItemFragancia), ItemID: fraganciaID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 {
			t.Fatalf("expected 1 movimiento for this fragancia, got %d", result.Total)
		}
	})

	t.Run("filtra_por_tipo", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, Tipo: "entrada_mercancia"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 2 {
			t.Fatalf("expected 2 entrada_mercancia movimientos, got %d", result.Total)
		}
	})
}
