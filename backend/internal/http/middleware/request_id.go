package middleware

import (
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// RequestID assigns a unique ID to each request, propagated via context and
// the X-Request-Id response header. Delegates to chi's implementation so
// every middleware in this package is imported from one place.
func RequestID(next http.Handler) http.Handler {
	return chimiddleware.RequestID(next)
}
