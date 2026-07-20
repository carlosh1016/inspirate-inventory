package cuadres

import (
	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// toDomainCuadre maps the base cuadres_caja columns. CerradoPor is left nil
// — callers that read via a joined row (Get*/List*) attach it separately
// from that row's CerradoPorNombre.
func toDomainCuadre(row generated.CuadresCaja) domaincuadres.Cuadre {
	return domaincuadres.Cuadre{
		ID:                  row.ID,
		SedeID:              row.SedeID,
		Fecha:               row.Fecha.Time,
		Estado:              domaincuadres.EstadoCuadre(row.Estado),
		FondoBase:           row.FondoBase,
		TotalEfectivo:       row.TotalEfectivo,
		TotalNequi:          row.TotalNequi,
		TotalDaviplata:      row.TotalDaviplata,
		TotalTransferencia:  row.TotalTransferencia,
		TotalOtros:          row.TotalOtros,
		TotalPagos:          row.TotalPagos,
		TotalConsignaciones: row.TotalConsignaciones,
		ValorTurno:          row.ValorTurno,
		SaldoCalculado:      row.SaldoCalculado,
		Observaciones:       repo.StringPtr(row.Observaciones),
		CerradoPorUsuarioID: repo.Int8Ptr(row.CerradoPorUsuarioID),
		CerradoAt:           repo.TimePtr(row.CerradoAt),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}
}

func cuadresCajaFromGetByIDRow(r generated.GetCuadreByIDRow) generated.CuadresCaja {
	return generated.CuadresCaja{
		ID: r.ID, SedeID: r.SedeID, Fecha: r.Fecha, Estado: r.Estado, FondoBase: r.FondoBase,
		TotalEfectivo: r.TotalEfectivo, TotalNequi: r.TotalNequi, TotalDaviplata: r.TotalDaviplata,
		TotalTransferencia: r.TotalTransferencia, TotalOtros: r.TotalOtros, TotalPagos: r.TotalPagos,
		TotalConsignaciones: r.TotalConsignaciones, ValorTurno: r.ValorTurno, SaldoCalculado: r.SaldoCalculado,
		Observaciones: r.Observaciones, CerradoPorUsuarioID: r.CerradoPorUsuarioID, CerradoAt: r.CerradoAt,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func cuadresCajaFromGetBySedeFechaRow(r generated.GetCuadreBySedeFechaRow) generated.CuadresCaja {
	return generated.CuadresCaja{
		ID: r.ID, SedeID: r.SedeID, Fecha: r.Fecha, Estado: r.Estado, FondoBase: r.FondoBase,
		TotalEfectivo: r.TotalEfectivo, TotalNequi: r.TotalNequi, TotalDaviplata: r.TotalDaviplata,
		TotalTransferencia: r.TotalTransferencia, TotalOtros: r.TotalOtros, TotalPagos: r.TotalPagos,
		TotalConsignaciones: r.TotalConsignaciones, ValorTurno: r.ValorTurno, SaldoCalculado: r.SaldoCalculado,
		Observaciones: r.Observaciones, CerradoPorUsuarioID: r.CerradoPorUsuarioID, CerradoAt: r.CerradoAt,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func cuadresCajaFromListRow(r generated.ListCuadresPaginatedRow) generated.CuadresCaja {
	return generated.CuadresCaja{
		ID: r.ID, SedeID: r.SedeID, Fecha: r.Fecha, Estado: r.Estado, FondoBase: r.FondoBase,
		TotalEfectivo: r.TotalEfectivo, TotalNequi: r.TotalNequi, TotalDaviplata: r.TotalDaviplata,
		TotalTransferencia: r.TotalTransferencia, TotalOtros: r.TotalOtros, TotalPagos: r.TotalPagos,
		TotalConsignaciones: r.TotalConsignaciones, ValorTurno: r.ValorTurno, SaldoCalculado: r.SaldoCalculado,
		Observaciones: r.Observaciones, CerradoPorUsuarioID: r.CerradoPorUsuarioID, CerradoAt: r.CerradoAt,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

// attachCerradoPor sets c.CerradoPor from a joined row's cerrado_por_nombre
// column, when present.
func attachCerradoPor(c *domaincuadres.Cuadre, cerradoPorNombre string, valid bool) {
	if !valid || c.CerradoPorUsuarioID == nil {
		return
	}
	c.CerradoPor = &domaincuadres.UsuarioBrief{ID: *c.CerradoPorUsuarioID, NombreCompleto: cerradoPorNombre}
}

func toDomainPagoCaja(r generated.GetPagosByCuadreRow) domaincuadres.PagoCaja {
	return domaincuadres.PagoCaja{
		ID:           r.ID,
		CuadreCajaID: r.CuadreCajaID,
		UsuarioID:    r.UsuarioID,
		Concepto:     r.Concepto,
		Monto:        r.Monto,
		CreatedAt:    r.CreatedAt.Time,
		Usuario:      &domaincuadres.UsuarioBrief{ID: r.UsuarioID, NombreCompleto: r.UsuarioNombre},
	}
}

func toDomainConsignacion(r generated.GetConsignacionesByCuadreRow) domaincuadres.Consignacion {
	return domaincuadres.Consignacion{
		ID:           r.ID,
		CuadreCajaID: r.CuadreCajaID,
		UsuarioID:    r.UsuarioID,
		Monto:        r.Monto,
		Banco:        repo.StringPtr(r.Banco),
		Referencia:   repo.StringPtr(r.Referencia),
		CreatedAt:    r.CreatedAt.Time,
		Usuario:      &domaincuadres.UsuarioBrief{ID: r.UsuarioID, NombreCompleto: r.UsuarioNombre},
	}
}
