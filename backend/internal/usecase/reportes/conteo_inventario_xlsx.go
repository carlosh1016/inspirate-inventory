package reportes

import (
	"bytes"
	"context"
	"errors"
	"time"

	domainconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/conteos"
	reporterepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
)

// ErrConteoNoEncontrado is returned when the requested conteoID doesn't
// belong to sedeID (or doesn't exist at all).
var ErrConteoNoEncontrado = errors.New("conteo no encontrado")

// GenerarConteoInventario builds the conteo mensual de inventario report (one
// sheet, one row per fragancia). It has no date range — it's a snapshot of
// one already-created conteo.
func (s *Service) GenerarConteoInventario(ctx context.Context, sedeID, conteoID int64) ([]byte, time.Time, error) {
	filtro := reporterepo.ConteoInventarioFiltro{SedeID: sedeID, ConteoID: conteoID}

	periodo, err := s.repo.ConteoInventarioPeriodo(ctx, filtro)
	if err != nil {
		if errors.Is(err, reporterepo.ErrNotFound) {
			return nil, time.Time{}, ErrConteoNoEncontrado
		}
		return nil, time.Time{}, wrapErr(err)
	}
	if err := checkCtx(ctx); err != nil {
		return nil, time.Time{}, err
	}

	items, err := s.repo.ConteoInventarioItems(ctx, filtro)
	if err != nil {
		return nil, time.Time{}, wrapErr(err)
	}

	b := NewXLSXBuilder(s.loc)
	sh := b.NewSheet("Conteo mensual")
	sh.WriteHeaders([]string{
		"Fragancia", "Saldo inicial (g)", "Gramos sistema (g)", "Gramos físico (g)",
		"Diferencia (g)", "% vendido", "Estado",
	})
	for _, it := range items {
		item := domainconteos.ConteoItem{
			FraganciaNombre: it.FraganciaNombre,
			SaldoInicial:    it.SaldoInicial,
			GramosSistema:   it.GramosSistema,
		}
		if it.GramosFisico.Valid {
			item.GramosFisico = &it.GramosFisico.Decimal
		}

		var gramosFisico interface{}
		if item.GramosFisico != nil {
			gramosFisico = Gramos(*item.GramosFisico)
		}
		var diferencia interface{}
		estado := "Sin contar"
		if d := item.Diferencia(); d != nil {
			diferencia = Gramos(*d)
			estado = "OK"
			if item.ExcedeLimite() {
				estado = "Revisar"
			}
		}
		var porcentaje interface{}
		if p := item.PorcentajeVendido(); p != nil {
			porcentaje = Porcentaje(p.Round(2))
		}

		sh.WriteRow(item.FraganciaNombre, Gramos(item.SaldoInicial), Gramos(item.GramosSistema), gramosFisico, diferencia, porcentaje, estado)
	}
	sh.AutoWidth()

	b.DeleteDefaultSheet()

	var buf bytes.Buffer
	if err := b.Render(&buf); err != nil {
		return nil, time.Time{}, err
	}
	return buf.Bytes(), periodo, nil
}
