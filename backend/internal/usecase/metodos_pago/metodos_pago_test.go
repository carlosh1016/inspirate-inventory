package metodospago_test

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/metodos_pago"
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

func TestCreateSuccess(t *testing.T) {
	env := newTestEnv(t)

	m, err := env.service.Create(context.Background(), usecase.CreateInput{Nombre: "Efectivo", Codigo: "EFEC", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Nombre != "Efectivo" || m.Codigo != "EFEC" {
		t.Fatalf("expected Efectivo/EFEC, got %q/%q", m.Nombre, m.Codigo)
	}
	if !m.Activo {
		t.Fatalf("expected new metodo_pago to be activo")
	}
}

func TestCreateDuplicateCodigoConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Efectivo", Codigo: "EFEC", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding first metodo_pago: %v", err)
	}

	_, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Efectivo Caja 2", Codigo: "efec", RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestCreateDuplicateNombreConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Efectivo", Codigo: "EFEC", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding first metodo_pago: %v", err)
	}

	_, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "EFECTIVO", Codigo: "EFEC2", RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestUpdateOwnCodigoDoesNotCollide(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	m, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Efectivo", Codigo: "EFEC", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding metodo_pago: %v", err)
	}

	updated, err := env.service.Update(ctx, usecase.UpdateInput{TargetID: m.ID, Codigo: strPtr("EFEC"), Nombre: strPtr("Efectivo Mostrador"), RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("unexpected error updating own codigo: %v", err)
	}
	if updated.Nombre != "Efectivo Mostrador" {
		t.Fatalf("expected nombre updated, got %q", updated.Nombre)
	}
}

func TestUpdateCodigoCollisionConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	efectivo, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Efectivo", Codigo: "EFEC", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding efectivo: %v", err)
	}
	if _, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Nequi", Codigo: "NEQ", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding nequi: %v", err)
	}

	_, err = env.service.Update(ctx, usecase.UpdateInput{TargetID: efectivo.ID, Codigo: strPtr("NEQ"), RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestDeleteWithoutVentasHardDeletes(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	m, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Efectivo", Codigo: "EFEC", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding metodo_pago: %v", err)
	}

	if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: m.ID, RequesterID: env.requesterID}); err != nil {
		t.Fatalf("unexpected error deleting metodo_pago without ventas: %v", err)
	}

	var count int
	if err := env.pool.QueryRow(ctx, `SELECT COUNT(*) FROM metodos_pago WHERE id = $1`, m.ID).Scan(&count); err != nil {
		t.Fatalf("checking row existence: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected metodo_pago row to be physically deleted, found %d rows", count)
	}
}

func TestDeleteWithVentasSoftDeletes(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	m, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Efectivo", Codigo: "EFEC", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding metodo_pago: %v", err)
	}
	seedVenta(t, env.pool, env.sedeID, env.requesterID, m.ID)

	if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: m.ID, RequesterID: env.requesterID}); err != nil {
		t.Fatalf("unexpected error deleting metodo_pago with ventas: %v", err)
	}

	var deletedAtIsNull bool
	if err := env.pool.QueryRow(ctx, `SELECT deleted_at IS NULL FROM metodos_pago WHERE id = $1`, m.ID).Scan(&deletedAtIsNull); err != nil {
		t.Fatalf("expected the row to still exist (soft-deleted), got: %v", err)
	}
	if deletedAtIsNull {
		t.Fatal("expected deleted_at to be set (soft delete), got NULL")
	}

	_, err = env.service.Get(ctx, m.ID)
	if err == nil {
		t.Fatal("expected soft-deleted metodo_pago to be not found via Get")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestListWithFilters(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Efectivo", Codigo: "EFEC", RequesterID: env.requesterID}); err != nil {
		t.Fatalf("seeding efectivo: %v", err)
	}
	nequi, err := env.service.Create(ctx, usecase.CreateInput{Nombre: "Nequi", Codigo: "NEQ", RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("seeding nequi: %v", err)
	}

	t.Run("filtra_por_q", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{Q: "neq"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ID != nequi.ID {
			t.Fatalf("expected only nequi, got total=%d items=%d", result.Total, len(result.Items))
		}
	})

	t.Run("include_deleted_expone_las_eliminadas", func(t *testing.T) {
		seedVenta(t, env.pool, env.sedeID, env.requesterID, nequi.ID)
		if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: nequi.ID, RequesterID: env.requesterID}); err != nil {
			t.Fatalf("deleting nequi: %v", err)
		}

		visible, err := env.service.List(ctx, usecase.ListInput{Activo: "all"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if visible.Total != 1 {
			t.Fatalf("expected soft-deleted nequi hidden by default, got total=%d", visible.Total)
		}

		withDeleted, err := env.service.List(ctx, usecase.ListInput{Activo: "all", IncludeDeleted: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if withDeleted.Total != 2 {
			t.Fatalf("expected both metodos_pago with include_deleted, got total=%d", withDeleted.Total)
		}
	})
}

func TestGetUnknownMetodoPagoNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Get(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdateUnknownMetodoPagoNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Update(context.Background(), usecase.UpdateInput{TargetID: 999999, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteUnknownMetodoPagoNotFound(t *testing.T) {
	env := newTestEnv(t)

	err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: 999999, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}
