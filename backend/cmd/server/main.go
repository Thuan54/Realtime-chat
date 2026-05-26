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
	"github.com/user/chat-app/internal/db"
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

	// 2. Database & Auto-Migrations (Blocking startup)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Error("startup failed: DATABASE_URL environment variable is required")
		os.Exit(1)
	}
	migrationDir := os.Getenv("MIGRATIONS_DIR")
	if migrationDir == "" {
		migrationDir = "./migrations"
	}

	ctx := context.Background()
	dbPool, err := db.RunMigrations(ctx, dbURL, migrationDir, log)
	if err != nil {
		log.Error("startup aborted: migration failure", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbPool.Close()

	// 3. Configure Gin & Middleware
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(server.LoggerMiddleware(log))

	// 4. Register Health & Ready (Wired to real DB pool)
	checker := &health.Checker{
		DBPing:    func(ctx context.Context) error { return dbPool.PingContext(ctx) },
		RedisPing: func(ctx context.Context) error { return nil }, // Stub until Redis client is wired
	}
	r.GET("/health", checker.Health)
	r.GET("/ready", checker.Ready)

	// 5. HTTP Server Setup
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 6. Graceful Shutdown
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

	log.Info("server started", slog.String("addr", addr), slog.String("db_pool", "connected"))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}

	<-idleConnsClosed
	log.Info("server stopped successfully")
}
