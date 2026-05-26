package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/chat-app/internal/health"
	"github.com/user/chat-app/internal/logger"
	"github.com/user/chat-app/internal/server"
)

func main() {
	// 1. Initialize structured logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	log := logger.New(logLevel, "http-server")

	// 2. Configure Gin
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(server.LoggerMiddleware(log))

	// 3. Register Health & Ready endpoints
	// TODO: Wire actual pgxPool.Ping and redisClient.Ping here
	checker := &health.Checker{
		DBPing:    func(ctx context.Context) error { return nil }, // Stub for bootstrap
		RedisPing: func(ctx context.Context) error { return nil }, // Stub for bootstrap
	}
	r.GET("/health", checker.Health)
	r.GET("/ready", checker.Ready)

	// 4. HTTP Server Setup
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 5. Graceful Shutdown (matches architecture doc SIGTERM drain requirement)
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Info("shutting down server gracefully", slog.String("signal", "SIGTERM"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server forced to shutdown", slog.String("error", err.Error()))
		}
		close(idleConnsClosed)
	}()

	log.Info("starting server", slog.String("addr", addr))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}

	<-idleConnsClosed
	log.Info("server stopped successfully")
}
