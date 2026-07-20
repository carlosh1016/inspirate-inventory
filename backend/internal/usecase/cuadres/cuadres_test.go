package cuadres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	usecasecuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/cuadres"
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

// --- 1/2/3. Abrir ---

func TestAbrirCuadreExito(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{
		SedeID: env.sedeID, RequesterID: env.adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Cuadre.Estado != domaincuadres.EstadoAbierto {
		t.Fatalf("expected abierto, got %v", out.Cuadre.Estado)
	}
	if out.Cuadre.FondoBase.String() != "100000" {
		t.Fatalf("expected default fondo_base 100000, got %s", out.Cuadre.FondoBase)
	}
	if len(out.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %+v", out.Warnings)
	}
}

func TestAbrirCuadreYaExisteAbiertoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error opening the first time: %v", err)
	}

	_, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestAbrirCuadreYaExisteCerradoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	_, err = env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

func TestAbrirCuadreConAnteriorAbiertoDaWarning(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	ayer := time.Now().AddDate(0, 0, -1)
	if _, err := env.pool.Exec(ctx,
		`INSERT INTO cuadres_caja (sede_id, fecha, estado, fondo_base) VALUES ($1, $2::date, 'abierto', 100000)`,
		env.sedeID, ayer,
	); err != nil {
		t.Fatalf("seeding cuadre de ayer: %v", err)
	}

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Warnings) != 1 || out.Warnings[0].Codigo != "cuadre_anterior_abierto" {
		t.Fatalf("expected 1 cuadre_anterior_abierto warning, got %+v", out.Warnings)
	}
}

// --- 4/5. Cerrar ---

func TestCerrarCuadreCongelaTotales(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	seedVenta(t, env, env.metodoEfectivo, "50000")
	seedVenta(t, env, env.metodoNequi, "20000")
	if _, err := env.service.AddPago(ctx, usecasecuadres.AddPagoInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Concepto: "papel higiénico", Monto: mustDecimal(t, "5000"),
	}); err != nil {
		t.Fatalf("unexpected error adding pago: %v", err)
	}

	valorTurno := mustDecimal(t, "10000")
	cerrado, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{
		TargetID: out.Cuadre.ID, ValorTurno: &valorTurno, RequesterID: env.adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}
	if cerrado.Estado != domaincuadres.EstadoCerrado {
		t.Fatalf("expected cerrado, got %v", cerrado.Estado)
	}
	if cerrado.TotalEfectivo.String() != "50000" || cerrado.TotalNequi.String() != "20000" {
		t.Fatalf("expected totals frozen from ventas, got efectivo=%s nequi=%s", cerrado.TotalEfectivo, cerrado.TotalNequi)
	}
	// saldo = 100000 + 50000 (efectivo) - 5000 (pagos) - 0 (consignaciones) - 10000 (turno) = 135000
	if cerrado.SaldoCalculado.String() != "135000" {
		t.Fatalf("expected saldo_calculado=135000, got %s", cerrado.SaldoCalculado)
	}
	if cerrado.CerradoPor == nil || cerrado.CerradoPor.ID != env.adminID {
		t.Fatalf("expected cerrado_por populated with admin, got %+v", cerrado.CerradoPor)
	}
	if cerrado.CerradoAt == nil {
		t.Fatal("expected cerrado_at populated")
	}

	// Nueva venta después de cerrar no debe alterar el snapshot congelado.
	seedVenta(t, env, env.metodoEfectivo, "999999")
	reread, err := env.service.GetByID(ctx, out.Cuadre.ID)
	if err != nil {
		t.Fatalf("unexpected error re-reading: %v", err)
	}
	if reread.TotalEfectivo.String() != "50000" {
		t.Fatalf("expected frozen total_efectivo unaffected by later ventas, got %s", reread.TotalEfectivo)
	}
}

func TestCerrarCuadreYaCerradoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error on first close: %v", err)
	}

	_, err = env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID})
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

// --- 6/7/8. GetHoy ---

