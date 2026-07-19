package movimientos_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
)

func assertCode(t *testing.T, err error, code domainerrors.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var domainErr *domainerrors.DomainError
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected *DomainError, got %T: %v", err, err)
	}
	if domainErr.Code != code {
		t.Fatalf("expected code %q, got %q (%v)", code, domainErr.Code, domainErr)
	}
}

func TestRegisterBatchEntradaSimple(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)

	result, err := env.service.RegisterBatch(ctx, env.sedeID, env.requesterID, []domainmovimientos.MovimientoInput{
		{
			TipoItem:  domainstock.TipoItemFragancia,
			ItemID:    fraganciaID,
			Tipo:      domainmovimientos.TipoEntradaMercancia,
			Ubicacion: domainstock.UbicacionBodega,
			Cantidad:  decimal.RequireFromString("50.00"),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 movimiento, got %d", len(result))
	}
	m := result[0]
	if !m.StockAnterior.IsZero() || m.StockPosterior.String() != "50" {
		t.Fatalf("expected stock_anterior=0 stock_posterior=50, got anterior=%s posterior=%s", m.StockAnterior, m.StockPosterior)
	}
	if m.ItemNombre != "Bleu de Chanel" {
		t.Fatalf("expected item nombre resolved, got %q", m.ItemNombre)
	}

	_, bodega, err := env.service.StockActual.GetStockTotal(ctx, env.sedeID, string(domainstock.TipoItemFragancia), fraganciaID)
	if err != nil {
		t.Fatalf("unexpected error reading stock: %v", err)
	}
	if bodega.String() != "50" {
		t.Fatalf("expected bodega=50 after entrada, got %s", bodega)
	}
}

func TestRegisterBatchEntradaMultiple(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	productoID := seedProducto(t, env, "Vela de Vainilla")

	result, err := env.service.RegisterBatch(ctx, env.sedeID, env.requesterID, []domainmovimientos.MovimientoInput{
		{TipoItem: domainstock.TipoItemFragancia, ItemID: fraganciaID, Tipo: domainmovimientos.TipoEntradaMercancia, Ubicacion: domainstock.UbicacionBodega, Cantidad: decimal.RequireFromString("50.00")},
		{TipoItem: domainstock.TipoItemProducto, ItemID: productoID, Tipo: domainmovimientos.TipoEntradaMercancia, Ubicacion: domainstock.UbicacionVitrina, Cantidad: decimal.RequireFromString("10")},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 movimientos, got %d", len(result))
	}
}

func TestRegisterBatchSalidaInsuficienteDevuelveExtraYNoAplicaNada(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	setStock(t, env, fraganciaID, "vitrina", "5.00")

	_, err := env.service.RegisterBatch(ctx, env.sedeID, env.requesterID, []domainmovimientos.MovimientoInput{
		{
			TipoItem:  domainstock.TipoItemFragancia,
			ItemID:    fraganciaID,
			Tipo:      domainmovimientos.TipoDanado,
			Ubicacion: domainstock.UbicacionVitrina,
			Cantidad:  decimal.RequireFromString("-20.00"),
		},
	})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)

	var domainErr *domainerrors.DomainError
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected *DomainError, got %T", err)
	}
	items, ok := domainErr.Extra.([]domainstock.StockInsuficienteItem)
	if !ok || len(items) != 1 {
		t.Fatalf("expected Extra to carry one StockInsuficienteItem, got %#v", domainErr.Extra)
	}
	if items[0].ItemID != fraganciaID || items[0].Disponible != "5" || items[0].Requerido != "20" {
		t.Fatalf("unexpected StockInsuficienteItem: %+v", items[0])
	}

	vitrina, _, err := env.service.StockActual.GetStockTotal(ctx, env.sedeID, string(domainstock.TipoItemFragancia), fraganciaID)
	if err != nil {
		t.Fatalf("unexpected error reading stock: %v", err)
	}
	if vitrina.String() != "5" {
		t.Fatalf("expected stock unchanged at 5 after the failed batch rolled back, got %s", vitrina)
	}
}

func TestRegisterBatchItemInexistenteFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_, err := env.service.RegisterBatch(ctx, env.sedeID, env.requesterID, []domainmovimientos.MovimientoInput{
		{TipoItem: domainstock.TipoItemFragancia, ItemID: 999999, Tipo: domainmovimientos.TipoEntradaMercancia, Ubicacion: domainstock.UbicacionBodega, Cantidad: decimal.RequireFromString("10")},
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestRegisterBatchItemSoftDeletedFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env)
	if _, err := env.pool.Exec(ctx, `UPDATE fragancias SET deleted_at = NOW() WHERE id = $1`, fraganciaID); err != nil {
		t.Fatalf("soft-deleting fragancia: %v", err)
	}

	_, err := env.service.RegisterBatch(ctx, env.sedeID, env.requesterID, []domainmovimientos.MovimientoInput{
		{TipoItem: domainstock.TipoItemFragancia, ItemID: fraganciaID, Tipo: domainmovimientos.TipoEntradaMercancia, Ubicacion: domainstock.UbicacionBodega, Cantidad: decimal.RequireFromString("10")},
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}
