package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Migration struct {
	Version string
	Name    string
	SQL     string
}

func RunMigrations(ctx context.Context, db *sql.DB) error {
	// Create migrations table if not exists
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	// Read all migration files
	migs, err := readMigrations()
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	// Execute each migration
	for _, mig := range migs {
		var exists bool
		err := db.QueryRowContext(ctx, `
			SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)
		`, mig.Version).Scan(&exists)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", mig.Version, err)
		}

		if exists {
			fmt.Printf("✓ Migration %s already applied\n", mig.Version)
			continue
		}

		fmt.Printf("Running migration %s...\n", mig.Version)
		_, err = db.ExecContext(ctx, mig.SQL)
		if err != nil {
			return fmt.Errorf("execute migration %s: %w", mig.Version, err)
		}

		_, err = db.ExecContext(ctx, `
			INSERT INTO schema_migrations (version) VALUES ($1)
		`, mig.Version)
		if err != nil {
			return fmt.Errorf("record migration %s: %w", mig.Version, err)
		}

		fmt.Printf("✓ Migration %s applied\n", mig.Version)
	}

	return nil
}

func readMigrations() ([]Migration, error) {
	var migs []Migration
	files, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return nil, err
	}

	sort.Strings(files)

	for _, file := range files {
		content, err := fs.ReadFile(migrations, file)
		if err != nil {
			return nil, err
		}

		version := strings.TrimSuffix(filepath.Base(file), ".sql")
		parts := strings.SplitN(version, "_", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid migration filename: %s", file)
		}

		migs = append(migs, Migration{
			Version: parts[0],
			Name:    parts[1],
			SQL:     string(content),
		})
	}

	return migs, nil
}
