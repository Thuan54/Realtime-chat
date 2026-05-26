package logger

import (
	"context"
	"log/slog"
	"os"
)

// Context keys for propagating metadata through request lifecycle
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
)

// New initializes a structured JSON logger with a global component tag
func New(level, component string) *slog.Logger {
	lvl := slog.LevelInfo
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	return slog.New(handler).With(slog.String("component", component))
}

// WithContext enriches the logger with request_id and user_id if present in ctx
func WithContext(ctx context.Context, log *slog.Logger) *slog.Logger {
	l := log
	if id, ok := ctx.Value(RequestIDKey).(string); ok && id != "" {
		l = l.With(slog.String("request_id", id))
	}
	if uid, ok := ctx.Value(UserIDKey).(string); ok && uid != "" {
		l = l.With(slog.String("user_id", uid))
	}
	return l
}
