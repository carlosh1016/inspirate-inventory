package fragancias_test

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/fragancias"
)

func strPtr(s string) *string { return &s }

func TestCreateSuccessInitializesStock(t *testing.T) {
	env := newTestEnv(t)

	result, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "212 VIP", Genero: "femenina",
		GramosMinimo: "10.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.NombreComercial != "212 VIP" {
		t.Errorf("unexpected nombre_comercial: %q", result.NombreComercial)
	}
	if !result.StockVitrina.IsZero() || !result.StockBodega.IsZero() {
		t.Errorf("expected zero initial stock, got vitrina=%s bodega=%s", result.StockVitrina, result.StockBodega)
	}

	vitrina, bodega, err := env.service.StockActual.GetStockTotal(context.Background(), env.sedeID, stockactual.TipoItemFragancia, result.ID)
	if err != nil {
		t.Fatalf("unexpected error reading stock: %v", err)
	}
	if !vitrina.IsZero() || !bodega.IsZero() {
		t.Errorf("expected zero stock rows in DB, got vitrina=%s bodega=%s", vitrina, bodega)
	}
}

func TestCreateDuplicateNombreConflicts(t *testing.T) {
	env := newTestEnv(t)

	in := usecase.CreateInput{SedeID: env.sedeID, NombreComercial: "Chance", Genero: "femenina", GramosMinimo: "5.00", RequesterID: env.requesterID}
	if _, err := env.service.Create(context.Background(), in); err != nil {
		t.Fatalf("unexpected error on first create: %v", err)
	}

	_, err := env.service.Create(context.Background(), in)
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestCreateDuplicateNombreCaseInsensitive(t *testing.T) {
	env := newTestEnv(t)

	if _, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "Sauvage", Genero: "masculina", GramosMinimo: "5.00", RequesterID: env.requesterID,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "SAUVAGE", Genero: "masculina", GramosMinimo: "5.00", RequesterID: env.requesterID,
	})
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestUpdateNombreCollisionConflicts(t *testing.T) {
	env := newTestEnv(t)

	f1, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "Black Opium", Genero: "femenina", GramosMinimo: "5.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "Libre", Genero: "femenina", GramosMinimo: "5.00", RequesterID: env.requesterID,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: f1.ID, NombreComercial: strPtr("Libre"), RequesterID: env.requesterID,
	})
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestUpdateOwnNombreDoesNotCollide(t *testing.T) {
	env := newTestEnv(t)

	f, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "Good Girl", Genero: "femenina", GramosMinimo: "5.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: f.ID, NombreComercial: strPtr("Good Girl"), GramosMinimo: strPtr("7.50"), RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error updating with its own current nombre: %v", err)
	}
	if updated.GramosMinimo.StringFixed(2) != "7.50" {
		t.Errorf("expected gramos_minimo=7.50, got %s", updated.GramosMinimo)
	}
}

