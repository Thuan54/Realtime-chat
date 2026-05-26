package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/chat-app/internal/logger"
)

// LoggerMiddleware creates a Gin handler that injects request IDs and logs structured JSON
func LoggerMiddleware(baseLog *slog.Logger) gin.HandlerFunc {
	if baseLog == nil {
		baseLog = slog.Default()
	}
	return func(c *gin.Context) {
		// Preserve or generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)

		// Attach to request context
		ctx := context.WithValue(c.Request.Context(), logger.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		reqLog := logger.WithContext(ctx, baseLog)
		start := time.Now()

		c.Next()

		// Log after handler execution
		latency := time.Since(start)
		reqLog.Info("http_request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency_ms", latency),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("client_ip", c.ClientIP()),
		)
	}
}
