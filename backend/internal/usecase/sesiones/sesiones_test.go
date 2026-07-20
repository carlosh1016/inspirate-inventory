package sesiones_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecasesesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
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

// seedSesionCerrada inserts a closed sesion directly, bypassing Entrada/
// Salida (which stamp time.Now()) so tests can control exact timestamps.
func seedSesionCerrada(t *testing.T, env *testEnv, usuarioID int64, entrada, salida time.Time) {
	t.Helper()
	_, err := env.pool.Exec(context.Background(),
		`INSERT INTO sesiones_laborales (sede_id, usuario_id, entrada_at, salida_at, horas_trabajadas)
		 VALUES ($1, $2, $3::timestamptz, $4::timestamptz, $4::timestamptz - $3::timestamptz)`,
		env.sedeID, usuarioID, entrada, salida,
	)
	if err != nil {
		t.Fatalf("seeding sesion cerrada: %v", err)
	}
}

// --- 1/2. Entrada ---

func TestEntradaCreaSesionAbierta(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	sesion, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sesion.SalidaAt != nil || sesion.HorasTrabajadas != nil {
		t.Fatalf("expected an open sesion, got %+v", sesion)
	}
}

func TestEntradaCuandoYaHayAbiertaFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID}); err != nil {
		t.Fatalf("unexpected error on first entrada: %v", err)
	}

	_, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

// --- 3/4. Salida ---

func TestSalidaCierraSesionYCalculaHoras(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID}); err != nil {
		t.Fatalf("unexpected error on entrada: %v", err)
	}

	sesion, err := env.service.Salida(ctx, env.vendedoraID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sesion.SalidaAt == nil {
		t.Fatal("expected salida_at populated")
	}
	if sesion.HorasTrabajadas == nil {
		t.Fatal("expected horas_trabajadas populated")
	}
}

func TestSalidaSinSesionAbiertaFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_, err := env.service.Salida(ctx, env.vendedoraID)
	if err == nil {
		t.Fatal("expected a not_found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

// --- 5. Sesiones simultáneas ---

func TestDosVendedorasConSesionesAbiertasSimultaneamente(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID}); err != nil {
		t.Fatalf("unexpected error for vendedora 1: %v", err)
	}
	if _, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedora2}); err != nil {
		t.Fatalf("unexpected error for vendedora 2: %v", err)
	}

	result, err := env.service.List(ctx, usecasesesiones.ListInput{SedeID: env.sedeID, Abiertas: true})
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
	if result.Total != 2 {
		t.Fatalf("expected 2 open sesiones, got %d", result.Total)
	}
}

// --- 7/8/9. List ---

