package sesiones

import (
	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	domainsesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/sesiones"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// toDomainSesion maps the base sesiones_laborales columns. Usuario is left
// nil — callers that read via a joined row (Get*/List*) attach it
// separately from that row's usuario_nombre.
func toDomainSesion(row generated.SesionesLaborale) domainsesiones.Sesion {
	return domainsesiones.Sesion{
		ID:              row.ID,
		SedeID:          row.SedeID,
		UsuarioID:       row.UsuarioID,
		EntradaAt:       row.EntradaAt.Time,
		SalidaAt:        repo.TimePtr(row.SalidaAt),
		HorasTrabajadas: repo.IntervalToDuration(row.HorasTrabajadas),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}

func sesionFromListRow(r generated.ListSesionesRow) domainsesiones.Sesion {
	s := toDomainSesion(generated.SesionesLaborale{
		ID: r.ID, SedeID: r.SedeID, UsuarioID: r.UsuarioID, EntradaAt: r.EntradaAt,
		SalidaAt: r.SalidaAt, HorasTrabajadas: r.HorasTrabajadas, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	})
	s.Usuario = &domaincuadres.UsuarioBrief{ID: r.UsuarioID, NombreCompleto: r.UsuarioNombre}
	return s
}
