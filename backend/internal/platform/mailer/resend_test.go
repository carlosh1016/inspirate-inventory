package mailer

import (
	"bytes"
	"strings"
	"testing"
)

func TestPasswordResetTemplateRenders(t *testing.T) {
	var body bytes.Buffer
	err := passwordResetTemplate.Execute(&body, map[string]string{
		"Nombre": "Ana Pérez",
		"URL":    "https://app.inspirate.co/reset-password?token=abc123",
	})
	if err != nil {
		t.Fatalf("unexpected error executing template: %v", err)
	}

	rendered := body.String()
	if !strings.Contains(rendered, "Ana Pérez") {
		t.Errorf("expected rendered template to contain the name, got: %s", rendered)
	}
	if !strings.Contains(rendered, "https://app.inspirate.co/reset-password?token=abc123") {
		t.Errorf("expected rendered template to contain the reset URL, got: %s", rendered)
	}
}
