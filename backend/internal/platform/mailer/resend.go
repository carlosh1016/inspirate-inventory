package mailer

import "context"

// ResendMailer implementa Mailer usando la API de Resend.
type ResendMailer struct {
	apiKey string
	from   string
}

// NewResendMailer crea un Mailer respaldado por Resend.
func NewResendMailer(apiKey, from string) *ResendMailer {
	return &ResendMailer{apiKey: apiKey, from: from}
}

// Send envía un correo a través de la API de Resend.
// TODO(module-1+): implementar la llamada HTTP real a Resend.
func (m *ResendMailer) Send(ctx context.Context, msg Message) error {
	return nil
}
