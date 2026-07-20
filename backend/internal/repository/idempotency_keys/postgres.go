package idempotencykeys

import (
	"context"
	"errors"

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

func (r *postgresRepository) GetByKey(ctx context.Context, key string) (generated.IdempotencyKey, error) {
	row, err := r.q.GetIdempotencyKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.IdempotencyKey{}, ErrNotFound
		}
		return generated.IdempotencyKey{}, err
	}
	return row, nil
}

func (r *postgresRepository) Insert(ctx context.Context, params InsertParams) error {
	return r.q.InsertIdempotencyKey(ctx, generated.InsertIdempotencyKeyParams{
		Key:            params.Key,
		UsuarioID:      params.UsuarioID,
		Endpoint:       params.Endpoint,
		ResponseStatus: params.ResponseStatus,
		ResponseBody:   params.ResponseBody,
		RequestHash:    params.RequestHash,
		ExpiresAt:      repo.Timestamptz(params.ExpiresAt),
	})
}

func (r *postgresRepository) DeleteExpired(ctx context.Context) error {
	return r.q.DeleteExpiredIdempotencyKeys(ctx)
}
