package reportes_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainreportes "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/reportes"
	usecasereportes "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/reportes"
)

// --- XLSX inspection helpers ---

func openXLSX(t *testing.T, data []byte) *excelize.File {
	t.Helper()
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("opening generated xlsx: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f
}

func rawCell(t *testing.T, f *excelize.File, sheet, cell string) string {
	t.Helper()
	v, err := f.GetCellValue(sheet, cell, excelize.Options{RawCellValue: true})
	if err != nil {
		t.Fatalf("reading raw cell %s!%s: %v", sheet, cell, err)
	}
	return v
}

func cellValue(t *testing.T, f *excelize.File, sheet, cell string) string {
	t.Helper()
	v, err := f.GetCellValue(sheet, cell)
	if err != nil {
		t.Fatalf("reading cell %s!%s: %v", sheet, cell, err)
	}
	return v
}

func rowCount(t *testing.T, f *excelize.File, sheet string) int {
	t.Helper()
	rows, err := f.GetRows(sheet)
	if err != nil {
		t.Fatalf("reading rows of %s: %v", sheet, err)
	}
	return len(rows)
}

func ptrTime(t time.Time) *time.Time { return &t }

// rangoParams builds a "rango" report covering [now+desdeDias, now+hastaDias].
func (e *testEnv) rangoParams(desdeDias, hastaDias int) domainreportes.ReporteParams {
	now := time.Now().In(e.loc)
	return domainreportes.ReporteParams{
		Periodo:    domainreportes.PeriodoRango,
		FechaDesde: ptrTime(now.AddDate(0, 0, desdeDias)),
		FechaHasta: ptrTime(now.AddDate(0, 0, hastaDias)),
	}
}

// --- Tests ---

func TestGenerarVentasResumen(t *testing.T) {
	e := newTestEnv(t)
	prod := e.seedProducto(t, "Crema", 10000, 0)
	for _, total := range []float64{10000, 20000, 30000} {
		v := e.seedVenta(t, e.efectivoID, total)
		e.seedVentaItemProducto(t, v, prod, 1, total)
	}

	data, err := e.service.GenerarVentas(context.Background(), e.sedeID, e.rangoParams(-1, 1))
	if err != nil {
		t.Fatalf("GenerarVentas: %v", err)
	}
	f := openXLSX(t, data)

	if got := f.GetSheetList(); len(got) != 3 || got[0] != "Resumen" || got[1] != "Ventas" || got[2] != "Items" {
		t.Fatalf("unexpected sheets: %v", got)
	}
	if title := cellValue(t, f, "Resumen", "A1"); !bytes.Contains([]byte(title), []byte("Reporte de ventas")) {
		t.Errorf("A1 = %q, want it to contain 'Reporte de ventas'", title)
	}
	if got := rawCell(t, f, "Resumen", "B4"); got != "60000" {
		t.Errorf("Total de ventas (B4) = %q, want 60000", got)
	}
	if got := rawCell(t, f, "Resumen", "B5"); got != "3" {
		t.Errorf("Número de ventas (B5) = %q, want 3", got)
	}
	if got := rawCell(t, f, "Resumen", "B6"); got != "20000" {
		t.Errorf("Ticket promedio (B6) = %q, want 20000", got)
	}
	if got := rawCell(t, f, "Resumen", "B10"); got != "60000" {
		t.Errorf("Efectivo (B10) = %q, want 60000", got)
	}
	if got := rowCount(t, f, "Ventas"); got != 4 {
		t.Errorf("Ventas rows = %d, want 4 (header + 3)", got)
	}
	if got := rowCount(t, f, "Items"); got != 4 {
		t.Errorf("Items rows = %d, want 4 (header + 3)", got)
	}
}

func TestGenerarVentasRangoSinVentas(t *testing.T) {
	e := newTestEnv(t)

	data, err := e.service.GenerarVentas(context.Background(), e.sedeID, e.rangoParams(-10, -9))
	if err != nil {
		t.Fatalf("GenerarVentas: %v", err)
	}
	f := openXLSX(t, data)

	if got := rawCell(t, f, "Resumen", "B4"); got != "0" {
		t.Errorf("Total de ventas (B4) = %q, want 0", got)
	}
	if got := rowCount(t, f, "Ventas"); got != 1 {
		t.Errorf("Ventas rows = %d, want 1 (headers only)", got)
	}
	if got := rowCount(t, f, "Items"); got != 1 {
		t.Errorf("Items rows = %d, want 1 (headers only)", got)
	}
}

func TestGenerarStockAlertas(t *testing.T) {
	e := newTestEnv(t)
	frag := e.seedFragancia(t, "Amber Bajo", 100) // mínimo 100g
	e.seedStock(t, "fragancia", frag, "vitrina", 10)
	prod := e.seedProducto(t, "Crema Suficiente", 5000, 5) // mínimo 5 uds
	e.seedStock(t, "producto", prod, "vitrina", 50)

	data, err := e.service.GenerarStock(context.Background(), e.sedeID, usecasereportes.StockParams{})
	if err != nil {
		t.Fatalf("GenerarStock: %v", err)
	}
	f := openXLSX(t, data)

	if got := f.GetSheetList(); len(got) != 4 {
		t.Fatalf("unexpected sheets: %v", got)
	}
	if got := rowCount(t, f, "Alertas"); got != 2 {
		t.Errorf("Alertas rows = %d, want 2 (header + 1 fragancia bajo mínimo)", got)
	}
	if got := cellValue(t, f, "Alertas", "A2"); got != "Fragancia" {
		t.Errorf("Alertas A2 = %q, want 'Fragancia'", got)
	}
	if got := cellValue(t, f, "Alertas", "B2"); got != "Amber Bajo" {
		t.Errorf("Alertas B2 = %q, want 'Amber Bajo'", got)
	}
	if got := rawCell(t, f, "Alertas", "F2"); got != "90" {
		t.Errorf("Faltante (F2) = %q, want 90", got)
	}
}

func TestGenerarMovimientosFiltroTipo(t *testing.T) {
	e := newTestEnv(t)
	prod := e.seedProducto(t, "Crema", 10000, 0)
	e.seedMovimiento(t, "venta", "producto", prod, "vitrina", 1, "", nil)
	e.seedMovimiento(t, "entrada_mercancia", "producto", prod, "bodega", 10, "", nil)

	// Sin filtro: ambos movimientos.
	all, err := e.service.GenerarMovimientos(context.Background(), e.sedeID, e.rangoParams(-1, 1), usecasereportes.MovimientosFiltros{})
	if err != nil {
		t.Fatalf("GenerarMovimientos (all): %v", err)
	}
	if got := rowCount(t, openXLSX(t, all), "Movimientos"); got != 3 {
		t.Errorf("Movimientos rows (sin filtro) = %d, want 3 (header + 2)", got)
	}

	// Filtro tipo=venta: solo uno.
	filtered, err := e.service.GenerarMovimientos(context.Background(), e.sedeID, e.rangoParams(-1, 1), usecasereportes.MovimientosFiltros{Tipo: "venta"})
	if err != nil {
		t.Fatalf("GenerarMovimientos (filtered): %v", err)
	}
	f := openXLSX(t, filtered)
	if got := rowCount(t, f, "Movimientos"); got != 2 {
		t.Errorf("Movimientos rows (tipo=venta) = %d, want 2 (header + 1)", got)
	}
	if got := cellValue(t, f, "Movimientos", "C2"); got != "venta" {
		t.Errorf("Movimientos C2 = %q, want 'venta'", got)
	}
}

func TestGenerarCuadresSoloCerrados(t *testing.T) {
	e := newTestEnv(t)
	now := time.Now().In(e.loc)
	e.seedCuadre(t, "cerrado", now, 500000)
	e.seedCuadre(t, "abierto", now.AddDate(0, 0, -1), 200000)

	data, err := e.service.GenerarCuadres(context.Background(), e.sedeID, e.rangoParams(-2, 1))
	if err != nil {
		t.Fatalf("GenerarCuadres: %v", err)
	}
	f := openXLSX(t, data)

	if got := f.GetSheetList(); len(got) != 3 || got[0] != "Cuadres" {
		t.Fatalf("unexpected sheets: %v", got)
	}
	if got := rowCount(t, f, "Cuadres"); got != 2 {
		t.Errorf("Cuadres rows = %d, want 2 (header + 1 cerrado)", got)
	}
	if got := cellValue(t, f, "Cuadres", "B2"); got != "cerrado" {
		t.Errorf("Cuadres B2 (estado) = %q, want 'cerrado'", got)
	}
}

func TestGenerarSesionesResumenPromedio(t *testing.T) {
	e := newTestEnv(t)
	now := time.Now().In(e.loc)
	day := func(offset, hour int) time.Time {
		d := now.AddDate(0, 0, offset)
		return time.Date(d.Year(), d.Month(), d.Day(), hour, 0, 0, 0, e.loc)
	}
	e.seedSesionCerrada(t, e.vendedoraID, day(-2, 8), day(-2, 13)) // 5h
	e.seedSesionCerrada(t, e.vendedoraID, day(-1, 8), day(-1, 13)) // 5h

	data, err := e.service.GenerarSesiones(context.Background(), e.sedeID, e.rangoParams(-5, 1))
	if err != nil {
		t.Fatalf("GenerarSesiones: %v", err)
	}
	f := openXLSX(t, data)

	if got := cellValue(t, f, "Resumen por vendedora", "B2"); got != "10:00:00" {
		t.Errorf("Total horas (B2) = %q, want 10:00:00", got)
	}
	if got := rawCell(t, f, "Resumen por vendedora", "C2"); got != "2" {
		t.Errorf("Días trabajados (C2) = %q, want 2", got)
	}
	if got := cellValue(t, f, "Resumen por vendedora", "D2"); got != "05:00:00" {
		t.Errorf("Promedio horas/día (D2) = %q, want 05:00:00", got)
	}
	if got := rowCount(t, f, "Detalle de sesiones"); got != 3 {
		t.Errorf("Detalle rows = %d, want 3 (header + 2)", got)
	}
}

func TestGenerarVentasTimeoutContextCancelado(t *testing.T) {
	e := newTestEnv(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := e.service.GenerarVentas(ctx, e.sedeID, e.rangoParams(-1, 1))
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
	var de *domainerrors.DomainError
	if !errors.As(err, &de) || de.Code != domainerrors.CodeGatewayTimeout {
		t.Fatalf("expected CodeGatewayTimeout DomainError, got %v", err)
	}
}
