package middleware

import (
	"net"
	"net/http"
	"strings"
)

// AuditContext extracts the caller's real IP (honoring X-Forwarded-For, set
// by a trusted reverse proxy) and User-Agent into the request context so
// usecases can record them in `auditoria` without depending on net/http.
func AuditContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := withIP(r.Context(), realIP(r))
		ctx = withUserAgent(ctx, r.UserAgent())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func realIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if ip := strings.TrimSpace(strings.Split(fwd, ",")[0]); ip != "" {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
