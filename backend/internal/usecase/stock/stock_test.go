package stock_test

import (
	"context"
	"testing"

	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/stock"
)

func TestListNoFiltersCombinesAllThreeTypes(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	seedFragancia(t, env, "Bleu de Chanel", true)
	modelo := seedModeloEnvase(t, env, "Spray", "3.00")
	seedVarianteEnvase(t, env, modelo, "Rojo", 5, true)
	seedProducto(t, env, "Vela de Vainilla", 5, true)

	result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 3 {
		t.Fatalf("expected 3 items across all types, got %d", result.Total)
	}

	seen := map[string]bool{}
	for _, item := range result.Items {
		seen[item.TipoItem] = true
	}
	for _, tipo := range []string{"fragancia", "variante_envase", "producto"} {
		if !seen[tipo] {
			t.Errorf("expected %s to be present in the unified list", tipo)
		}
	}
}

func TestListFiltraPorTipoItem(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	seedFragancia(t, env, "Bleu de Chanel", true)
	seedProducto(t, env, "Vela de Vainilla", 5, true)

	result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, TipoItem: "producto"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 || len(result.Items) != 1 || result.Items[0].TipoItem != "producto" {
		t.Fatalf("expected only producto, got total=%d items=%+v", result.Total, result.Items)
	}
}

func TestListFiltraPorStockBajo(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	bajo := seedFragancia(t, env, "Bleu de Chanel", true)
	sobre := seedFragancia(t, env, "Miss Dior", true)
	setVitrinaStock(t, env, "fragancia", bajo, "2.00")
	setVitrinaStock(t, env, "fragancia", sobre, "20.00")

	result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, StockBajo: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ItemID != bajo {
		t.Fatalf("expected only the fragancia under gramos_minimo, got total=%d items=%+v", result.Total, result.Items)
	}
}

func TestListFiltraPorStockCero(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	cero := seedFragancia(t, env, "Bleu de Chanel", true)
	conStock := seedFragancia(t, env, "Miss Dior", true)
	setVitrinaStock(t, env, "fragancia", conStock, "5.00")

	result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, StockCero: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ItemID != cero {
		t.Fatalf("expected only the fragancia with zero stock, got total=%d items=%+v", result.Total, result.Items)
	}
}

func TestListPaginacion(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	for i := 0; i < 12; i++ {
		seedProducto(t, env, "Producto "+string(rune('A'+i)), 0, true)
	}

	page2, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, Page: 2, PageSize: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page2.Total != 12 {
		t.Fatalf("expected total=12, got %d", page2.Total)
	}
	if len(page2.Items) != 5 {
		t.Fatalf("expected 5 items on page 2, got %d", len(page2.Items))
	}
	if page2.Page != 2 || page2.PageSize != 5 {
		t.Fatalf("expected page=2 page_size=5, got page=%d page_size=%d", page2.Page, page2.PageSize)
	}
}

func TestListSoftDeletedNoAparecen(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	id := seedFragancia(t, env, "Bleu de Chanel", true)
	if _, err := env.pool.Exec(ctx, `UPDATE fragancias SET deleted_at = NOW() WHERE id = $1`, id); err != nil {
		t.Fatalf("soft-deleting fragancia: %v", err)
	}

	result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 {
		t.Fatalf("expected soft-deleted fragancia to be excluded, got total=%d", result.Total)
	}
}

func TestListInactivosSoloConIncludeInactivos(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	seedFragancia(t, env, "Bleu de Chanel", false)

	hidden, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hidden.Total != 0 {
		t.Fatalf("expected inactive fragancia hidden by default, got total=%d", hidden.Total)
	}

	shown, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, IncludeInactivos: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shown.Total != 1 {
		t.Fatalf("expected inactive fragancia visible with include_inactivos, got total=%d", shown.Total)
	}
}
