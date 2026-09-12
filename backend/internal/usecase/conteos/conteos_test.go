package conteos_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecaseconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/conteos"
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

func mustDecimal(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatalf("parsing decimal %q: %v", s, err)
	}
	return d
}

func TestCrearConteoGeneraFilasDesdeStockActual(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	f1 := seedFragancia(t, env, "Chanel No. 5", "100.00")
	f2 := seedFragancia(t, env, "Dior Sauvage", "50.50")

	conteo, err := env.service.Crear(ctx, usecaseconteos.CrearInput{
		SedeID: env.sedeID, RequesterID: env.adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conteo.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(conteo.Items))
	}

	byFragancia := map[int64]decimal.Decimal{}
	for _, it := range conteo.Items {
		byFragancia[it.FraganciaID] = it.GramosSistema
		// No previous conteo exists yet — saldo_inicial falls back to the
		// system's current stock, same as gramos_sistema.
		if !it.SaldoInicial.Equal(it.GramosSistema) {
			t.Fatalf("expected saldo_inicial == gramos_sistema on first conteo, got %s vs %s", it.SaldoInicial, it.GramosSistema)
		}
		if it.GramosFisico != nil {
			t.Fatalf("expected gramos_fisico nil before counting, got %v", it.GramosFisico)
		}
	}
	if !byFragancia[f1].Equal(mustDecimal(t, "100.00")) {
		t.Fatalf("expected f1 gramos_sistema 100.00, got %s", byFragancia[f1])
	}
	if !byFragancia[f2].Equal(mustDecimal(t, "50.50")) {
		t.Fatalf("expected f2 gramos_sistema 50.50, got %s", byFragancia[f2])
	}
}

func TestCrearConteoSumaEntradasDelMesAlSaldoInicialPrimerConteo(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	// Fragancia nueva: entra mercancía (100g) y se vende parte (40g) dentro
	// del mismo mes, sin conteo anterior. El stock actual (60g) por sí solo
	// escondería la venta; saldo_inicial debe reconstruir los 100g de base.
	f1 := seedFragancia(t, env, "Chanel No. 5", "60.00")
	seedMovimiento(t, env, f1, "entrada_mercancia", "100.00")
	seedMovimiento(t, env, f1, "venta", "-40.00")

	conteo, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	item := conteo.Items[0]
	if !item.SaldoInicial.Equal(mustDecimal(t, "100.00")) {
		t.Fatalf("expected saldo_inicial 100.00 (60 stock + 40 reconstruida por la venta), got %s", item.SaldoInicial)
	}
	if !item.GramosSistema.Equal(mustDecimal(t, "60.00")) {
		t.Fatalf("expected gramos_sistema 60.00, got %s", item.GramosSistema)
	}
	pct := item.PorcentajeVendido()
	if pct == nil || !pct.Equal(mustDecimal(t, "40")) {
		t.Fatalf("expected 40%% vendido, got %v", pct)
	}
}

func TestCrearConteoSumaEntradasDelMesSobreSaldoInicialPrevio(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	f1 := seedFragancia(t, env, "Chanel No. 5", "150.00")

	anterior := mesAnterior(env)
	conteoAnterior, err := env.service.Crear(ctx, usecaseconteos.CrearInput{
		SedeID: env.sedeID, RequesterID: env.adminID, Periodo: &anterior,
	})
	if err != nil {
		t.Fatalf("unexpected error creating conteo anterior: %v", err)
	}
	if _, err := env.service.ActualizarGramosFisico(ctx, usecaseconteos.ActualizarGramosFisicoInput{
		ConteoID: conteoAnterior.ID, ItemID: conteoAnterior.Items[0].ID, GramosFisico: mustDecimal(t, "200.00"), RequesterID: env.adminID,
	}); err != nil {
		t.Fatalf("unexpected error updating gramos_fisico: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecaseconteos.CerrarInput{TargetID: conteoAnterior.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error cerrando conteo anterior: %v", err)
	}

	// Este mes: se compran 100g más (entrada_mercancia) — deben sumarse al
	// saldo_inicial heredado (200), no perderse en el stock actual.
	seedMovimiento(t, env, f1, "entrada_mercancia", "100.00")

	conteoActual, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error creating conteo actual: %v", err)
	}
	item := conteoActual.Items[0]
	if !item.SaldoInicial.Equal(mustDecimal(t, "300.00")) {
		t.Fatalf("expected saldo_inicial 300.00 (200 heredado + 100 comprados este mes), got %s", item.SaldoInicial)
	}
	if !item.GramosSistema.Equal(mustDecimal(t, "150.00")) {
		t.Fatalf("expected gramos_sistema 150.00 (stock actual sin ajustar), got %s", item.GramosSistema)
	}
}

