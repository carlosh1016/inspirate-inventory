package auditoria

import (
	"context"
	"encoding/json"
	"errors"
	"net/netip"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	domainauditoria "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auditoria"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) Insert(ctx context.Context, e Entry) error {
	var tablaAfectada string
	if e.TablaAfectada != nil {
		tablaAfectada = *e.TablaAfectada
	}

	return r.q.InsertAuditoria(ctx, generated.InsertAuditoriaParams{
		UsuarioID:     repo.Int8(e.UsuarioID),
		Accion:        e.Accion,
		TablaAfectada: repo.Text(tablaAfectada),
		RegistroID:    repo.Int8(e.RegistroID),
		DatosAntes:    e.DatosAntes,
		DatosDespues:  e.DatosDespues,
		Ip:            repo.InetPtr(e.IP),
		UserAgent:     repo.Text(e.UserAgent),
	})
}

func (r *postgresRepository) List(ctx context.Context, f ListFiltro) ([]domainauditoria.Evento, error) {
	rows, err := r.q.ListAuditoriaPaginated(ctx, generated.ListAuditoriaPaginatedParams{
		Limit:         f.Limit,
		Offset:        f.Offset,
		UsuarioID:     f.UsuarioID,
		Accion:        f.Accion,
		TablaAfectada: f.TablaAfectada,
		FechaDesde:    repo.TimestamptzPtr(f.FechaDesde),
		FechaHasta:    repo.TimestamptzPtr(f.FechaHasta),
	})
	if err != nil {
		return nil, err
	}
	out := make([]domainauditoria.Evento, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapEvento(
			row.ID, row.UsuarioID, row.Accion, row.TablaAfectada, row.RegistroID,
			row.DatosAntes, row.DatosDespues, row.Ip, row.UserAgent, row.CreatedAt, row.UsuarioNombre,
		))
	}
	return out, nil
}

func (r *postgresRepository) Count(ctx context.Context, f ListFiltro) (int64, error) {
	return r.q.CountAuditoria(ctx, generated.CountAuditoriaParams{
		UsuarioID:     f.UsuarioID,
		Accion:        f.Accion,
		TablaAfectada: f.TablaAfectada,
		FechaDesde:    repo.TimestamptzPtr(f.FechaDesde),
		FechaHasta:    repo.TimestamptzPtr(f.FechaHasta),
	})
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (domainauditoria.Evento, error) {
	row, err := r.q.GetAuditoriaByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainauditoria.Evento{}, ErrNotFound
		}
		return domainauditoria.Evento{}, err
	}
	return mapEvento(
		row.ID, row.UsuarioID, row.Accion, row.TablaAfectada, row.RegistroID,
		row.DatosAntes, row.DatosDespues, row.Ip, row.UserAgent, row.CreatedAt, row.UsuarioNombre,
	), nil
}

func (r *postgresRepository) AccionesDistintas(ctx context.Context) ([]string, error) {
	return r.q.GetAccionesDistintas(ctx)
}

// mapEvento builds a domain Evento from the shared column set of the
// list/get rows (both select a.* plus usuario_nombre).
func mapEvento(
	id int64,
	usuarioID pgtype.Int8,
	accion string,
	tablaAfectada pgtype.Text,
	registroID pgtype.Int8,
	datosAntes, datosDespues []byte,
	ip *netip.Addr,
	userAgent pgtype.Text,
	createdAt pgtype.Timestamptz,
	usuarioNombre pgtype.Text,
) domainauditoria.Evento {
	ev := domainauditoria.Evento{
		ID:            id,
		UsuarioID:     repo.Int8Ptr(usuarioID),
		Accion:        accion,
		TablaAfectada: repo.StringPtr(tablaAfectada),
		RegistroID:    repo.Int8Ptr(registroID),
		DatosAntes:    json.RawMessage(datosAntes),
		DatosDespues:  json.RawMessage(datosDespues),
		IP:            repo.InetString(ip),
		UserAgent:     repo.StringPtr(userAgent),
		CreatedAt:     createdAt.Time,
	}
	if usuarioID.Valid && usuarioNombre.Valid {
		ev.Usuario = &domainauditoria.UsuarioBrief{ID: usuarioID.Int64, NombreCompleto: usuarioNombre.String}
	}
	return ev
}
