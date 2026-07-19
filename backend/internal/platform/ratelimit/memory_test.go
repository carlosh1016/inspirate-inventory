package ratelimit_test

import (
	"context"
	"testing"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/ratelimit"
)

func TestLoginLimiterBlocksAfterFiveRequests(t *testing.T) {
	l := ratelimit.NewLoginLimiter()
	ctx := context.Background()
	key := "1.2.3.4:test@inspirate.co"

	for i := 0; i < 5; i++ {
		allowed, _, err := l.Allow(ctx, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}

	allowed, retryAfter, err := l.Allow(ctx, key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected 6th request to be blocked")
	}
	if retryAfter <= 0 {
		t.Fatalf("expected positive retryAfter, got %v", retryAfter)
	}
}

func TestLoginLimiterKeysAreIndependent(t *testing.T) {
	l := ratelimit.NewLoginLimiter()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		if _, _, err := l.Allow(ctx, "key-a"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	allowed, _, err := l.Allow(ctx, "key-b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected a different key to have its own limit")
	}
}

func TestPasswordResetLimiterAllowsThreeRequests(t *testing.T) {
	l := ratelimit.NewPasswordResetLimiter()
	ctx := context.Background()
	key := "1.2.3.4:ana@inspirate.co"

	for i := 0; i < 3; i++ {
		allowed, _, err := l.Allow(ctx, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}

	allowed, _, err := l.Allow(ctx, key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected 4th request to be blocked")
	}
}
