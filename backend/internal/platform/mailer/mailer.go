// Package mailer sends transactional email (currently just password reset).
package mailer

import "context"

// Mailer sends the password-reset email to a user.
type Mailer interface {
	SendPasswordReset(ctx context.Context, to, resetURL, nombreCompleto string) error
}
