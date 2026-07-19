// Package logger provides structured logging via log/slog: JSON in
// production, human-readable text in development.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a *slog.Logger. env selects the handler ("production" -> JSON,
// anything else -> text). level accepts "debug", "info", "warn", "error"
// (default: "info").
func New(env, level string) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}

	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
