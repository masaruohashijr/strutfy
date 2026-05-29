package database

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"log"
	"github.com/uptrace/bun"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(ctx context.Context, db *bun.DB) error {
    log.Println("Running migrations...")
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`); err != nil {
		return err
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)

	for _, file := range files {
	    log.Printf("Applying migration: %s", file)
		applied, err := migrationApplied(ctx, db, file)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		sqlBytes, err := migrationFiles.ReadFile("migrations/" + file)
		if err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}

		if _, err := db.ExecContext(ctx,
			`INSERT INTO schema_migrations (version) VALUES (?)`,
			file,
		); err != nil {
			return err
		}
	}

	return nil
}

func migrationApplied(ctx context.Context, db *bun.DB, version string) (bool, error) {
	var exists bool

	err := db.NewSelect().
		Table("schema_migrations").
		ColumnExpr("COUNT(*) > 0").
		Where("version = ?", version).
		Scan(ctx, &exists)

	return exists, err
}