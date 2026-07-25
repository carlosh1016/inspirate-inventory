package auditoria

import (
	"encoding/json"
	"time"

	domainauditoria "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auditoria"
)

// UsuarioBriefResponse is the actor identity embedded in an evento.
type UsuarioBriefResponse struct {
	ID             int64  `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
}

// EventoResponse is the JSON shape of one audit evento. datos_antes and
// datos_despues are passed through as raw JSON (objects, not strings); a NULL
// payload serializes to null.
type EventoResponse struct {
	ID            int64                 `json:"id"`
	Usuario       *UsuarioBriefResponse `json:"usuario"`
	Accion        string                `json:"accion"`
	TablaAfectada *string               `json:"tabla_afectada"`
	RegistroID    *int64                `json:"registro_id"`
	DatosAntes    json.RawMessage       `json:"datos_antes"`
	DatosDespues  json.RawMessage       `json:"datos_despues"`
	IP            *string               `json:"ip"`
	UserAgent     *string               `json:"user_agent"`
	CreatedAt     time.Time             `json:"created_at"`
}

func toEventoResponse(ev domainauditoria.Evento) EventoResponse {
	resp := EventoResponse{
		ID:            ev.ID,
		Accion:        ev.Accion,
		TablaAfectada: ev.TablaAfectada,
		RegistroID:    ev.RegistroID,
		DatosAntes:    ev.DatosAntes,
		DatosDespues:  ev.DatosDespues,
		IP:            ev.IP,
		UserAgent:     ev.UserAgent,
		CreatedAt:     ev.CreatedAt,
	}
	if ev.Usuario != nil {
		resp.Usuario = &UsuarioBriefResponse{ID: ev.Usuario.ID, NombreCompleto: ev.Usuario.NombreCompleto}
	}
	return resp
}

func toEventoResponses(eventos []domainauditoria.Evento) []EventoResponse {
	out := make([]EventoResponse, 0, len(eventos))
	for _, ev := range eventos {
		out = append(out, toEventoResponse(ev))
	}
	return out
}
