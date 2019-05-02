package cockroachdb

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/51st-state/api/pkg/apis/role"
	"github.com/51st-state/api/pkg/rbac"
)

// CreateSchema for the cockroachdb repository
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS role_info (
            id TEXT PRIMARY KEY,
            title TEXT NOT NULL DEFAULT '',
            description TEXT NOT NULL DEFAULT ''
        );
        CREATE UNIQUE INDEX IF NOT EXISTS role_info_idx_id ON role_info (id);`,
	)
	return
}

type db struct {
	database *sql.DB
}

// NewRepository for storage of role information in a cockroachdb database
func NewRepository(d *sql.DB) role.Repository {
	return &db{
		d,
	}
}

type complete struct {
	role.Identifier
	role.Incomplete
}

func newComplete(id role.Identifier, inc role.Incomplete) role.Complete {
	return &complete{id, inc}
}

func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          rbac.RoleID    `json:"id"`
		Title       string         `json:"title"`
		Description string         `json:"description"`
		Rules       rbac.RoleRules `json:"rules"`
	}{
		c.ID(),
		c.Data().Title,
		c.Data().Description,
		c.Data().Rules,
	})
}

func (d *db) Get(ctx context.Context, id role.Identifier) (role.Complete, error) {
	inc := role.NewIncomplete("", "", make(rbac.RoleRules, 0))

	if err := d.database.QueryRowContext(
		ctx,
		`SELECT title,
        description
        FROM role_info
        WHERE id = $1`,
		id.ID(),
	).Scan(
		&inc.Data().Title,
		&inc.Data().Description,
	); err != nil {
		return nil, err
	}

	return newComplete(id, inc), nil
}

func (d *db) Update(ctx context.Context, c role.Complete) error {
	_, err := d.database.ExecContext(
		ctx,
		`UPDATE role_info
        SET title = $1,
        description = $2
        WHERE id = $3`,
		c.Data().Title,
		c.Data().Description,
		c.ID(),
	)
	return err
}

func (d *db) Create(ctx context.Context, c role.Complete) error {
	_, err := d.database.ExecContext(
		ctx,
		`INSERT INTO role_info (
            id,
            title,
            description
        ) SELECT $1,
        $2,
        $3`,
		c.ID(),
		c.Data().Title,
		c.Data().Description,
	)
	return err
}

func (d *db) Delete(ctx context.Context, id role.Identifier) error {
	_, err := d.database.ExecContext(
		ctx,
		`DELETE FROM role_info
        WHERE id = $1`,
		id.ID(),
	)
	return err
}
