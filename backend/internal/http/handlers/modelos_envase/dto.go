package modelosenvase

import (
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ModeloEnvaseResponse never includes deleted_at.
type ModeloEnvaseResponse struct {
	ID                 int64     `json:"id"`
	Tipo               string    `json:"tipo"`
	TamanoOz           string    `json:"tamano_oz"`
	EquivGramos        string    `json:"equiv_gramos"`
	PrecioSolo         string    `json:"precio_solo"`
	PrecioConFragancia string    `json:"precio_con_fragancia"`
	PrecioRecarga      string    `json:"precio_recarga"`
	Activo             bool      `json:"activo"`
	TieneVariantes     bool      `json:"tiene_variantes"`
	VariantesActivas   int64     `json:"variantes_activas"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func toModeloEnvaseResponseFromGet(m generated.GetModeloEnvaseByIDRow) ModeloEnvaseResponse {
	return ModeloEnvaseResponse{
		ID:                 m.ID,
		Tipo:               m.Tipo,
		TamanoOz:           m.TamanoOz.String(),
		EquivGramos:        m.EquivGramos.String(),
		PrecioSolo:         m.PrecioSolo.String(),
		PrecioConFragancia: m.PrecioConFragancia.String(),
		PrecioRecarga:      m.PrecioRecarga.String(),
		Activo:             m.Activo,
		TieneVariantes:     m.TieneVariantes,
		VariantesActivas:   m.VariantesActivas,
		CreatedAt:          m.CreatedAt.Time,
		UpdatedAt:          m.UpdatedAt.Time,
	}
}

func toModeloEnvaseResponseFromList(m generated.ListModelosEnvasePaginatedRow) ModeloEnvaseResponse {
	return ModeloEnvaseResponse{
		ID:                 m.ID,
		Tipo:               m.Tipo,
		TamanoOz:           m.TamanoOz.String(),
		EquivGramos:        m.EquivGramos.String(),
		PrecioSolo:         m.PrecioSolo.String(),
		PrecioConFragancia: m.PrecioConFragancia.String(),
		PrecioRecarga:      m.PrecioRecarga.String(),
		Activo:             m.Activo,
		TieneVariantes:     m.TieneVariantes,
		VariantesActivas:   m.VariantesActivas,
		CreatedAt:          m.CreatedAt.Time,
		UpdatedAt:          m.UpdatedAt.Time,
	}
}

// CreateModeloEnvaseRequest is the payload for POST /modelos-envase.
// TieneVariantes defaults to true when omitted (nil) — most modelos vary by
// grosor. Set it to false only for a self-contained modelo like "envase de
// lujo" that never varies.
type CreateModeloEnvaseRequest struct {
	Tipo               string `json:"tipo" validate:"required,min=2,max=100"`
	TamanoOz           string `json:"tamano_oz" validate:"required,numeric"`
	EquivGramos        string `json:"equiv_gramos" validate:"required,numeric"`
	PrecioSolo         string `json:"precio_solo" validate:"required,numeric"`
	PrecioConFragancia string `json:"precio_con_fragancia" validate:"required,numeric"`
	PrecioRecarga      string `json:"precio_recarga" validate:"required,numeric"`
	TieneVariantes     *bool  `json:"tiene_variantes,omitempty"`
}

// UpdateModeloEnvaseRequest is the payload for PATCH /modelos-envase/:id. A
// nil field means "leave unchanged".
type UpdateModeloEnvaseRequest struct {
	Tipo               *string `json:"tipo,omitempty" validate:"omitempty,min=2,max=100"`
	TamanoOz           *string `json:"tamano_oz,omitempty" validate:"omitempty,numeric"`
	EquivGramos        *string `json:"equiv_gramos,omitempty" validate:"omitempty,numeric"`
	PrecioSolo         *string `json:"precio_solo,omitempty" validate:"omitempty,numeric"`
	PrecioConFragancia *string `json:"precio_con_fragancia,omitempty" validate:"omitempty,numeric"`
	PrecioRecarga      *string `json:"precio_recarga,omitempty" validate:"omitempty,numeric"`
}
