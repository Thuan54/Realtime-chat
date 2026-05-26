package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/chat-app/internal/logger"
)

func TestNew_DoesNotPanic(t *testing.T) {
	// Verify logger initialization with valid levels
	assert.NotPanics(t, func() { logger.New("debug", "test") })
	assert.NotPanics(t, func() { logger.New("info", "test") })
	assert.NotPanics(t, func() { logger.New("warn", "test") })
	assert.NotPanics(t, func() { logger.New("error", "test") })
	assert.NotPanics(t, func() { logger.New("invalid", "test") }) // defaults to info
}

func TestWithContext_ExtractsMetadata(t *testing.T) {
	baseLog := slog.Default()
	ctx := context.Background()

	// Without keys, should return same logger
	enriched := logger.WithContext(ctx, baseLog)
	assert.Equal(t, baseLog, enriched, "logger should remain unchanged without context keys")

	// With request_id only
	ctx = context.WithValue(ctx, logger.RequestIDKey, "req-abc-123")
	enriched = logger.WithContext(ctx, baseLog)
	assert.NotEqual(t, baseLog, enriched, "logger should be enriched with request_id")

	// With user_id only
	ctx = context.WithValue(ctx, logger.UserIDKey, "user-456")
	enriched = logger.WithContext(ctx, baseLog)
	assert.NotEqual(t, baseLog, enriched, "logger should be enriched with user_id")

	// With both
	ctx = context.Background()
	ctx = context.WithValue(ctx, logger.RequestIDKey, "req-abc-123")
	ctx = context.WithValue(ctx, logger.UserIDKey, "user-456")
	enriched = logger.WithContext(ctx, baseLog)
	assert.NotEqual(t, baseLog, enriched, "logger should be enriched with both keys")
}
