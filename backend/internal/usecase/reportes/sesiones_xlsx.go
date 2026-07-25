package reportes

import (
	"bytes"
	"context"
	"time"

	domainreportes "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/reportes"
	reporterepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
)

// GenerarSesiones builds the sesiones laborales report (Resumen por vendedora,
// Detalle de sesiones sheets) for the resolved range.
func (s *Service) GenerarSesiones(ctx context.Context, sedeID int64, params domainreportes.ReporteParams) ([]byte, error) {
	desde, hasta, err := s.resolverRango(params)
	if err != nil {
		return nil, err
	}
	usuarioID := int64(0)
	if params.UsuarioID != nil {
		usuarioID = *params.UsuarioID
	}
	filtro := reporterepo.RangoFiltro{SedeID: sedeID, Desde: desde, Hasta: hasta, UsuarioID: usuarioID}

	resumen, err := s.repo.SesionesResumen(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	detalle, err := s.repo.SesionesDetalle(ctx, filtro)
	if err != nil {
		return nil, wrapErr(err)
	}
	if err := checkCtx(ctx); err != nil {
		return nil, err
	}

	b := NewXLSXBuilder(s.loc)

	// --- Hoja "Resumen por vendedora" ---
	rh := b.NewSheet("Resumen por vendedora")
	rh.WriteHeaders([]string{"Vendedora", "Total horas", "Días trabajados", "Promedio horas por día", "Sesiones"})
	for _, r := range resumen {
		var promedio time.Duration
		if r.TotalHoras != nil && r.DiasTrabajados > 0 {
			promedio = *r.TotalHoras / time.Duration(r.DiasTrabajados)
		}
		rh.WriteRow(
			r.NombreCompleto,
			r.TotalHoras,
			r.DiasTrabajados,
			promedio,
			r.SesionesCount,
		)
	}
	rh.AutoWidth()

	// --- Hoja "Detalle de sesiones" ---
	dh := b.NewSheet("Detalle de sesiones")
	dh.WriteHeaders([]string{"Vendedora", "Fecha entrada", "Hora entrada", "Fecha salida", "Hora salida", "Horas trabajadas"})
	for _, d := range detalle {
		var fechaSalida, horaSalida interface{}
		if d.SalidaAt != nil {
			fechaSalida = Fecha(*d.SalidaAt)
			horaSalida = Hora(*d.SalidaAt)
		}
		dh.WriteRow(
			d.NombreCompleto,
			Fecha(d.EntradaAt),
			Hora(d.EntradaAt),
			fechaSalida,
			horaSalida,
			d.HorasTrabajadas,
		)
	}
	dh.AutoWidth()

	b.DeleteDefaultSheet()

	var buf bytes.Buffer
	if err := b.Render(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
