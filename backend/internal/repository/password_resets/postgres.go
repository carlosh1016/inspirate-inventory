package passwordresets

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

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

func (r *postgresRepository) Insert(ctx context.Context, usuarioID int64, tokenHash string, expiresAt time.Time) (generated.PasswordReset, error) {
	return r.q.InsertPasswordReset(ctx, generated.InsertPasswordResetParams{
		UsuarioID: usuarioID,
		TokenHash: tokenHash,
		ExpiresAt: repo.Timestamptz(expiresAt),
	})
}

func (r *postgresRepository) GetByHash(ctx context.Context, tokenHash string) (generated.PasswordReset, error) {
	pr, err := r.q.GetPasswordResetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.PasswordReset{}, ErrNotFound
		}
		return generated.PasswordReset{}, err
	}
	return pr, nil
}

func (r *postgresRepository) MarkUsed(ctx context.Context, id int64) error {
	return r.q.MarkPasswordResetUsed(ctx, id)
}
