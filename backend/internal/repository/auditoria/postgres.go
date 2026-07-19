package auditoria

import (
	"context"

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
