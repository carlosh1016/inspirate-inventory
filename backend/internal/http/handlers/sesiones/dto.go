package sesiones

import (
	"fmt"
	"time"

	domaincuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/cuadres"
	domainsesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/sesiones"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
)

// UsuarioBriefResponse is the (id, nombre_completo) of the vendedora a
// sesion belongs to.
type UsuarioBriefResponse struct {
	ID             int64  `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
}

func toUsuarioBriefResponse(u *domaincuadres.UsuarioBrief) *UsuarioBriefResponse {
	if u == nil {
		return nil
	}
	return &UsuarioBriefResponse{ID: u.ID, NombreCompleto: u.NombreCompleto}
}

// SesionResponse is the response shape for entrada/salida/list/update.
type SesionResponse struct {
	ID              int64                 `json:"id"`
	Usuario         *UsuarioBriefResponse `json:"usuario,omitempty"`
	EntradaAt       time.Time             `json:"entrada_at"`
	SalidaAt        *time.Time            `json:"salida_at"`
	HorasTrabajadas *string               `json:"horas_trabajadas"`
}

func toSesionResponse(s domainsesiones.Sesion) SesionResponse {
	return SesionResponse{
		ID:              s.ID,
		Usuario:         toUsuarioBriefResponse(s.Usuario),
		EntradaAt:       s.EntradaAt,
		SalidaAt:        s.SalidaAt,
		HorasTrabajadas: formatHorasTrabajadas(s.HorasTrabajadas),
	}
}

// formatHorasTrabajadas renders a duration as "HH:MM:SS", or "Nd HH:MM:SS"
// when it spans one or more full days (a forgotten clock-out) — nil while
// the sesion is still open.
func formatHorasTrabajadas(d *time.Duration) *string {
	if d == nil {
		return nil
	}
	total := *d
	days := int64(total / (24 * time.Hour))
	rem := total % (24 * time.Hour)
	hours := int64(rem / time.Hour)
	rem %= time.Hour
	minutes := int64(rem / time.Minute)
	rem %= time.Minute
	seconds := int64(rem / time.Second)

	var s string
	if days > 0 {
		s = fmt.Sprintf("%dd %02d:%02d:%02d", days, hours, minutes, seconds)
	} else {
		s = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return &s
}

// UpdateSesionRequest is the payload for PATCH /sesiones-laborales/:id.
type UpdateSesionRequest struct {
	EntradaAt *time.Time `json:"entrada_at,omitempty"`
	SalidaAt  *time.Time `json:"salida_at,omitempty"`
}

// ResumenSesionResponse is one vendedora's aggregated worked time.
type ResumenSesionResponse struct {
	Usuario        UsuarioBriefResponse `json:"usuario"`
	TotalHoras     string               `json:"total_horas"`
	SesionesCount  int64                `json:"sesiones_count"`
	DiasTrabajados int64                `json:"dias_trabajados"`
}

func toResumenSesionResponse(item usecase.ResumenItem) ResumenSesionResponse {
	total := formatHorasTrabajadas(&item.TotalHoras)
	return ResumenSesionResponse{
		Usuario:        UsuarioBriefResponse{ID: item.UsuarioID, NombreCompleto: item.NombreCompleto},
		TotalHoras:     *total,
		SesionesCount:  item.SesionesCount,
		DiasTrabajados: item.DiasTrabajados,
	}
}