func TestGetHoySinCuadreRetornaNil(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	cuadre, err := env.service.GetHoy(ctx, env.sedeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cuadre != nil {
		t.Fatalf("expected nil cuadre, got %+v", cuadre)
	}
}

func TestGetHoyAbiertoRecalculaTotalesEnVivo(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	if _, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}

	cuadre, err := env.service.GetHoy(ctx, env.sedeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cuadre.TotalEfectivo.IsZero() {
		t.Fatalf("expected total_efectivo=0 before any venta, got %s", cuadre.TotalEfectivo)
	}

	seedVenta(t, env, env.metodoEfectivo, "30000")

	cuadre, err = env.service.GetHoy(ctx, env.sedeID)
	if err != nil {
		t.Fatalf("unexpected error re-reading: %v", err)
	}
	if cuadre.TotalEfectivo.String() != "30000" {
		t.Fatalf("expected total_efectivo=30000 live, got %s", cuadre.TotalEfectivo)
	}
}

func TestGetHoyCerradoRetornaSnapshot(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	seedVenta(t, env, env.metodoEfectivo, "15000")
	if _, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	cuadre, err := env.service.GetHoy(ctx, env.sedeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cuadre.Estado != domaincuadres.EstadoCerrado {
		t.Fatalf("expected cerrado, got %v", cuadre.Estado)
	}
	if cuadre.TotalEfectivo.String() != "15000" {
		t.Fatalf("expected frozen total_efectivo=15000, got %s", cuadre.TotalEfectivo)
	}
}

// --- 9/10/11/12. Pagos ---

func TestAddPagoMientrasAbiertoActualizaTotales(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}

	pago, err := env.service.AddPago(ctx, usecasecuadres.AddPagoInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Concepto: "papel higiénico", Monto: mustDecimal(t, "5000"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pago.Usuario == nil || pago.Usuario.ID != env.vendedoraID {
		t.Fatalf("expected pago.Usuario populated, got %+v", pago.Usuario)
	}

	cuadre, err := env.service.GetByID(ctx, out.Cuadre.ID)
	if err != nil {
		t.Fatalf("unexpected error re-reading: %v", err)
	}
	if cuadre.TotalPagos.String() != "5000" || len(cuadre.Pagos) != 1 {
		t.Fatalf("expected total_pagos=5000 with 1 pago, got total=%s len=%d", cuadre.TotalPagos, len(cuadre.Pagos))
	}
}

func TestAddPagoCuandoCerradoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	_, err = env.service.AddPago(ctx, usecasecuadres.AddPagoInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Concepto: "papel higiénico", Monto: mustDecimal(t, "5000"),
	})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestDeletePagoMientrasAbierto(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	pago, err := env.service.AddPago(ctx, usecasecuadres.AddPagoInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Concepto: "papel higiénico", Monto: mustDecimal(t, "5000"),
	})
	if err != nil {
		t.Fatalf("unexpected error adding pago: %v", err)
	}

	if err := env.service.DeletePago(ctx, out.Cuadre.ID, pago.ID); err != nil {
		t.Fatalf("unexpected error deleting: %v", err)
	}

	cuadre, err := env.service.GetByID(ctx, out.Cuadre.ID)
	if err != nil {
		t.Fatalf("unexpected error re-reading: %v", err)
	}
	if len(cuadre.Pagos) != 0 {
		t.Fatalf("expected 0 pagos after delete, got %d", len(cuadre.Pagos))
	}
}

func TestDeletePagoCuandoCerradoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	pago, err := env.service.AddPago(ctx, usecasecuadres.AddPagoInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Concepto: "papel higiénico", Monto: mustDecimal(t, "5000"),
	})
	if err != nil {
		t.Fatalf("unexpected error adding pago: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	err = env.service.DeletePago(ctx, out.Cuadre.ID, pago.ID)
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

// --- Consignaciones (mismo comportamiento que pagos) ---

func TestAddYDeleteConsignacionMientrasAbierto(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	banco := "Bancolombia"
	consig, err := env.service.AddConsignacion(ctx, usecasecuadres.AddConsignacionInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Monto: mustDecimal(t, "50000"), Banco: &banco,
	})
	if err != nil {
		t.Fatalf("unexpected error adding consignacion: %v", err)
	}

	cuadre, err := env.service.GetByID(ctx, out.Cuadre.ID)
	if err != nil {
		t.Fatalf("unexpected error re-reading: %v", err)
	}
	if cuadre.TotalConsignaciones.String() != "50000" || len(cuadre.Consignaciones) != 1 {
		t.Fatalf("expected total_consignaciones=50000 with 1 entry, got total=%s len=%d", cuadre.TotalConsignaciones, len(cuadre.Consignaciones))
	}

	if err := env.service.DeleteConsignacion(ctx, out.Cuadre.ID, consig.ID); err != nil {
		t.Fatalf("unexpected error deleting: %v", err)
	}
	cuadre, err = env.service.GetByID(ctx, out.Cuadre.ID)
	if err != nil {
		t.Fatalf("unexpected error re-reading after delete: %v", err)
	}
	if len(cuadre.Consignaciones) != 0 {
		t.Fatalf("expected 0 consignaciones after delete, got %d", len(cuadre.Consignaciones))
	}
}

func TestDeleteConsignacionCuandoCerradoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	consig, err := env.service.AddConsignacion(ctx, usecasecuadres.AddConsignacionInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Monto: mustDecimal(t, "50000"),
	})
	if err != nil {
		t.Fatalf("unexpected error adding consignacion: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	err = env.service.DeleteConsignacion(ctx, out.Cuadre.ID, consig.ID)
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

// --- CajaStatusService (consumido por usecase/ventas) ---

func TestCajaStatusServiceBloqueaVentaSiCuadreCerrado(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	cajaStatus := usecasecuadres.NewCajaStatusService(cuadresrepo.NewPostgres(env.pool))
	hoy := time.Now()
	err = cajaStatus.VerificarPuedeRegistrarVenta(ctx, env.sedeID, hoy)
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)
}

func TestCajaStatusServicePermiteVentaSiCuadreAbiertoOSinCuadre(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	cajaStatus := usecasecuadres.NewCajaStatusService(cuadresrepo.NewPostgres(env.pool))
	hoy := time.Now()

	if err := cajaStatus.VerificarPuedeRegistrarVenta(ctx, env.sedeID, hoy); err != nil {
		t.Fatalf("expected no error without any cuadre, got %v", err)
	}

	if _, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	if err := cajaStatus.VerificarPuedeRegistrarVenta(ctx, env.sedeID, hoy); err != nil {
		t.Fatalf("expected no error with an abierto cuadre, got %v", err)
	}
}

// --- List ---

func TestListCuadres(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	if _, err := env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, RequesterID: env.adminID}); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	result, err := env.service.List(ctx, usecasecuadres.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 1 || len(result.Cuadres) != 1 {
		t.Fatalf("expected 1 cuadre, got total=%d len=%d", result.Total, len(result.Cuadres))
	}
	if result.Cuadres[0].CerradoPor == nil || result.Cuadres[0].CerradoPor.ID != env.adminID {
		t.Fatalf("expected cerrado_por populated in list row, got %+v", result.Cuadres[0].CerradoPor)
	}

	filtered, err := env.service.List(ctx, usecasecuadres.ListInput{SedeID: env.sedeID, Estado: "abierto"})
	if err != nil {
		t.Fatalf("unexpected error filtering by estado: %v", err)
	}
	if filtered.Total != 0 {
		t.Fatalf("expected 0 cuadres abiertos, got %d", filtered.Total)
	}
}

