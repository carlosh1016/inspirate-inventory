package auth

import (
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

type usuarioResponse struct {
	ID             int64      `json:"id"`
	SedeID         int64      `json:"sede_id"`
	NombreCompleto string     `json:"nombre_completo"`
	Correo         string     `json:"correo"`
	Rol            string     `json:"rol"`
	IsActive       bool       `json:"is_active"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

func toUsuarioResponse(u generated.Usuario) usuarioResponse {
	var lastLogin *time.Time
	if u.LastLoginAt.Valid {
		t := u.LastLoginAt.Time
		lastLogin = &t
	}

	return usuarioResponse{
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
