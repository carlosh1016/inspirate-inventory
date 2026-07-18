// Package mailer define la interfaz de envío de correos y su implementación
// para Resend, además de un mock para tests.
package mailer

import "context"

// Message representa un correo a enviar.
type Message struct {
	To      string
	Subject string
	HTML    string
}

// Mailer envía correos electrónicos.
type Mailer interface {
	Send(ctx context.Context, msg Message) error
}
