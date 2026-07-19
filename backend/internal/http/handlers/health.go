// Package handlers contains the HTTP handlers organized by domain.
package handlers

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/db"
)

const apiVersion = "0.1.0"

// Health reports service status, including DB connectivity. Returns 503
// with status "degraded" when the database is unreachable.
func Health(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		dbCheck := "ok"
		httpStatus := http.StatusOK

		if err := db.HealthCheck(r.Context(), pool); err != nil {
			status = "degraded"
			dbCheck = "error"
			httpStatus = http.StatusServiceUnavailable
		}

		response.WriteJSON(w, httpStatus, map[string]any{
			"status":    status,
			"version":   apiVersion,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"checks": map[string]string{
				"db": dbCheck,
			},
		})
	}
}
