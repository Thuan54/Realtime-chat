package db

import (
	"context"
	"database/sql"
	"log/slog"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// RunMigrations applies pending SQL migrations from the specified directory
func RunMigrations(ctx context.Context, db *sql.DB, migrationDir string, log *slog.Logger) error {
	absDir, err := filepath.Abs(migrationDir)
	if err != nil {
		return err
	}

	goose.SetDialect("postgres")
	// Suppress default goose logging to prevent duplicate output
	// we handle it via structured slog
	goose.SetLogger(goose.NopLogger())

	log.Info("applying database migrations", slog.String("dir", absDir))
	if err := goose.UpContext(ctx, db, migrationDir); err != nil {
		log.Error("migration failed", slog.String("error", err.Error()))
		return err
	}

	log.Info("database migrations applied successfully")
	return nil
}
