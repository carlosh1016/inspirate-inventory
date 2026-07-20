package sesiones

import (
	"context"
	"errors"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainsesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/sesiones"
	sesionesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/sesiones"
)

// EntradaInput is the request context for POST /sesiones-laborales/entrada.
type EntradaInput struct {
	SedeID    int64
	UsuarioID int64
}

// Entrada opens a new sesion for UsuarioID. A usuario may only have one
// open sesion at a time (across every sede) — attempting a second entrada
// without closing the first is a 409.
func (s *Service) Entrada(ctx context.Context, in EntradaInput) (*domainsesiones.Sesion, error) {
	if _, err := s.Sesiones.GetAbiertaPorUsuario(ctx, in.UsuarioID); err == nil {
		return nil, domainerrors.NewConflict("Sesión ya abierta", "Ya tienes una sesión abierta. Ciérrala primero.")
	} else if !errors.Is(err, sesionesrepo.ErrNotFound) {
		return nil, internalErr(err)
	}

	row, err := s.Sesiones.Insert(ctx, in.SedeID, in.UsuarioID, time.Now())
	if err != nil {
		return nil, internalErr(err)
	}

	sesion := toDomainSesion(row)
	return &sesion, nil
}
