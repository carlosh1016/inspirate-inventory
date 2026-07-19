// Package middleware holds the HTTP middleware chain (auth, roles, audit
// context, request id) shared by every route group.
package middleware

import "context"

type ctxKey int

const (
	userCtxKey ctxKey = iota
	ipCtxKey
	userAgentCtxKey
)

// AuthenticatedUser is the identity extracted from a valid access token by
// Auth and stored in the request context for handlers and RequireRole.
type AuthenticatedUser struct {
	ID     int64
	Rol    string
	SedeID int64
}

// WithUser stores the authenticated user in ctx.
func WithUser(ctx context.Context, u AuthenticatedUser) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

// UserFromContext returns the authenticated user placed by Auth. ok is
// false for unauthenticated requests.
func UserFromContext(ctx context.Context) (AuthenticatedUser, bool) {
	u, ok := ctx.Value(userCtxKey).(AuthenticatedUser)
	return u, ok
}

func withIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ipCtxKey, ip)
}

// IPFromContext returns the caller IP recorded by AuditContext.
func IPFromContext(ctx context.Context) string {
	ip, _ := ctx.Value(ipCtxKey).(string)
	return ip
}

func withUserAgent(ctx context.Context, ua string) context.Context {
	return context.WithValue(ctx, userAgentCtxKey, ua)
}

// UserAgentFromContext returns the caller's User-Agent recorded by AuditContext.
func UserAgentFromContext(ctx context.Context) string {
	ua, _ := ctx.Value(userAgentCtxKey).(string)
	return ua
}
