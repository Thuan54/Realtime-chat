package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// New initializes a PostgreSQL connection pool via pgx driver,
// applies pooling limits, and retries on startup.
func New(ctx context.Context, log *slog.Logger, dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	var pingErr error
	for attempt := 1; attempt <= 5; attempt++ {
		pingErr = db.PingContext(ctx)
		if pingErr == nil {
			log.Info("database connected successfully")
			return db, nil
		}
		log.Warn("database ping failed, retrying...",
			slog.Int("attempt", attempt),
			slog.String("error", pingErr.Error()))
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}

	return nil, fmt.Errorf("failed to connect to database after retries: %w", pingErr)
}
