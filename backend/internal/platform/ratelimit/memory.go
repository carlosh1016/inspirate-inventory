package ratelimit

import (
	"context"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

const (
	loginLimit          = 5
	loginPeriod         = time.Minute
	passwordResetLimit  = 3
	passwordResetPeriod = time.Hour
)

// memoryLimiter adapts ulule/limiter's in-process store to our Limiter
// interface. State resets on process restart, which is fine for this
// single-instance deployment.
type memoryLimiter struct {
	l *limiter.Limiter
}

func newMemoryLimiter(period time.Duration, limit int64) *memoryLimiter {
	return &memoryLimiter{
		l: limiter.New(memory.NewStore(), limiter.Rate{Period: period, Limit: limit}),
	}
}

// Allow increments the counter for key and reports whether the request is
// still within the limit.
func (m *memoryLimiter) Allow(ctx context.Context, key string) (bool, time.Duration, error) {
	res, err := m.l.Get(ctx, key)
	if err != nil {
		return false, 0, err
	}

	if !res.Reached {
		return true, 0, nil
	}

	retryAfter := time.Until(time.Unix(res.Reset, 0))
	if retryAfter < 0 {
		retryAfter = 0
	}
	return false, retryAfter, nil
}

// NewLoginLimiter allows 5 requests per minute per key.
func NewLoginLimiter() Limiter {
	return newMemoryLimiter(loginPeriod, loginLimit)
}

// NewPasswordResetLimiter allows 3 requests per hour per key.
func NewPasswordResetLimiter() Limiter {
	return newMemoryLimiter(passwordResetPeriod, passwordResetLimit)
}