func TestListFiltraPorUsuario(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID}); err != nil {
		t.Fatalf("unexpected error for vendedora 1: %v", err)
	}
	if _, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedora2}); err != nil {
		t.Fatalf("unexpected error for vendedora 2: %v", err)
	}

	scoped, err := env.service.List(ctx, usecasesesiones.ListInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scoped.Total != 1 || scoped.Sesiones[0].UsuarioID != env.vendedoraID {
		t.Fatalf("expected exactly vendedora 1's sesion, got %+v", scoped.Sesiones)
	}

	all, err := env.service.List(ctx, usecasesesiones.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if all.Total != 2 {
		t.Fatalf("expected 2 sesiones without a usuario filter, got %d", all.Total)
	}
}

func TestListFiltraAbiertas(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID}); err != nil {
		t.Fatalf("unexpected error on entrada: %v", err)
	}
	seedSesionCerrada(t, env, env.vendedora2, time.Now().Add(-8*time.Hour), time.Now().Add(-1*time.Hour))

	abiertas, err := env.service.List(ctx, usecasesesiones.ListInput{SedeID: env.sedeID, Abiertas: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if abiertas.Total != 1 || abiertas.Sesiones[0].UsuarioID != env.vendedoraID {
		t.Fatalf("expected only the open sesion, got %+v", abiertas.Sesiones)
	}

	todas, err := env.service.List(ctx, usecasesesiones.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todas.Total != 2 {
		t.Fatalf("expected 2 sesiones without the abiertas filter, got %d", todas.Total)
	}
}

// --- 10/11. Update ---

func TestUpdateRecalculaHoras(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	entrada := time.Now().Add(-10 * time.Hour)
	salida := time.Now().Add(-2 * time.Hour)
	seedSesionCerrada(t, env, env.vendedoraID, entrada, salida)

	sesiones, err := env.service.List(ctx, usecasesesiones.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
	if len(sesiones.Sesiones) != 1 {
		t.Fatalf("expected 1 sesion seeded, got %d", len(sesiones.Sesiones))
	}
	sesionID := sesiones.Sesiones[0].ID

	nuevaEntrada := entrada.Add(-1 * time.Hour) // one more hour worked
	updated, err := env.service.Update(ctx, usecasesesiones.UpdateInput{
		TargetID: sesionID, EntradaAt: &nuevaEntrada, RequesterID: env.adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.HorasTrabajadas == nil || updated.HorasTrabajadas.Hours() != 9 {
		t.Fatalf("expected horas_trabajadas recalculated to 9h, got %v", updated.HorasTrabajadas)
	}
}

func TestUpdateEntradaPosteriorASalidaFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	entrada := time.Now().Add(-10 * time.Hour)
	salida := time.Now().Add(-2 * time.Hour)
	seedSesionCerrada(t, env, env.vendedoraID, entrada, salida)

	sesiones, err := env.service.List(ctx, usecasesesiones.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
	sesionID := sesiones.Sesiones[0].ID

	entradaInvalida := salida.Add(1 * time.Hour)
	_, err = env.service.Update(ctx, usecasesesiones.UpdateInput{
		TargetID: sesionID, EntradaAt: &entradaInvalida, RequesterID: env.adminID,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

// --- 12/13. Resumen ---

func TestResumenTotalesPorVendedora(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	base := time.Now().Add(-48 * time.Hour)
	seedSesionCerrada(t, env, env.vendedoraID, base, base.Add(8*time.Hour))
	seedSesionCerrada(t, env, env.vendedoraID, base.Add(24*time.Hour), base.Add(24*time.Hour+6*time.Hour))
	seedSesionCerrada(t, env, env.vendedora2, base, base.Add(4*time.Hour))

	resumen, err := env.service.Resumen(ctx, usecasesesiones.ResumenInput{
		FechaDesde: base.Add(-time.Hour), FechaHasta: time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resumen) != 2 {
		t.Fatalf("expected 2 vendedoras in resumen, got %d", len(resumen))
	}

	var v1, v2 *usecasesesiones.ResumenItem
	for i := range resumen {
		switch resumen[i].UsuarioID {
		case env.vendedoraID:
			v1 = &resumen[i]
		case env.vendedora2:
			v2 = &resumen[i]
		}
	}
	if v1 == nil || v1.SesionesCount != 2 || v1.DiasTrabajados != 2 || v1.TotalHoras != 14*time.Hour {
		t.Fatalf("expected vendedora 1: 2 sesiones, 2 dias, 14h total, got %+v", v1)
	}
	if v2 == nil || v2.SesionesCount != 1 || v2.TotalHoras != 4*time.Hour {
		t.Fatalf("expected vendedora 2: 1 sesion, 4h total, got %+v", v2)
	}
}

func TestResumenIgnoraSesionesAbiertas(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Entrada(ctx, usecasesesiones.EntradaInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID}); err != nil {
		t.Fatalf("unexpected error on entrada: %v", err)
	}

	resumen, err := env.service.Resumen(ctx, usecasesesiones.ResumenInput{
		FechaDesde: time.Now().Add(-time.Hour), FechaHasta: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resumen) != 0 {
		t.Fatalf("expected open sesiones excluded from resumen, got %+v", resumen)
	}
}
