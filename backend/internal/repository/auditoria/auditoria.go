// Package auditoria is the persistence port for the auditoria table. Only
// Insert exists in this tanda — listing/filtering audit records is Tanda 6.
package auditoria

import "context"

// Entry is one row to record in `auditoria`. UsuarioID is nil for events
// like a failed login with an unknown correo.
type Entry struct {
	UsuarioID     *int64
	Accion        string
	TablaAfectada *string
	RegistroID    *int64
	DatosAntes    []byte
	DatosDespues  []byte
	IP            string
	UserAgent     string
}

// Repository is the persistence port for auditoria, consumed by usecases.
type Repository interface {
	Insert(ctx context.Context, e Entry) error
}
