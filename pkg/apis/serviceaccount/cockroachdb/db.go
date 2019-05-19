package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/51st-state/api/pkg/apis/serviceaccount/key"

	"github.com/google/uuid"

	"github.com/51st-state/api/pkg/apis/serviceaccount"
)

// CreateSchema creates a new cockroachdb schema in a cockroachdb database for the serviceaccount service
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS service_accounts (
            guid UUID PRIMARY KEY,
            name TEXT NOT NULL DEFAULT '',
            description TEXT NOT NULL DEFAULT ''
        );
        CREATE UNIQUE INDEX IF NOT EXISTS service_accounts_idx_guid ON service_accounts (guid);`,
	)
	return
}

type db struct {
	db *sql.DB
}

// NewRepository creates a new cockroachdb db storage repository
func NewRepository(d *sql.DB) serviceaccount.Repository {
	return &db{d}
}

func (d *db) Get(ctx context.Context, id serviceaccount.Identifier) (serviceaccount.Complete, error) {
	c := &complete{
		id,
		serviceaccount.NewIncomplete("", ""),
	}
	if err := d.db.QueryRowContext(
		ctx,
		`SELECT name,
        description
        FROM service_accounts
        WHERE guid = $1`,
		id.GUID(),
	).Scan(
		&c.Data().Name,
		&c.Data().Description,
	); err != nil {
		return nil, err
	}

	return c, nil
}

func (d *db) Create(ctx context.Context, inc serviceaccount.Incomplete) (serviceaccount.Complete, error) {
	rand, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	if _, err := d.db.ExecContext(
		ctx,
		`INSERT INTO service_accounts (
            guid,
            name,
            description
        ) SELECT $1,
        $2,
        $3`,
		rand.String(),
		inc.Data().Name,
		inc.Data().Description,
	); err != nil {
		return nil, err
	}

	return &complete{
		key.NewIdentifier(rand.String()),
		inc,
	}, nil
}

func (d *db) Update(ctx context.Context, c serviceaccount.Complete) error {
	_, err := d.db.ExecContext(
		ctx,
		`UPDATE service_accounts
        SET name = $1,
        description = $2
        WHERE guid = $3`,
		c.Data().Name,
		c.Data().Description,
		c.GUID(),
	)
	return err
}

func (d *db) Delete(ctx context.Context, id serviceaccount.Identifier) error {
	_, err := d.db.ExecContext(
		ctx,
		`DELETE FROM service_accounts
        WHERE guid = $1`,
		id.GUID(),
	)
	return err
}
