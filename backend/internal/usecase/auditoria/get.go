package auditoria

import (
	"context"
	"errors"

	domainauditoria "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auditoria"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	auditoriarepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
)

// Get returns a single audit evento by id, or a 404 DomainError when absent.
func (s *Service) Get(ctx context.Context, id int64) (domainauditoria.Evento, error) {
	ev, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, auditoriarepo.ErrNotFound) {
			return domainauditoria.Evento{}, domainerrors.NewNotFound(
				"Evento no encontrado", "No existe un evento de auditoría con ese identificador.")
		}
		return domainauditoria.Evento{}, err
	}
	return ev, nil
}
