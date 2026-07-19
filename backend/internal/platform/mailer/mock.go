package mailer

import (
	"context"
	"sync"
)

// SentPasswordReset records one SendPasswordReset call captured by MockMailer.
type SentPasswordReset struct {
	To             string
	ResetURL       string
	NombreCompleto string
}

// MockMailer is an in-memory Mailer for tests and non-production
// environments without a Resend API key: it records every call instead of
// sending real email.
type MockMailer struct {
	mu   sync.Mutex
	sent []SentPasswordReset
}

// NewMock builds a MockMailer.
func NewMock() *MockMailer {
	return &MockMailer{}
}

// SendPasswordReset records the call and always succeeds.
func (m *MockMailer) SendPasswordReset(_ context.Context, to, resetURL, nombreCompleto string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sent = append(m.sent, SentPasswordReset{To: to, ResetURL: resetURL, NombreCompleto: nombreCompleto})
	return nil
}

// Sent returns a copy of every SendPasswordReset call recorded so far.
func (m *MockMailer) Sent() []SentPasswordReset {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]SentPasswordReset, len(m.sent))
	copy(out, m.sent)
	return out
}
