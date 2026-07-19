package metodospago

import (
	"context"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// CreateInput is the request payload plus the requester's context. Create
// is admin-only, enforced at the router.
type CreateInput struct {
	Nombre      string
	Codigo      string
	RequesterID int64
	IP          string
	UserAgent   string
}

// Create registers a new metodo_pago, enforcing nombre and codigo
// uniqueness.
func (s *Service) Create(ctx context.Context, in CreateInput) (generated.MetodosPago, error) {
	existsCodigo, err := s.MetodosPago.ExistsCodigo(ctx, in.Codigo, 0)
	if err != nil {
		return generated.MetodosPago{}, internalErr(err)
	}
	if existsCodigo {
		return generated.MetodosPago{}, domainerrors.NewConflict(
			"Código en uso",
			"Ya existe un método de pago con ese código.",
		)
	}

	existsNombre, err := s.MetodosPago.ExistsNombre(ctx, in.Nombre, 0)
	if err != nil {
		return generated.MetodosPago{}, internalErr(err)
	}
	if existsNombre {
		return generated.MetodosPago{}, domainerrors.NewConflict(
			"Nombre en uso",
			"Ya existe un método de pago con ese nombre.",
		)
	}

	m, err := s.MetodosPago.Insert(ctx, in.Nombre, in.Codigo)
	if err != nil {
		return generated.MetodosPago{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "metodo_pago_creado", in.IP, in.UserAgent, &m.ID, nil, snapshot(m))

	return m, nil
}