func TestCrearConteoConflictoSiYaExisteParaElPeriodo(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	seedFragancia(t, env, "Chanel No. 5", "100.00")

	if _, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error on first Crear: %v", err)
	}

	_, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID})
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestCrearConteoHeredaSaldoInicialDelConteoCerradoAnterior(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	f1 := seedFragancia(t, env, "Chanel No. 5", "100.00")

	anterior := mesAnterior(env)
	conteoAnterior, err := env.service.Crear(ctx, usecaseconteos.CrearInput{
		SedeID: env.sedeID, RequesterID: env.adminID, Periodo: &anterior,
	})
	if err != nil {
		t.Fatalf("unexpected error creating conteo anterior: %v", err)
	}
	itemAnterior := conteoAnterior.Items[0]

	// Cuenta física del mes anterior: 98g (2g menos que el sistema).
	if _, err := env.service.ActualizarGramosFisico(ctx, usecaseconteos.ActualizarGramosFisicoInput{
		ConteoID: conteoAnterior.ID, ItemID: itemAnterior.ID, GramosFisico: mustDecimal(t, "98.00"), RequesterID: env.adminID,
	}); err != nil {
		t.Fatalf("unexpected error updating gramos_fisico: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecaseconteos.CerrarInput{TargetID: conteoAnterior.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error cerrando conteo anterior: %v", err)
	}

	// El stock del sistema cambió (se vendieron 10g) desde el mes anterior.
	if _, err := env.pool.Exec(ctx,
		`UPDATE stock_actual SET cantidad = '88.00' WHERE tipo_item = 'fragancia' AND item_id = $1 AND ubicacion = 'bodega'`,
		f1,
	); err != nil {
		t.Fatalf("updating stock: %v", err)
	}

	conteoActual, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error creating conteo actual: %v", err)
	}
	item := conteoActual.Items[0]
	if !item.SaldoInicial.Equal(mustDecimal(t, "98.00")) {
		t.Fatalf("expected saldo_inicial 98.00 (inherited from last month's gramos_fisico), got %s", item.SaldoInicial)
	}
	if !item.GramosSistema.Equal(mustDecimal(t, "88.00")) {
		t.Fatalf("expected gramos_sistema 88.00 (current stock), got %s", item.GramosSistema)
	}
}

func TestActualizarGramosFisicoCalculaDiferenciaYLimite(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	seedFragancia(t, env, "Chanel No. 5", "100.00")

	conteo, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	item := conteo.Items[0]

	// Dentro del límite: diferencia de 1.5g.
	updated, err := env.service.ActualizarGramosFisico(ctx, usecaseconteos.ActualizarGramosFisicoInput{
		ConteoID: conteo.ID, ItemID: item.ID, GramosFisico: mustDecimal(t, "98.50"), RequesterID: env.adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.ExcedeLimite() {
		t.Fatal("expected 1.5g difference to be within the 2g limit")
	}
	if d := updated.Diferencia(); d == nil || !d.Equal(mustDecimal(t, "-1.50")) {
		t.Fatalf("expected diferencia -1.50, got %v", d)
	}

	// Excede el límite: diferencia de 5g.
	updated, err = env.service.ActualizarGramosFisico(ctx, usecaseconteos.ActualizarGramosFisicoInput{
		ConteoID: conteo.ID, ItemID: item.ID, GramosFisico: mustDecimal(t, "95.00"), RequesterID: env.adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated.ExcedeLimite() {
		t.Fatal("expected 5g difference to exceed the 2g limit")
	}
}

func TestActualizarGramosFisicoRechazaNegativo(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	seedFragancia(t, env, "Chanel No. 5", "100.00")

	conteo, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = env.service.ActualizarGramosFisico(ctx, usecaseconteos.ActualizarGramosFisicoInput{
		ConteoID: conteo.ID, ItemID: conteo.Items[0].ID, GramosFisico: mustDecimal(t, "-1"), RequesterID: env.adminID,
	})
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestCerrarBloqueaEdicionesPosteriores(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	seedFragancia(t, env, "Chanel No. 5", "100.00")

	conteo, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecaseconteos.CerrarInput{TargetID: conteo.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error cerrando: %v", err)
	}

	_, err = env.service.ActualizarGramosFisico(ctx, usecaseconteos.ActualizarGramosFisicoInput{
		ConteoID: conteo.ID, ItemID: conteo.Items[0].ID, GramosFisico: mustDecimal(t, "99"), RequesterID: env.adminID,
	})
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestCerrarDosVecesEsConflicto(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	seedFragancia(t, env, "Chanel No. 5", "100.00")

	conteo, err := env.service.Crear(ctx, usecaseconteos.CrearInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecaseconteos.CerrarInput{TargetID: conteo.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error on first Cerrar: %v", err)
	}

	_, err = env.service.Cerrar(ctx, usecaseconteos.CerrarInput{TargetID: conteo.ID, RequesterID: env.adminID})
	assertCode(t, err, domainerrors.CodeConflict)
}
