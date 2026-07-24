// Package logger wraps slog with a small project-friendly setup.
package logger

import (
	"log/slog"
	"os"
)

// New returns a structured logger. In production it emits JSON; in development
// it emits human-readable text.
func New(env string) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}
