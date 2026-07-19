package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"

	"github.com/resend/resend-go/v2"
)

//go:embed templates/password_reset.html
var templatesFS embed.FS

var passwordResetTemplate = template.Must(template.ParseFS(templatesFS, "templates/password_reset.html"))

// ResendMailer sends email through the Resend API.
type ResendMailer struct {
	client *resend.Client
	from   string
}

// NewResendMailer builds a Mailer backed by Resend.
func NewResendMailer(apiKey, from string) *ResendMailer {
	return &ResendMailer{client: resend.NewClient(apiKey), from: from}
}

// SendPasswordReset renders templates/password_reset.html and sends it via Resend.
func (m *ResendMailer) SendPasswordReset(ctx context.Context, to, resetURL, nombreCompleto string) error {
	var body bytes.Buffer
	if err := passwordResetTemplate.Execute(&body, map[string]string{
		"Nombre": nombreCompleto,
		"URL":    resetURL,
	}); err != nil {
		return fmt.Errorf("rendering password reset template: %w", err)
	}

	_, err := m.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{to},
		Subject: "Restablecer tu contraseña — Inspírate Inventory",
		Html:    body.String(),
	})
	if err != nil {
		return fmt.Errorf("sending password reset email via resend: %w", err)
	}

	return nil
}
