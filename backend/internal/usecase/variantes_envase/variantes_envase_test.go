package variantesenvase_test

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/variantes_envase"
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

func strPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32 { return &i }

func TestCreateSuccessInitializesStock(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env.pool, "Spray", "3.00", true)

	v, err := env.service.Create(ctx, usecase.CreateInput{
		SedeID:         env.sedeID,
		ModeloEnvaseID: modeloID,
		Color:          "Rojo",
		StockMinimo:    5,
		RequesterID:    env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Color != "Rojo" {
		t.Fatalf("expected color Rojo, got %q", v.Color)
	}

	vitrina, bodega, err := env.service.StockActual.GetStockTotal(ctx, env.sedeID, stockactual.TipoItemVarianteEnvase, v.ID)
	if err != nil {
		t.Fatalf("unexpected error reading stock: %v", err)
	}
	if !vitrina.IsZero() || !bodega.IsZero() {
		t.Fatalf("expected zero stock on creation, got vitrina=%s bodega=%s", vitrina, bodega)
	}
}

func TestCreateWithInactiveModeloFails(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env.pool, "Spray", "3.00", false)

	_, err := env.service.Create(ctx, usecase.CreateInput{
		SedeID:         env.sedeID,
		ModeloEnvaseID: modeloID,
		Color:          "Rojo",
		RequesterID:    env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestCreateForModeloSinVariantesFails(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvaseSinVariantes(t, env.pool, "Envase de lujo", "1.00")

	_, err := env.service.Create(ctx, usecase.CreateInput{
		SedeID:         env.sedeID,
		ModeloEnvaseID: modeloID,
		Color:          "Delgado",
		RequesterID:    env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestCreateWithUnknownModeloFails(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID:         env.sedeID,
		ModeloEnvaseID: 999999,
		Color:          "Rojo",
		RequesterID:    env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestCreateDuplicateColorConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env.pool, "Spray", "3.00", true)

	in := usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloID, Color: "Rojo", RequesterID: env.requesterID}
	if _, err := env.service.Create(ctx, in); err != nil {
		t.Fatalf("seeding first variante: %v", err)
	}

	_, err := env.service.Create(ctx, in)
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestCreateDuplicateColorCaseInsensitive(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env.pool, "Spray", "3.00", true)

	if _, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloID, Color: "Rojo", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding first variante: %v", err)
	}

	_, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloID, Color: "ROJO", RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestUpdateColorCollisionConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env.pool, "Spray", "3.00", true)

	rojo, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloID, Color: "Rojo", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding rojo: %v", err)
	}
	if _, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloID, Color: "Azul", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding azul: %v", err)
	}

	_, err = env.service.Update(ctx, usecase.UpdateInput{TargetID: rojo.ID, Color: strPtr("Azul"), RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestUpdateOwnColorDoesNotCollide(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env.pool, "Spray", "3.00", true)

	rojo, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloID, Color: "Rojo", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding rojo: %v", err)
	}

	updated, err := env.service.Update(ctx, usecase.UpdateInput{
		TargetID:    rojo.ID,
		Color:       strPtr("Rojo"),
		StockMinimo: int32Ptr(8),
		RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error updating own color: %v", err)
	}
	if updated.StockMinimo != 8 {
		t.Fatalf("expected stock_minimo=8, got %d", updated.StockMinimo)
	}
}

func TestDeleteWithZeroStockSucceeds(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env.pool, "Spray", "3.00", true)

	v, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloID, Color: "Rojo", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding variante: %v", err)
	}

	if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: v.ID, RequesterID: env.requesterID}); err != nil {
		t.Fatalf("unexpected error deleting variante with zero stock: %v", err)
	}

	_, err = env.service.Get(ctx, v.ID)
	if err == nil {
		t.Fatal("expected deleted variante to be not found")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteWithStockFails(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env.pool, "Spray", "3.00", true)

	v, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloID, Color: "Rojo", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding variante: %v", err)
	}
	setStock(t, env.pool, v.ID, "3")

	err = env.service.Delete(ctx, usecase.DeleteInput{TargetID: v.ID, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestListWithFilters(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloA := seedModeloEnvase(t, env.pool, "Spray", "3.00", true)
	modeloB := seedModeloEnvase(t, env.pool, "Roll-on", "1.00", true)

	rojo, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloA, Color: "Rojo", StockMinimo: 5, RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding rojo: %v", err)
	}
	azul, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, ModeloEnvaseID: modeloB, Color: "Azul", StockMinimo: 5, RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding azul: %v", err)
	}
	// Give azul enough stock that it's excluded from the stock_bajo filter;
	// rojo stays at zero (below its stock_minimo of 5).
	setStock(t, env.pool, azul.ID, "10")

	t.Run("filtra_por_modelo_envase_id", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, ModeloEnvaseID: modeloB})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ID != azul.ID {
			t.Fatalf("expected only azul, got total=%d items=%d", result.Total, len(result.Items))
		}
	})

	t.Run("stock_bajo_detecta_bajo_stock_minimo", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, StockBajo: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ID != rojo.ID {
			t.Fatalf("expected only rojo under stock_minimo, got total=%d items=%d", result.Total, len(result.Items))
		}
	})

	t.Run("include_deleted_expone_las_eliminadas", func(t *testing.T) {
		if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: rojo.ID, RequesterID: env.requesterID}); err != nil {
			t.Fatalf("deleting rojo: %v", err)
		}

		visible, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, Activo: "all"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if visible.Total != 1 {
			t.Fatalf("expected deleted rojo hidden by default, got total=%d", visible.Total)
		}

		withDeleted, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, Activo: "all", IncludeDeleted: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if withDeleted.Total != 2 {
			t.Fatalf("expected both variantes with include_deleted, got total=%d", withDeleted.Total)
		}
	})
}

func TestGetUnknownVarianteNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Get(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdateUnknownVarianteNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Update(context.Background(), usecase.UpdateInput{TargetID: 999999, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteUnknownVarianteNotFound(t *testing.T) {
	env := newTestEnv(t)

	err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: 999999, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}
