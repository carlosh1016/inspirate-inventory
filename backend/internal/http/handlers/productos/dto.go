package productos

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

// ProductoResponse never includes sede_id or deleted_at.
type ProductoResponse struct {
	ID          int64         `json:"id"`
	Nombre      string        `json:"nombre"`
	Categoria   string        `json:"categoria"`
	Precio      string        `json:"precio"`
	StockMinimo int32         `json:"stock_minimo"`
	Activo      bool          `json:"activo"`
	Stock       StockResponse `json:"stock"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func toProductoResponseFromGet(p generated.GetProductoByIDRow) ProductoResponse {
	return ProductoResponse{
		ID:          p.ID,
		Nombre:      p.Nombre,
		Categoria:   string(p.Categoria),
		Precio:      p.Precio.String(),
		StockMinimo: p.StockMinimo,
		Activo:      p.Activo,
		Stock: StockResponse{
			Vitrina: p.StockVitrina.String(),
			Bodega:  p.StockBodega.String(),
			Total:   p.StockVitrina.Add(p.StockBodega).String(),
		},
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}
}

func toProductoResponseFromList(p generated.ListProductosPaginatedRow) ProductoResponse {
	return ProductoResponse{
		ID:          p.ID,
		Nombre:      p.Nombre,
		Categoria:   string(p.Categoria),
		Precio:      p.Precio.String(),
		StockMinimo: p.StockMinimo,
		Activo:      p.Activo,
		Stock: StockResponse{
			Vitrina: p.StockVitrina.String(),
			Bodega:  p.StockBodega.String(),
			Total:   p.StockVitrina.Add(p.StockBodega).String(),
		},
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}
}

// CreateProductoRequest is the payload for POST /productos.
type CreateProductoRequest struct {
	Nombre      string `json:"nombre" validate:"required,min=2,max=200"`
	Categoria   string `json:"categoria" validate:"required,oneof=crema feromona hogar regalo otro"`
	Precio      string `json:"precio" validate:"required,numeric"`
	StockMinimo *int32 `json:"stock_minimo,omitempty" validate:"omitempty,min=0"`
}

// UpdateProductoRequest is the payload for PATCH /productos/:id. A nil
// field means "leave unchanged". A vendedora may only set stock_minimo —
// enforced by the usecase, not this DTO, since the wire shape is the same
// for both roles.
type UpdateProductoRequest struct {
	Nombre      *string `json:"nombre,omitempty" validate:"omitempty,min=2,max=200"`
	Categoria   *string `json:"categoria,omitempty" validate:"omitempty,oneof=crema feromona hogar regalo otro"`
	Precio      *string `json:"precio,omitempty" validate:"omitempty,numeric"`
	StockMinimo *int32  `json:"stock_minimo,omitempty" validate:"omitempty,min=0"`
}
