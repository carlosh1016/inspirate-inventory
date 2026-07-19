package variantesenvase

import (
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// StockResponse is the vitrina/bodega/total stock snapshot for an item.
type StockResponse struct {
	Vitrina string `json:"vitrina"`
	Bodega  string `json:"bodega"`
	Total   string `json:"total"`
}

// VarianteEnvaseResponse never includes sede_id or deleted_at.
type VarianteEnvaseResponse struct {
	ID             int64         `json:"id"`
	ModeloEnvaseID int64         `json:"modelo_envase_id"`
	Color          string        `json:"color"`
	StockMinimo    int32         `json:"stock_minimo"`
	Activo         bool          `json:"activo"`
	Stock          StockResponse `json:"stock"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

func toVarianteEnvaseResponseFromGet(v generated.GetVarianteEnvaseByIDRow) VarianteEnvaseResponse {
	return VarianteEnvaseResponse{
		ID:             v.ID,
		ModeloEnvaseID: v.ModeloEnvaseID,
		Color:          v.Color,
		StockMinimo:    v.StockMinimo,
		Activo:         v.Activo,
		Stock: StockResponse{
			Vitrina: v.StockVitrina.String(),
			Bodega:  v.StockBodega.String(),
			Total:   v.StockVitrina.Add(v.StockBodega).String(),
		},
		CreatedAt: v.CreatedAt.Time,
		UpdatedAt: v.UpdatedAt.Time,
	}
}

func toVarianteEnvaseResponseFromList(v generated.ListVariantesEnvasePaginatedRow) VarianteEnvaseResponse {
	return VarianteEnvaseResponse{
		ID:             v.ID,
		ModeloEnvaseID: v.ModeloEnvaseID,
		Color:          v.Color,
		StockMinimo:    v.StockMinimo,
		Activo:         v.Activo,
		Stock: StockResponse{
			Vitrina: v.StockVitrina.String(),
			Bodega:  v.StockBodega.String(),
			Total:   v.StockVitrina.Add(v.StockBodega).String(),
		},
		CreatedAt: v.CreatedAt.Time,
		UpdatedAt: v.UpdatedAt.Time,
	}
}

// CreateVarianteEnvaseRequest is the payload for POST /variantes-envase.
type CreateVarianteEnvaseRequest struct {
	ModeloEnvaseID int64  `json:"modelo_envase_id" validate:"required,min=1"`
	Color          string `json:"color" validate:"required,min=1,max=100"`
	StockMinimo    *int32 `json:"stock_minimo,omitempty" validate:"omitempty,min=0"`
}

// UpdateVarianteEnvaseRequest is the payload for PATCH /variantes-envase/:id.
// A nil field means "leave unchanged". Both admin and vendedora may set
// either field — variantes_envase has no field this endpoint restricts by
// role.
type UpdateVarianteEnvaseRequest struct {
	Color       *string `json:"color,omitempty" validate:"omitempty,min=1,max=100"`
	StockMinimo *int32  `json:"stock_minimo,omitempty" validate:"omitempty,min=0"`
}
