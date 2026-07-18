// Package logger provee un logger estructurado basado en log/slog.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New crea un *slog.Logger que escribe JSON a stdout con el nivel indicado.
// level acepta: "debug", "info", "warn", "error" (default: "info").
func New(level string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
	})
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
