package usuarios

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

const defaultPageSize = 20

type postgresRepository struct {
	q *generated.Queries
}

// NewPostgres builds a Repository backed by Postgres via sqlc/pgx. db may be
// a *pgxpool.Pool or a pgx.Tx, so callers can run within a transaction.
func NewPostgres(db generated.DBTX) Repository {
	return &postgresRepository{q: generated.New(db)}
}

func (r *postgresRepository) GetByCorreo(ctx context.Context, correo string) (generated.Usuario, error) {
	u, err := r.q.GetUsuarioByCorreo(ctx, correo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Usuario{}, ErrNotFound
		}
		return generated.Usuario{}, err
	}
	return u, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (generated.Usuario, error) {
	u, err := r.q.GetUsuarioByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Usuario{}, ErrNotFound
		}
		return generated.Usuario{}, err
	}
	return u, nil
}

func (r *postgresRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	return r.q.UpdateLastLogin(ctx, id)
}

func (r *postgresRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	return r.q.UpdatePassword(ctx, generated.UpdatePasswordParams{
		PasswordHash: passwordHash,
		ID:           id,
	})
}

func (r *postgresRepository) ListPaginated(ctx context.Context, filter ListFilter) ([]generated.Usuario, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	offset := (page - 1) * pageSize

	rows, err := r.q.ListUsuariosPaginated(ctx, generated.ListUsuariosPaginatedParams{
		Limit:          int32(pageSize),
		Offset:         int32(offset),
		IncludeDeleted: filter.IncludeDeleted,
		Rol:            filter.Rol,
		Activo:         filter.Activo,
		Q:              filter.Q,
		SortCol:        filter.SortCol,
		SortDir:        filter.SortDir,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := r.q.CountUsuarios(ctx, generated.CountUsuariosParams{
		IncludeDeleted: filter.IncludeDeleted,
		Rol:            filter.Rol,
		Activo:         filter.Activo,
		Q:              filter.Q,
	})
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *postgresRepository) CountActiveAdmins(ctx context.Context) (int64, error) {
	return r.q.CountActiveAdmins(ctx)
}

func (r *postgresRepository) ExistsCorreo(ctx context.Context, correo string) (bool, error) {
	return r.q.ExistsCorreo(ctx, correo)
}

func (r *postgresRepository) Insert(ctx context.Context, sedeID int64, nombreCompleto, correo, passwordHash, rol string) (generated.Usuario, error) {
	return r.q.InsertUsuario(ctx, generated.InsertUsuarioParams{
		SedeID:         sedeID,
		NombreCompleto: nombreCompleto,
		Correo:         correo,
		PasswordHash:   passwordHash,
		Rol:            generated.RolEnum(rol),
	})
}

func (r *postgresRepository) Update(ctx context.Context, id int64, fields UpdateFields) (generated.Usuario, error) {
	params := generated.UpdateUsuarioParams{ID: id}
	if fields.NombreCompleto != nil {
		params.NombreCompleto = pgtype.Text{String: *fields.NombreCompleto, Valid: true}
	}
	if fields.Correo != nil {
		params.Correo = pgtype.Text{String: *fields.Correo, Valid: true}
	}
	if fields.Rol != nil {
		params.Rol = generated.NullRolEnum{RolEnum: generated.RolEnum(*fields.Rol), Valid: true}
	}

	u, err := r.q.UpdateUsuario(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.Usuario{}, ErrNotFound
		}
		return generated.Usuario{}, err
	}
	return u, nil
}

func (r *postgresRepository) Activate(ctx context.Context, id int64) error {
	return r.q.ActivateUsuario(ctx, id)
}

func (r *postgresRepository) Deactivate(ctx context.Context, id int64) error {
	return r.q.DeactivateUsuario(ctx, id)
}

func (r *postgresRepository) SoftDelete(ctx context.Context, id int64) error {
	return r.q.SoftDeleteUsuario(ctx, id)
}
