// Package ratelimit provides per-key request throttling for sensitive,
// unauthenticated endpoints (login, password reset request).
package ratelimit

import (
	"context"
	"time"
)

// Limiter checks whether a request identified by key is allowed under a
// rate limit policy.
type Limiter interface {
	// Allow reports whether the request is allowed. When it isn't,
	// retryAfter is how long the caller should wait before trying again.
	Allow(ctx context.Context, key string) (allowed bool, retryAfter time.Duration, err error)
}
