package usuarios

import (
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// UsuarioResponse never includes password_hash or deleted_at.
type UsuarioResponse struct {
	ID             int64      `json:"id"`
	SedeID         int64      `json:"sede_id"`
	NombreCompleto string     `json:"nombre_completo"`
	Correo         string     `json:"correo"`
	Rol            string     `json:"rol"`
	IsActive       bool       `json:"is_active"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

func toUsuarioResponse(u generated.Usuario) UsuarioResponse {
	var lastLogin *time.Time
	if u.LastLoginAt.Valid {
		t := u.LastLoginAt.Time
		lastLogin = &t
	}

	return UsuarioResponse{
		ID:             u.ID,
		SedeID:         u.SedeID,
		NombreCompleto: u.NombreCompleto,
		Correo:         u.Correo,
		Rol:            string(u.Rol),
		IsActive:       u.IsActive,
		LastLoginAt:    lastLogin,
		CreatedAt:      u.CreatedAt.Time,
	}
}

// CreateUsuarioRequest is the payload for POST /usuarios.
type CreateUsuarioRequest struct {
	NombreCompleto string `json:"nombre_completo" validate:"required,min=3,max=200"`
	Correo         string `json:"correo" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8,max=100"`
	Rol            string `json:"rol" validate:"required,oneof=admin vendedora"`
}

// UpdateUsuarioRequest is the payload for PATCH /usuarios/:id. A nil field
// means "leave unchanged".
type UpdateUsuarioRequest struct {
	NombreCompleto *string `json:"nombre_completo,omitempty" validate:"omitempty,min=3,max=200"`
	Correo         *string `json:"correo,omitempty" validate:"omitempty,email"`
	Rol            *string `json:"rol,omitempty" validate:"omitempty,oneof=admin vendedora"`
}

// UpdatePasswordRequest is the payload for PATCH /usuarios/:id/password.
type UpdatePasswordRequest struct {
	PasswordActual string `json:"password_actual" validate:"omitempty"`
	PasswordNueva  string `json:"password_nueva" validate:"required,min=8,max=100"`
}
