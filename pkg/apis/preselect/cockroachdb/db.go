package cockroachdb

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/51st-state/api/pkg/apis/preselect"
)

// CreateSchema for the cockroachdb repository
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS preselections (
            sex integer NOT NULL DEFAULT 0,
            componentId integer NOT NULL DEFAULT 0,
            drawableId integer NOT NULL DEFAULT 0,
            textureId integer NOT NULL DEFAULT 0,
            accepted integer NOT NULL DEFAULT 0,
            UNIQUE(sex, componentId, drawableId, textureId)
        );
        CREATE INDEX IF NOT EXISTS preselections_idx_sex ON preselections (sex);
        CREATE INDEX IF NOT EXISTS preselections_idx_componentId ON preselections (componentId);
        CREATE INDEX IF NOT EXISTS preselections_idx_drawableId ON preselections (drawableId);
        CREATE INDEX IF NOT EXISTS preselections_idx_textureId ON preselections (textureId);`,
	)
	return
}

type complete struct {
	preselect.Identifier
	preselect.Incomplete
}

func newComplete(id preselect.Identifier, inc preselect.Incomplete) preselect.Complete {
	return &complete{id, inc}
}

func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sex         uint64 `json:"sex"`
		ComponentID uint64 `json:"component_id"`
		DrawableID  uint64 `json:"drawable_id"`
		TextureID   uint64 `json:"texture_id"`
		Accepted    uint8  `json:"accepted"`
	}{
		c.Sex(),
		c.ComponentID(),
		c.DrawableID(),
		c.TextureID(),
		c.Data().Accepted,
	})
}

type identifier struct {
	sex, componentID, drawableID, textureID uint64
}

func newIdentifier(s, c, d, t uint64) preselect.Identifier {
	return &identifier{s, c, d, t}
}

func (i *identifier) Sex() uint64 {
	return i.sex
}

func (i *identifier) ComponentID() uint64 {
	return i.componentID
}

func (i *identifier) DrawableID() uint64 {
	return i.drawableID
}

func (i *identifier) TextureID() uint64 {
	return i.textureID
}

type repository struct {
	database *sql.DB
}

// NewRepository for a preselect service
func NewRepository(db *sql.DB) preselect.Repository {
	return &repository{db}
}

func (r *repository) GetLeft(ctx context.Context) (uint64, error) {
	var count uint64
	if err := r.database.QueryRowContext(
		ctx,
		`SELECT COUNT(*)
        FROM preselections
        WHERE accepted = 0`,
	).Scan(
		&count,
	); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) GetNext(ctx context.Context) (preselect.Complete, error) {
	id := &identifier{}

	if err := r.database.QueryRowContext(
		ctx,
		`SELECT sex,
        componentId,
        drawableId,
        textureId
        FROM preselections
        WHERE accepted = 0`,
	).Scan(
		&id.sex,
		&id.componentID,
		&id.drawableID,
		&id.textureID,
	); err != nil {
		return nil, err
	}

	return newComplete(
		id,
		preselect.NewIncomplete(0),
	), nil
}

func (r *repository) Create(ctx context.Context, c ...preselect.Complete) error {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, v := range c {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO preselections (
                sex,
                componentId,
                drawableId,
                textureId,
                accepted
            ) SELECT $1,
            $2,
            $3,
            $4,
            $5
            WHERE NOT EXISTS (
                SELECT NULL
                FROM preselections
                WHERE sex = $6
                AND componentId = $7
                AND drawableId = $8
                AND textureId = $9
            )`,
			v.Sex(),
			v.ComponentID(),
			v.DrawableID(),
			v.TextureID(),
			v.Data().Accepted,
			v.Sex(),
			v.ComponentID(),
			v.DrawableID(),
			v.TextureID(),
		); err != nil {
			return func() error {
				if err := tx.Rollback(); err != nil {
					return err
				}

				return err
			}()
		}
	}

	return tx.Commit()
}

func (r *repository) Update(ctx context.Context, c ...preselect.Complete) error {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, v := range c {
		if _, err := tx.ExecContext(
			ctx,
			`UPDATE preselections
            SET accepted = $1
            WHERE sex = $2
            AND componentId = $3
            AND drawableId = $4
            AND textureId = $5`,
			v.Data().Accepted,
			v.Sex(),
			v.ComponentID(),
			v.DrawableID(),
			v.TextureID(),
		); err != nil {
			return func() error {
				if err := tx.Rollback(); err != nil {
					return err
				}

				return err
			}()
		}
	}

	return tx.Commit()
}
