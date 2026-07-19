package metodospago

import (
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// MetodoPagoResponse never includes deleted_at.
type MetodoPagoResponse struct {
	ID        int64     `json:"id"`
	Nombre    string    `json:"nombre"`
	Codigo    string    `json:"codigo"`
	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toMetodoPagoResponse(m generated.MetodosPago) MetodoPagoResponse {
	return MetodoPagoResponse{
		ID:        m.ID,
		Nombre:    m.Nombre,
		Codigo:    m.Codigo,
		Activo:    m.Activo,
		CreatedAt: m.CreatedAt.Time,
		UpdatedAt: m.UpdatedAt.Time,
	}
}

// CreateMetodoPagoRequest is the payload for POST /metodos-pago.
type CreateMetodoPagoRequest struct {
	Nombre string `json:"nombre" validate:"required,min=2,max=100"`
	Codigo string `json:"codigo" validate:"required,min=1,max=50"`
}

// UpdateMetodoPagoRequest is the payload for PATCH /metodos-pago/:id. A nil
// field means "leave unchanged".
type UpdateMetodoPagoRequest struct {
	Nombre *string `json:"nombre,omitempty" validate:"omitempty,min=2,max=100"`
	Codigo *string `json:"codigo,omitempty" validate:"omitempty,min=1,max=50"`
}