func TestDeleteWithZeroStockSucceeds(t *testing.T) {
	env := newTestEnv(t)

	f, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "J'adore", Genero: "femenina", GramosMinimo: "5.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: f.ID, RequesterID: env.requesterID}); err != nil {
		t.Fatalf("unexpected error deleting with zero stock: %v", err)
	}

	_, err = env.service.Get(context.Background(), f.ID)
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteWithStockFails(t *testing.T) {
	env := newTestEnv(t)

	f, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "Le Male", Genero: "masculina", GramosMinimo: "5.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	setStock(t, env.pool, f.ID, "3.50")

	err = env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: f.ID, RequesterID: env.requesterID})
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestRestoreSoftDeletedFragancia(t *testing.T) {
	env := newTestEnv(t)

	f, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "Alien", Genero: "femenina", GramosMinimo: "5.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: f.ID, RequesterID: env.requesterID}); err != nil {
		t.Fatalf("unexpected error deleting: %v", err)
	}

	restored, err := env.service.Restore(context.Background(), usecase.RestoreInput{TargetID: f.ID, RequesterID: env.requesterID})
	if err != nil {
		t.Fatalf("unexpected error restoring: %v", err)
	}
	if !restored.Activo {
		t.Error("expected restored fragancia to be active")
	}

	// Restoring something that is NOT deleted (or doesn't exist) must fail.
	_, err = env.service.Restore(context.Background(), usecase.RestoreInput{TargetID: f.ID, RequesterID: env.requesterID})
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestListWithFilters(t *testing.T) {
	env := newTestEnv(t)

	masc, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "Bleu de Chanel", Genero: "masculina", GramosMinimo: "10.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fem, err := env.service.Create(context.Background(), usecase.CreateInput{
		SedeID: env.sedeID, NombreComercial: "Miss Dior", Genero: "femenina", GramosMinimo: "10.00", RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Give it stock at/above its own gramos_minimo so it never shows up as
	// "low stock" — only masc (kept at zero, then explicitly lowered) should.
	setStock(t, env.pool, fem.ID, "10.00")

	t.Run("filtra por genero", func(t *testing.T) {
		result, err := env.service.List(context.Background(), usecase.ListInput{SedeID: env.sedeID, Genero: "masculina"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 {
			t.Fatalf("expected 1 masculina, got %d", result.Total)
		}
	})

	t.Run("filtra por q", func(t *testing.T) {
		result, err := env.service.List(context.Background(), usecase.ListInput{SedeID: env.sedeID, Q: "dior"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 {
			t.Fatalf("expected 1 match for 'dior', got %d", result.Total)
		}
	})

	t.Run("stock_bajo detecta bajo gramos_minimo", func(t *testing.T) {
		setStock(t, env.pool, masc.ID, "2.00")

		result, err := env.service.List(context.Background(), usecase.ListInput{SedeID: env.sedeID, StockBajo: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 {
			t.Fatalf("expected 1 fragancia under gramos_minimo, got %d", result.Total)
		}
		if result.Items[0].ID != masc.ID {
			t.Fatalf("expected the low-stock fragancia to be %d, got %d", masc.ID, result.Items[0].ID)
		}
	})

	t.Run("include_deleted expone las eliminadas", func(t *testing.T) {
		// The previous subtest left masc with stock=2; clear it so delete succeeds.
		setStock(t, env.pool, masc.ID, "0")
		if err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: masc.ID, RequesterID: env.requesterID}); err != nil {
			t.Fatalf("unexpected error deleting: %v", err)
		}

		// By default (IncludeDeleted=false), a soft-deleted fragancia is
		// invisible — this is a distinct dimension from the activo flag.
		activeOnly, err := env.service.List(context.Background(), usecase.ListInput{SedeID: env.sedeID, Activo: "all"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if activeOnly.Total != 1 {
			t.Fatalf("expected 1 non-deleted fragancia after deleting one, got %d", activeOnly.Total)
		}

		result, err := env.service.List(context.Background(), usecase.ListInput{SedeID: env.sedeID, Activo: "all", IncludeDeleted: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 2 {
			t.Fatalf("expected 2 fragancias with include_deleted=true, got %d", result.Total)
		}
	})
}

func TestGetUnknownFraganciaNotFound(t *testing.T) {
	env := newTestEnv(t)
	_, err := env.service.Get(context.Background(), 9_999_999)
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdateUnknownFraganciaNotFound(t *testing.T) {
	env := newTestEnv(t)
	_, err := env.service.Update(context.Background(), usecase.UpdateInput{
		TargetID: 9_999_999, NombreComercial: strPtr("Cualquiera"), RequesterID: env.requesterID,
	})
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteUnknownFraganciaNotFound(t *testing.T) {
	env := newTestEnv(t)
	err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: 9_999_999, RequesterID: env.requesterID})
	assertCode(t, err, domainerrors.CodeNotFound)
}

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
