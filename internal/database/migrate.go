package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// RunSQLMigrations applies *.up.sql files from dir once (tracked in schema_migrations).
func RunSQLMigrations(db *gorm.DB, dir string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if _, err := sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`); err != nil {
		return fmt.Errorf("schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir %q: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".up.sql") {
			files = append(files, name)
		}
	}
	sort.Strings(files)

	for _, name := range files {
		version := strings.TrimSuffix(name, ".up.sql")
		var count int
		if err := sqlDB.QueryRow(`SELECT COUNT(1) FROM schema_migrations WHERE version = $1`, version).Scan(&count); err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}

		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(body)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", name, err)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
