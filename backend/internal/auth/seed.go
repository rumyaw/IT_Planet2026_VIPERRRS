package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureAdmin creates an initial ADMIN user for curator access if env vars are provided.
func EnsureAdmin(ctx context.Context, pool *pgxpool.Pool, adminEmail string, adminPassword string) error {
	if adminEmail == "" || adminPassword == "" {
		return nil
	}

	var exists bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, adminEmail).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}

	pwHash, err := HashPassword(adminPassword)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO users (email, password_hash, role, status, display_name)
		 VALUES ($1,$2,'ADMIN','ACTIVE','Администратор')`,
		adminEmail, pwHash,
	)
	return err
}

