package mailer_test

import (
	"context"
	"testing"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/mailer"
)

func TestMockMailerRecordsSentEmails(t *testing.T) {
	m := mailer.NewMock()

	err := m.SendPasswordReset(context.Background(), "ana@inspirate.co", "https://app.inspirate.co/reset?token=abc", "Ana Pérez")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sent := m.Sent()
	if len(sent) != 1 {
		t.Fatalf("expected 1 sent email, got %d", len(sent))
	}
	if sent[0].To != "ana@inspirate.co" {
		t.Errorf("unexpected To: %q", sent[0].To)
	}
	if sent[0].ResetURL != "https://app.inspirate.co/reset?token=abc" {
		t.Errorf("unexpected ResetURL: %q", sent[0].ResetURL)
	}
	if sent[0].NombreCompleto != "Ana Pérez" {
		t.Errorf("unexpected NombreCompleto: %q", sent[0].NombreCompleto)
	}
}

var _ mailer.Mailer = (*mailer.MockMailer)(nil)
var _ mailer.Mailer = (*mailer.ResendMailer)(nil)
