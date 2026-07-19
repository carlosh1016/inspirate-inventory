package productos_test

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/productos"
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

	p, err := env.service.Create(ctx, usecase.CreateInput{
		SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Nombre != "Vela de Vainilla" {
		t.Fatalf("expected nombre Vela de Vainilla, got %q", p.Nombre)
	}
	if p.StockVitrina.Sign() != 0 || p.StockBodega.Sign() != 0 {
		t.Fatalf("expected zero stock on creation, got vitrina=%s bodega=%s", p.StockVitrina, p.StockBodega)
	}
}

func TestCreateDuplicateNombreCategoriaConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	in := usecase.CreateInput{SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID}
	if _, err := env.service.Create(ctx, in); err != nil {
		t.Fatalf("seeding first producto: %v", err)
	}

	_, err := env.service.Create(ctx, in)
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestCreateDuplicateNombreCategoriaCaseInsensitive(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding first producto: %v", err)
	}

	_, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "VELA DE VAINILLA", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestCreateSameNombreDifferentCategoriaAllowed(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Kit Regalo", Categoria: "regalo", Precio: "40000.00", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding first producto: %v", err)
	}

	_, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Kit Regalo", Categoria: "hogar", Precio: "40000.00", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("expected same nombre under a different categoria to be allowed, got: %v", err)
	}
}

func TestUpdateAsAdminCanEditAllFields(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding producto: %v", err)
	}

	updated, err := env.service.Update(ctx, usecase.UpdateInput{
		TargetID: p.ID, Nombre: strPtr("Vela de Lavanda"), Precio: strPtr("27000.00"), RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error updating as admin: %v", err)
	}
	if updated.Nombre != "Vela de Lavanda" || updated.Precio.StringFixed(2) != "27000.00" {
		t.Fatalf("expected updated nombre/precio, got nombre=%q precio=%q", updated.Nombre, updated.Precio.StringFixed(2))
	}
}

func TestUpdateAsVendedoraRestrictedToStockMinimo(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding producto: %v", err)
	}

	_, err = env.service.Update(ctx, usecase.UpdateInput{
		TargetID: p.ID, Nombre: strPtr("Vela de Lavanda"), RequesterID: env.requesterID, RequesterRole: usecase.RolVendedora,
	})
	if err == nil {
		t.Fatal("expected a forbidden error, got nil")
	}
	assertCode(t, err, domainerrors.CodeForbidden)

	updated, err := env.service.Update(ctx, usecase.UpdateInput{
		TargetID: p.ID, StockMinimo: int32Ptr(3), RequesterID: env.requesterID, RequesterRole: usecase.RolVendedora,
	})
	if err != nil {
		t.Fatalf("expected vendedora to edit stock_minimo, got: %v", err)
	}
	if updated.StockMinimo != 3 {
		t.Fatalf("expected stock_minimo=3, got %d", updated.StockMinimo)
	}
}

func TestUpdateNombreCategoriaCollisionConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	vela, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding vela: %v", err)
	}
	if _, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Difusor", Categoria: "hogar", Precio: "30000.00", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding difusor: %v", err)
	}

	_, err = env.service.Update(ctx, usecase.UpdateInput{TargetID: vela.ID, Nombre: strPtr("Difusor"), RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestDeleteWithZeroStockSucceeds(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding producto: %v", err)
	}

	if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: p.ID, RequesterID: env.requesterID}); err != nil {
		t.Fatalf("unexpected error deleting producto with zero stock: %v", err)
	}

	_, err = env.service.Get(ctx, p.ID)
	if err == nil {
		t.Fatal("expected deleted producto to be not found")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteWithStockFails(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding producto: %v", err)
	}
	setStock(t, env.pool, p.ID, "5")

	err = env.service.Delete(ctx, usecase.DeleteInput{TargetID: p.ID, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestListWithFilters(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	vela, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Vela de Vainilla", Categoria: "hogar", Precio: "25000.00", StockMinimo: 5, RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding vela: %v", err)
	}
	crema, err := env.service.Create(ctx, usecase.CreateInput{SedeID: env.sedeID, Nombre: "Crema Corporal", Categoria: "crema", Precio: "18000.00", StockMinimo: 5, RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding crema: %v", err)
	}
	// Give crema enough stock that it's excluded from stock_bajo; vela stays
	// at zero (below its stock_minimo of 5).
	setStock(t, env.pool, crema.ID, "10")

	t.Run("filtra_por_categoria", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, Categoria: "crema"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ID != crema.ID {
			t.Fatalf("expected only crema, got total=%d items=%d", result.Total, len(result.Items))
		}
	})

	t.Run("stock_bajo_detecta_bajo_stock_minimo", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, StockBajo: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ID != vela.ID {
			t.Fatalf("expected only vela under stock_minimo, got total=%d items=%d", result.Total, len(result.Items))
		}
	})

	t.Run("include_deleted_expone_las_eliminadas", func(t *testing.T) {
		if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: vela.ID, RequesterID: env.requesterID}); err != nil {
			t.Fatalf("deleting vela: %v", err)
		}

		visible, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, Activo: "all"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if visible.Total != 1 {
			t.Fatalf("expected deleted vela hidden by default, got total=%d", visible.Total)
		}

		withDeleted, err := env.service.List(ctx, usecase.ListInput{SedeID: env.sedeID, Activo: "all", IncludeDeleted: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if withDeleted.Total != 2 {
			t.Fatalf("expected both productos with include_deleted, got total=%d", withDeleted.Total)
		}
	})
}

func TestGetUnknownProductoNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Get(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdateUnknownProductoNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Update(context.Background(), usecase.UpdateInput{TargetID: 999999, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteUnknownProductoNotFound(t *testing.T) {
	env := newTestEnv(t)

	err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: 999999, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}
