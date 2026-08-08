package fragancias

import (
	"time"

	commonrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// StockResponse is the vitrina/bodega/total stock snapshot for an item.
type StockResponse struct {
	Vitrina string `json:"vitrina"`
	Bodega  string `json:"bodega"`
	Total   string `json:"total"`
}

// FraganciaResponse never includes sede_id or deleted_at.
type FraganciaResponse struct {
	ID                int64         `json:"id"`
	NombreComercial   string        `json:"nombre_comercial"`
	NombreAlternativo *string       `json:"nombre_alternativo"`
	Genero            string        `json:"genero"`
	NumeroGenero      int32         `json:"numero_genero"`
	GramosMinimo      string        `json:"gramos_minimo"`
	Activo            bool          `json:"activo"`
	Stock             StockResponse `json:"stock"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

func toFraganciaResponseFromGet(f generated.GetFraganciaByIDRow) FraganciaResponse {
	return FraganciaResponse{
		ID:                f.ID,
		NombreComercial:   f.NombreComercial,
		NombreAlternativo: commonrepo.StringPtr(f.NombreAlternativo),
		Genero:            string(f.Genero),
		NumeroGenero:      f.NumeroGenero,
		GramosMinimo:      f.GramosMinimo.String(),
		Activo:            f.Activo,
		Stock: StockResponse{
			Vitrina: f.StockVitrina.String(),
			Bodega:  f.StockBodega.String(),
			Total:   f.StockVitrina.Add(f.StockBodega).String(),
		},
		CreatedAt: f.CreatedAt.Time,
		UpdatedAt: f.UpdatedAt.Time,
	}
}

func toFraganciaResponseFromList(f generated.ListFraganciasPaginatedRow) FraganciaResponse {
	return FraganciaResponse{
		ID:                f.ID,
		NombreComercial:   f.NombreComercial,
		NombreAlternativo: commonrepo.StringPtr(f.NombreAlternativo),
		Genero:            string(f.Genero),
		NumeroGenero:      f.NumeroGenero,
		GramosMinimo:      f.GramosMinimo.String(),
		Activo:            f.Activo,
		Stock: StockResponse{
			Vitrina: f.StockVitrina.String(),
			Bodega:  f.StockBodega.String(),
			Total:   f.StockVitrina.Add(f.StockBodega).String(),
		},
		CreatedAt: f.CreatedAt.Time,
		UpdatedAt: f.UpdatedAt.Time,
	}
}

// CreateFraganciaRequest is the payload for POST /fragancias.
type CreateFraganciaRequest struct {
	NombreComercial   string  `json:"nombre_comercial" validate:"required,min=2,max=200"`
	NombreAlternativo *string `json:"nombre_alternativo,omitempty" validate:"omitempty,max=200"`
	Genero            string  `json:"genero" validate:"required,oneof=masculina femenina"`
	NumeroGenero      int32   `json:"numero_genero" validate:"required,min=1"`
	GramosMinimo      string  `json:"gramos_minimo" validate:"required,numeric"`
}

// UpdateFraganciaRequest is the payload for PATCH /fragancias/:id. A nil
// field means "leave unchanged".
type UpdateFraganciaRequest struct {
	NombreComercial   *string `json:"nombre_comercial,omitempty" validate:"omitempty,min=2,max=200"`
	NombreAlternativo *string `json:"nombre_alternativo,omitempty" validate:"omitempty,max=200"`
	Genero            *string `json:"genero,omitempty" validate:"omitempty,oneof=masculina femenina"`
	NumeroGenero      *int32  `json:"numero_genero,omitempty" validate:"omitempty,min=1"`
	GramosMinimo      *string `json:"gramos_minimo,omitempty" validate:"omitempty,numeric"`
}

// SiguienteNumeroResponse is the payload for GET /fragancias/siguiente-numero.
type SiguienteNumeroResponse struct {
	Siguiente int32 `json:"siguiente"`
}
