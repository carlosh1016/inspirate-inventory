package modelosenvase_test

import (
	"context"
	"errors"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/modelos_envase"
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

func createModelo(t *testing.T, env *testEnv, tipo, tamanoOz string) usecase.CreateInput {
	t.Helper()
	return usecase.CreateInput{
		Tipo:               tipo,
		TamanoOz:           tamanoOz,
		EquivGramos:        "50.00",
		PrecioSolo:         "10000.00",
		PrecioConFragancia: "25000.00",
		PrecioRecarga:      "15000.00",
		RequesterID:        env.requesterID,
		IP:                 "127.0.0.1",
		UserAgent:          "test-agent",
	}
}

func TestCreateSuccess(t *testing.T) {
	env := newTestEnv(t)

	m, err := env.service.Create(context.Background(), createModelo(t, env, "Spray", "3.00"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Tipo != "Spray" {
		t.Fatalf("expected tipo Spray, got %q", m.Tipo)
	}
	if m.VariantesActivas != 0 {
		t.Fatalf("expected 0 variantes_activas, got %d", m.VariantesActivas)
	}
	if !m.Activo {
		t.Fatalf("expected new modelo to be activo")
	}
}

func TestCreateDuplicateTipoTamanoConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Create(ctx, createModelo(t, env, "Spray", "3.00")); err != nil {
		t.Fatalf("seeding first modelo: %v", err)
	}

	_, err := env.service.Create(ctx, createModelo(t, env, "Spray", "3.00"))
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestCreateDuplicateTipoTamanoCaseInsensitive(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Create(ctx, createModelo(t, env, "Spray", "3.00")); err != nil {
		t.Fatalf("seeding first modelo: %v", err)
	}

	_, err := env.service.Create(ctx, createModelo(t, env, "SPRAY", "3.00"))
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestUpdateTipoTamanoCollisionConflicts(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	spray, err := env.service.Create(ctx, createModelo(t, env, "Spray", "3.00"))
	if err != nil {
		t.Fatalf("seeding spray: %v", err)
	}
	if _, err := env.service.Create(ctx, createModelo(t, env, "Roll-on", "3.00")); err != nil {
		t.Fatalf("seeding roll-on: %v", err)
	}

	_, err = env.service.Update(ctx, usecase.UpdateInput{
		TargetID:    spray.ID,
		Tipo:        strPtr("Roll-on"),
		RequesterID: env.requesterID,
	})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestUpdateOwnTipoTamanoDoesNotCollide(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	spray, err := env.service.Create(ctx, createModelo(t, env, "Spray", "3.00"))
	if err != nil {
		t.Fatalf("seeding spray: %v", err)
	}

	updated, err := env.service.Update(ctx, usecase.UpdateInput{
		TargetID:    spray.ID,
		Tipo:        strPtr("Spray"),
		PrecioSolo:  strPtr("11000.00"),
		RequesterID: env.requesterID,
	})
	if err != nil {
		t.Fatalf("unexpected error updating own tipo: %v", err)
	}
	if updated.PrecioSolo.StringFixed(2) != "11000.00" {
		t.Fatalf("expected precio_solo 11000.00, got %q", updated.PrecioSolo.StringFixed(2))
	}
}

func TestDeleteWithoutVariantesSucceeds(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	m, err := env.service.Create(ctx, createModelo(t, env, "Spray", "3.00"))
	if err != nil {
		t.Fatalf("seeding modelo: %v", err)
	}

	if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: m.ID, RequesterID: env.requesterID}); err != nil {
		t.Fatalf("unexpected error deleting modelo without variantes: %v", err)
	}

	_, err = env.service.Get(ctx, m.ID)
	if err == nil {
		t.Fatal("expected deleted modelo to be not found")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeleteWithVariantesActivasFails(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	m, err := env.service.Create(ctx, createModelo(t, env, "Spray", "3.00"))
	if err != nil {
		t.Fatalf("seeding modelo: %v", err)
	}

	seedVarianteEnvase(t, env.pool, m.ID, env.sedeID, "Rojo")

	err = env.service.Delete(ctx, usecase.DeleteInput{TargetID: m.ID, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestListWithFilters(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Create(ctx, createModelo(t, env, "Spray", "3.00")); err != nil {
		t.Fatalf("seeding spray: %v", err)
	}
	rollon, err := env.service.Create(ctx, createModelo(t, env, "Roll-on", "1.00"))
	if err != nil {
		t.Fatalf("seeding roll-on: %v", err)
	}

	t.Run("filtra_por_tipo", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{Tipo: "Roll-on"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ID != rollon.ID {
			t.Fatalf("expected only roll-on, got total=%d items=%d", result.Total, len(result.Items))
		}
	})

	t.Run("filtra_por_q", func(t *testing.T) {
		result, err := env.service.List(ctx, usecase.ListInput{Q: "spr"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 1 {
			t.Fatalf("expected 1 match for q=spr, got %d", result.Total)
		}
	})

	t.Run("include_deleted_expone_las_eliminadas", func(t *testing.T) {
		if err := env.service.Delete(ctx, usecase.DeleteInput{TargetID: rollon.ID, RequesterID: env.requesterID}); err != nil {
			t.Fatalf("deleting roll-on: %v", err)
		}

		visible, err := env.service.List(ctx, usecase.ListInput{Activo: "all"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if visible.Total != 1 {
			t.Fatalf("expected deleted roll-on hidden by default, got total=%d", visible.Total)
		}

		withDeleted, err := env.service.List(ctx, usecase.ListInput{Activo: "all", IncludeDeleted: true})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if withDeleted.Total != 2 {
			t.Fatalf("expected both modelos with include_deleted, got total=%d", withDeleted.Total)
		}
	})
}

func TestGetUnknownModeloNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Get(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdateUnknownModeloNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Update(context.Background(), usecase.UpdateInput{TargetID: 999999, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestCreateSinVariantesAutoCreaVarianteOculta(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	in := createModelo(t, env, "Envase de lujo", "1.00")
	in.SinVariantes = true
	in.SedeID = env.sedeID

	m, err := env.service.Create(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.TieneVariantes {
		t.Fatal("expected TieneVariantes=false")
	}
	if m.VariantesActivas != 1 {
		t.Fatalf("expected exactly 1 auto-created variante, got %d", m.VariantesActivas)
	}

	var color string
	var vitrina, bodega int
	err = env.pool.QueryRow(ctx,
		`SELECT v.color,
		   (SELECT COUNT(*) FROM stock_actual WHERE tipo_item = 'variante_envase' AND item_id = v.id AND ubicacion = 'vitrina'),
		   (SELECT COUNT(*) FROM stock_actual WHERE tipo_item = 'variante_envase' AND item_id = v.id AND ubicacion = 'bodega')
		 FROM variantes_envase v WHERE v.modelo_envase_id = $1`,
		m.ID,
	).Scan(&color, &vitrina, &bodega)
	if err != nil {
		t.Fatalf("querying auto-created variante: %v", err)
	}
	if color != "Único" {
		t.Fatalf("expected sentinel color Único, got %q", color)
	}
	if vitrina != 1 || bodega != 1 {
		t.Fatalf("expected stock rows seeded for both ubicaciones, got vitrina=%d bodega=%d", vitrina, bodega)
	}
}

func TestDeleteUnknownModeloNotFound(t *testing.T) {
	env := newTestEnv(t)

	err := env.service.Delete(context.Background(), usecase.DeleteInput{TargetID: 999999, RequesterID: env.requesterID})
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}
