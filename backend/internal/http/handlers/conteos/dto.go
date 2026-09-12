package conteos

import (
	"time"

	domainconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/conteos"
)

// UsuarioBriefResponse is the (id, nombre_completo) of who did something.
type UsuarioBriefResponse struct {
	ID             int64  `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
}

func toUsuarioBriefResponse(u *domainconteos.UsuarioBrief) *UsuarioBriefResponse {
	if u == nil {
		return nil
	}
	return &UsuarioBriefResponse{ID: u.ID, NombreCompleto: u.NombreCompleto}
}

// ConteoItemResponse is one fragancia's row within a conteo. Diferencia,
// ExcedeLimite and PorcentajeVendido are computed server-side (see
// internal/domain/conteos.ConteoItem) so the frontend never reimplements
// the 2g business rule.
type ConteoItemResponse struct {
	ID                int64   `json:"id"`
	FraganciaID       int64   `json:"fragancia_id"`
	FraganciaNombre   string  `json:"fragancia_nombre"`
	SaldoInicial      string  `json:"saldo_inicial"`
	GramosSistema     string  `json:"gramos_sistema"`
	GramosFisico      *string `json:"gramos_fisico"`
	Diferencia        *string `json:"diferencia"`
	ExcedeLimite      bool    `json:"excede_limite"`
	PorcentajeVendido *string `json:"porcentaje_vendido"`
}

func toConteoItemResponse(i domainconteos.ConteoItem) ConteoItemResponse {
	var gramosFisico *string
	if i.GramosFisico != nil {
		s := i.GramosFisico.String()
		gramosFisico = &s
	}
	var diferencia *string
	if d := i.Diferencia(); d != nil {
		s := d.String()
		diferencia = &s
	}
	var porcentaje *string
	if p := i.PorcentajeVendido(); p != nil {
		s := p.StringFixed(2)
		porcentaje = &s
	}
	return ConteoItemResponse{
		ID:                i.ID,
		FraganciaID:       i.FraganciaID,
		FraganciaNombre:   i.FraganciaNombre,
		SaldoInicial:      i.SaldoInicial.String(),
		GramosSistema:     i.GramosSistema.String(),
		GramosFisico:      gramosFisico,
		Diferencia:        diferencia,
		ExcedeLimite:      i.ExcedeLimite(),
		PorcentajeVendido: porcentaje,
	}
}

// ConteoResponse is the full response shape for GET/POST/PATCH/cerrar.
type ConteoResponse struct {
	ID         int64                 `json:"id"`
	SedeID     int64                 `json:"sede_id"`
	Periodo    string                `json:"periodo"`
	Estado     string                `json:"estado"`
	CreadoPor  UsuarioBriefResponse  `json:"creado_por"`
	CerradoPor *UsuarioBriefResponse `json:"cerrado_por"`
	CerradoAt  *time.Time            `json:"cerrado_at"`
	CreatedAt  time.Time             `json:"created_at"`
	Items      []ConteoItemResponse  `json:"items"`
}

func toConteoResponse(c domainconteos.ConteoInventario) ConteoResponse {
	items := make([]ConteoItemResponse, len(c.Items))
	for i, it := range c.Items {
		items[i] = toConteoItemResponse(it)
	}
	return ConteoResponse{
		ID:         c.ID,
		SedeID:     c.SedeID,
		Periodo:    c.Periodo.Format("2006-01-02"),
		Estado:     string(c.Estado),
		CreadoPor:  *toUsuarioBriefResponse(&c.CreadoPor),
		CerradoPor: toUsuarioBriefResponse(c.CerradoPor),
		CerradoAt:  c.CerradoAt,
		CreatedAt:  c.CreatedAt,
		Items:      items,
	}
}

// ConteoListItemResponse is one row of GET /conteos-inventario — no items,
// same convention as cuadres' list response.
type ConteoListItemResponse struct {
	ID         int64                 `json:"id"`
	Periodo    string                `json:"periodo"`
	Estado     string                `json:"estado"`
	CreadoPor  UsuarioBriefResponse  `json:"creado_por"`
	CerradoPor *UsuarioBriefResponse `json:"cerrado_por"`
	CreatedAt  time.Time             `json:"created_at"`
}

func toConteoListItemResponse(c domainconteos.ConteoInventario) ConteoListItemResponse {
	return ConteoListItemResponse{
		ID:         c.ID,
		Periodo:    c.Periodo.Format("2006-01-02"),
		Estado:     string(c.Estado),
		CreadoPor:  *toUsuarioBriefResponse(&c.CreadoPor),
		CerradoPor: toUsuarioBriefResponse(c.CerradoPor),
		CreatedAt:  c.CreatedAt,
	}
}

// CrearConteoRequest is the payload for POST /conteos-inventario. Periodo is
// optional — omitted, it defaults to the current month.
type CrearConteoRequest struct {
	Periodo *string `json:"periodo,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// ActualizarItemRequest is the payload for PATCH
// /conteos-inventario/:id/items/:item_id.
type ActualizarItemRequest struct {
	GramosFisico string `json:"gramos_fisico" validate:"required,numeric"`
}
