package sesiones

import (
	"context"
	"errors"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainsesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/sesiones"
	sesionesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/sesiones"
)

// Salida closes usuarioID's currently open sesion, computing
// horas_trabajadas in SQL as salida_at - entrada_at.
func (s *Service) Salida(ctx context.Context, usuarioID int64) (*domainsesiones.Sesion, error) {
	abierta, err := s.Sesiones.GetAbiertaPorUsuario(ctx, usuarioID)
	if err != nil {
		if errors.Is(err, sesionesrepo.ErrNotFound) {
			return nil, domainerrors.NewNotFound("Sesión no encontrada", "No tienes ninguna sesión abierta.")
		}
		return nil, internalErr(err)
	}

	row, err := s.Sesiones.Cerrar(ctx, abierta.ID, time.Now())
	if err != nil {
		return nil, internalErr(err)
	}

	sesion := toDomainSesion(row)
	return &sesion, nil
}
