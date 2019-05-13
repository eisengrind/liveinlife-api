package cockroachdb

import (
	"context"
	"database/sql"
	"time"

	"github.com/51st-state/api/pkg/apis/auth"
)

type db struct {
	database *sql.DB
}

// NewRepository creates a new instance of a repository provider for a cockroachdb database
func NewRepository(d *sql.DB) auth.Repository {
	return &db{d}
}

// CreateSchema creates a new cockroachdb database schema
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS login_attempts (
            id TEXT,
            attemptedAt TIMESTAMPTZ NOT NULL DEFAULT NOW()
        );
        CREATE UNIQUE INDEX IF NOT EXISTS login_attempts_idx_id_attemptedAt ON login_attempts(id, attemptedAt);`,
	)
	return
}

func (d *db) LoginAttemptsCountSince(ctx context.Context, id string, t time.Time) (uint64, error) {
	var count uint64
	if err := d.database.QueryRowContext(
		ctx,
		`SELECT COUNT(*)
        FROM login_attempts
        WHERE id = $1
        AND attemptedAt >= $2`,
		id,
		t,
	).Scan(
		&count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

func (d *db) AddLoginAttempt(ctx context.Context, id string, t time.Time) error {
	_, err := d.database.ExecContext(
		ctx,
		`INSERT INTO login_attempts (
            id,
            attemptedAt
        ) SELECT $1,
        $2`,
		id,
		t,
	)
	return err
}
