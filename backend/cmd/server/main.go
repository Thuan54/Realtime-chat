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
	"github.com/user/chat-app/internal/cache"
	"github.com/user/chat-app/internal/db"
	"github.com/user/chat-app/internal/health"
	"github.com/user/chat-app/internal/logger"
	"github.com/user/chat-app/internal/server"
)

func main() {
	log := logger.New(os.Getenv("LOG_LEVEL"), "http-server")
	dbDSN := os.Getenv("DATABASE_URL")
	redisURL := os.Getenv("REDIS_URL")
	migrationDir := os.Getenv("MIGRATION_DIR")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize data stores
	database, err := db.New(ctx, log, dbDSN)
	if err != nil {
		log.Error("database init failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer database.Close()

	redisClient, err := cache.New(ctx, log, redisURL)
	if err != nil {
		log.Error("redis init failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer redisClient.Close()

	// Auto-apply migrations before accepting traffic
	if err := db.RunMigrations(ctx, database, migrationDir, log); err != nil {
		log.Error("migration failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Router setup
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(server.LoggerMiddleware(log))

	// Health/Ready endpoints with real dependency pingers
	checker := &health.Checker{
		DBPing:    database.PingContext,
		RedisPing: func(c context.Context) error { return redisClient.Ping(c).Err() },
	}
	r.GET("/health", checker.Health)
	r.GET("/ready", checker.Ready)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	srv := &http.Server{Addr: addr, Handler: r}

	// Shutdown gracefully
	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info("shutting down gracefully")
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutCancel()
		if err := srv.Shutdown(shutCtx); err != nil {
			log.Error("server shutdown error", slog.String("error", err.Error()))
		}
		close(idleConnsClosed)
	}()

	// Start server
	log.Info("starting server", slog.String("addr", addr))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	<-idleConnsClosed
	log.Info("server stopped successfully")
}
