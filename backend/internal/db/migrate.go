package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// RunMigrations opens a PostgreSQL connection, applies pending SQL migrations from the specified directory,
// and returns the initialized *sql.DB pool for application use.
func RunMigrations(ctx context.Context, dbURL, migrationDir string, log *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connectivity before attempting migrations
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	// Suppress default goose logging to prevent duplicate output; we handle it via structured slog
	goose.SetLogger(goose.NopLogger())

	log.Info("applying database migrations", slog.String("directory", migrationDir))
	if err := goose.UpContext(ctx, db, migrationDir); err != nil {
		db.Close()
		return nil, fmt.Errorf("migration execution failed: %w", err)
	}

	log.Info("database migrations applied successfully")
	return db, nil
}
