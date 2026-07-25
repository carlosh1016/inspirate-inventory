// Package auditoria holds the pure domain types for reading the audit log. The
// write side lives in repository/auditoria (used across many usecases since
// Tanda 1); this package models the read side (list/get) added in Tanda 6.
package auditoria

import (
	"encoding/json"
	"time"
)

// UsuarioBrief is the minimal actor identity attached to an evento.
type UsuarioBrief struct {
	ID             int64
	NombreCompleto string
}

// Evento is one audit-log record. UsuarioID/Usuario are nil for actor-less
// events (e.g. a failed login with an unknown correo). DatosAntes/DatosDespues
// are the raw JSONB payloads, passed through untouched.
type Evento struct {
	ID            int64
	UsuarioID     *int64
	Accion        string
	TablaAfectada *string
	RegistroID    *int64
	DatosAntes    json.RawMessage
	DatosDespues  json.RawMessage
	IP            *string
	UserAgent     *string
	CreatedAt     time.Time
	Usuario       *UsuarioBrief
}