// --- Validaciones y casos not-found ---

func TestGetByIDDesconocidoNotFound(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_, err := env.service.GetByID(ctx, 999999)
	if err == nil {
		t.Fatal("expected a not_found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestAbrirFondoBaseNegativoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	negativo := mustDecimal(t, "-1")
	_, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, FondoBase: &negativo, RequesterID: env.adminID})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestCerrarValorTurnoNegativoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	negativo := mustDecimal(t, "-1")
	_, err = env.service.Cerrar(ctx, usecasecuadres.CerrarInput{TargetID: out.Cuadre.ID, ValorTurno: &negativo, RequesterID: env.adminID})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestAddPagoMontoNoPositivoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	_, err = env.service.AddPago(ctx, usecasecuadres.AddPagoInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Concepto: "x", Monto: decimal.Zero,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestDeletePagoDesconocidoNotFound(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}

	err = env.service.DeletePago(ctx, out.Cuadre.ID, 999999)
	if err == nil {
		t.Fatal("expected a not_found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestDeletePagoDeOtroCuadreNotFound(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	pago, err := env.service.AddPago(ctx, usecasecuadres.AddPagoInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Concepto: "x", Monto: mustDecimal(t, "1000"),
	})
	if err != nil {
		t.Fatalf("unexpected error adding pago: %v", err)
	}

	err = env.service.DeletePago(ctx, 999999, pago.ID)
	if err == nil {
		t.Fatal("expected a not_found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestAddConsignacionMontoNoPositivoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}
	_, err = env.service.AddConsignacion(ctx, usecasecuadres.AddConsignacionInput{
		CuadreID: out.Cuadre.ID, UsuarioID: env.vendedoraID, Monto: decimal.Zero,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestDeleteConsignacionDesconocidaNotFound(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	out, err := env.service.Abrir(ctx, usecasecuadres.AbrirInput{SedeID: env.sedeID, RequesterID: env.adminID})
	if err != nil {
		t.Fatalf("unexpected error opening: %v", err)
	}

	err = env.service.DeleteConsignacion(ctx, out.Cuadre.ID, 999999)
	if err == nil {
		t.Fatal("expected a not_found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}
