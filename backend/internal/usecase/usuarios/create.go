package usuarios

import (
	"context"

	domainauth "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// CreateInput is the request payload plus the requesting admin's context.
// SedeID is inherited from the requester's own claims, not client input.
type CreateInput struct {
	SedeID         int64
	NombreCompleto string
	Correo         string
	Password       string
	Rol            string
	RequesterID    int64
	IP             string
	UserAgent      string
}

// Create registers a new usuario after checking the password policy and
// correo uniqueness.
func (s *Service) Create(ctx context.Context, in CreateInput) (generated.Usuario, error) {
	if err := domainauth.ValidatePassword(in.Password); err != nil {
		return generated.Usuario{}, err
	}

	exists, err := s.Usuarios.ExistsCorreo(ctx, in.Correo)
	if err != nil {
		return generated.Usuario{}, internalErr(err)
	}
	if exists {
		return generated.Usuario{}, domainerrors.NewConflict("Correo en uso", "Ya existe un usuario con ese correo.")
	}

	hash, err := domainauth.HashPassword(in.Password)
	if err != nil {
		return generated.Usuario{}, internalErr(err)
	}

	user, err := s.Usuarios.Insert(ctx, in.SedeID, in.NombreCompleto, in.Correo, hash, in.Rol)
	if err != nil {
		return generated.Usuario{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "usuario_creado", in.IP, in.UserAgent, &user.ID, nil, snapshot(user))

	return user, nil
}
