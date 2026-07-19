package refreshtokens

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

func (r *postgresRepository) Insert(ctx context.Context, usuarioID int64, tokenHash, ip, userAgent string, expiresAt time.Time) (generated.RefreshToken, error) {
	return r.q.InsertRefreshToken(ctx, generated.InsertRefreshTokenParams{
		UsuarioID: usuarioID,
		TokenHash: tokenHash,
		IpOrigen:  repo.InetPtr(ip),
		UserAgent: repo.Text(userAgent),
		ExpiresAt: repo.Timestamptz(expiresAt),
	})
}

func (r *postgresRepository) GetByHash(ctx context.Context, tokenHash string) (generated.RefreshToken, error) {
	rt, err := r.q.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return generated.RefreshToken{}, ErrNotFound
		}
		return generated.RefreshToken{}, err
	}
	return rt, nil
}

func (r *postgresRepository) Revoke(ctx context.Context, id int64) error {
	return r.q.RevokeRefreshToken(ctx, id)
}

func (r *postgresRepository) RevokeAllByUser(ctx context.Context, usuarioID int64) error {
	return r.q.RevokeAllRefreshTokensByUser(ctx, usuarioID)
}
