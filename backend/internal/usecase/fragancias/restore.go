package fragancias

import (
	"context"
	"errors"

	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// RestoreInput is the request context: who is restoring what.
type RestoreInput struct {
	TargetID    int64
	RequesterID int64
	IP          string
	UserAgent   string
}

// Restore un-deletes a previously soft-deleted fragancia. The underlying
// query's WHERE deleted_at IS NOT NULL means "not found" covers both a
// missing id and one that's already active.
func (s *Service) Restore(ctx context.Context, in RestoreInput) (generated.GetFraganciaByIDRow, error) {
	restored, err := s.Fragancias.Restore(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetFraganciaByIDRow{}, notFoundErr()
		}
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}

	result, err := s.Fragancias.GetByID(ctx, restored.ID)
	if err != nil {
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "fragancia_restaurada", in.IP, in.UserAgent, &restored.ID, nil, snapshot(restored))

	return result, nil
}
