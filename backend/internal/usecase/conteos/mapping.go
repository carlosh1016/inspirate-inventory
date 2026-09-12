package conteos

import (
	"github.com/jackc/pgx/v5/pgtype"

	domainconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/conteos"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// conteoBase is the shape shared by GetConteoByIDRow and
// ListConteosPaginatedRow — both join usuarios twice (creado_por/
// cerrado_por) on top of the base conteos_inventario columns.
type conteoBase struct {
	ID                  int64
	SedeID              int64
	Periodo             pgtype.Date
	Estado              generated.EstadoConteoEnum
	CreadoPorUsuarioID  int64
	CerradoPorUsuarioID pgtype.Int8
	CerradoAt           pgtype.Timestamptz
	CreatedAt           pgtype.Timestamptz
	UpdatedAt           pgtype.Timestamptz
	CreadoPorNombre     string
	CerradoPorNombre    pgtype.Text
}

func toDomainConteo(row conteoBase) domainconteos.ConteoInventario {
	c := domainconteos.ConteoInventario{
		ID:        row.ID,
		SedeID:    row.SedeID,
		Periodo:   row.Periodo.Time,
		Estado:    domainconteos.EstadoConteo(row.Estado),
		CreadoPor: domainconteos.UsuarioBrief{ID: row.CreadoPorUsuarioID, NombreCompleto: row.CreadoPorNombre},
		CerradoAt: repo.TimePtr(row.CerradoAt),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if cerradoPorID := repo.Int8Ptr(row.CerradoPorUsuarioID); cerradoPorID != nil && row.CerradoPorNombre.Valid {
		c.CerradoPor = &domainconteos.UsuarioBrief{ID: *cerradoPorID, NombreCompleto: row.CerradoPorNombre.String}
	}
	return c
}

func conteoFromGetByIDRow(r generated.GetConteoByIDRow) domainconteos.ConteoInventario {
	return toDomainConteo(conteoBase{
		ID: r.ID, SedeID: r.SedeID, Periodo: r.Periodo, Estado: r.Estado,
		CreadoPorUsuarioID: r.CreadoPorUsuarioID, CerradoPorUsuarioID: r.CerradoPorUsuarioID,
		CerradoAt: r.CerradoAt, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
		CreadoPorNombre: r.CreadoPorNombre, CerradoPorNombre: r.CerradoPorNombre,
	})
}

func conteoFromListRow(r generated.ListConteosPaginatedRow) domainconteos.ConteoInventario {
	return toDomainConteo(conteoBase{
		ID: r.ID, SedeID: r.SedeID, Periodo: r.Periodo, Estado: r.Estado,
		CreadoPorUsuarioID: r.CreadoPorUsuarioID, CerradoPorUsuarioID: r.CerradoPorUsuarioID,
		CerradoAt: r.CerradoAt, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
		CreadoPorNombre: r.CreadoPorNombre, CerradoPorNombre: r.CerradoPorNombre,
	})
}

func toDomainConteoItem(r generated.ListConteoItemsRow) domainconteos.ConteoItem {
	return domainconteos.ConteoItem{
		ID:              r.ID,
		ConteoID:        r.ConteoID,
		FraganciaID:     r.FraganciaID,
		FraganciaNombre: r.FraganciaNombre,
		SaldoInicial:    r.SaldoInicial,
		GramosSistema:   r.GramosSistema,
		GramosFisico:    repo.NullDecimalPtr(r.GramosFisico),
	}
}

func itemFromUpdateRow(r generated.UpdateConteoItemFisicoRow) domainconteos.ConteoItem {
	return domainconteos.ConteoItem{
		ID:              r.ID,
		ConteoID:        r.ConteoID,
		FraganciaID:     r.FraganciaID,
		FraganciaNombre: r.FraganciaNombre,
		SaldoInicial:    r.SaldoInicial,
		GramosSistema:   r.GramosSistema,
		GramosFisico:    repo.NullDecimalPtr(r.GramosFisico),
	}
}
