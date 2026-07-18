// Package handlers contiene los handlers HTTP organizados por dominio.
package handlers

import (
	"net/http"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

const apiVersion = "0.1.0"

// Health responde con el estado del servicio. Usado para health checks.
func Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"version":   apiVersion,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
