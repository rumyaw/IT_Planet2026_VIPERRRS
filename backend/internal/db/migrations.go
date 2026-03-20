package db

import (
	"database/sql"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func ApplyMigrationsFromDSN(dsn string, migrationsDir string) error {
	db, err := OpenStdlib(dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	return ApplyMigrationsDB(db, migrationsDir)
}

// ApplyMigrationsDB runs goose migrations using provided *sql.DB.
func ApplyMigrationsDB(db *sql.DB, migrationsDir string) error {
	// Ensure goose knows the SQL dialect.
	goose.SetDialect("postgres")
	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}
	return nil
}

// OpenStdlib opens a database/sql handle for goose/other tooling.
func OpenStdlib(dsn string) (*sql.DB, error) {
	return sql.Open("pgx", dsn) // driver is registered by stdlib import
}

// Make sure stdlib driver is linked.
var _ = stdlib.RegisterConnConfig

