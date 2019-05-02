package cockroachdb

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/51st-state/api/pkg/apis/user"
)

// CreateSchema for the cockroachdb repository
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            wcfUserId integer NOT NULL DEFAULT 0,
            gameSerialHash text NOT NULL DEFAULT '',
            banned boolean NOT NULL DEFAULT false,
            UNIQUE(id),
            UNIQUE(wcfUserId),
            UNIQUE(gameSerialHash)
        );
        CREATE UNIQUE INDEX IF NOT EXISTS users_idx_id ON users (id);
        CREATE UNIQUE INDEX IF NOT EXISTS users_idx_wcfUserId ON users (wcfUserId);
        CREATE UNIQUE INDEX IF NOT EXISTS users_idx_gameSerialHash ON users (gameSerialHash);`,
	)
	return
}

type identifier struct {
	uuid string
}

func newIdentifier(uuid string) user.Identifier {
	return &identifier{uuid}
}

func (i *identifier) UUID() string {
	return i.uuid
}

type complete struct {
	user.Identifier
	user.Incomplete
}

func newComplete(id user.Identifier, inc user.Incomplete) user.Complete {
	return &complete{
		id,
		inc,
	}
}

type completeUser struct {
	UUID           string         `json:"uuid"`
	WCFUserID      user.WCFUserID `json:"wcf_user_id"`
	GameSerialHash string         `json:"game_serial_hash"`
	Banned         bool           `json:"banned"`
}

// MarshalJSON information of a user
func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(&completeUser{
		c.UUID(),
		c.Data().WCFUserID,
		c.Data().GameSerialHash,
		c.Data().Banned,
	})
}

type repository struct {
	database *sql.DB
}

// NewRepository for the user service for cockroachdb
func NewRepository(db *sql.DB) user.Repository {
	return &repository{db}
}

func (r *repository) Get(ctx context.Context, id user.Identifier) (user.Complete, error) {
	inc := user.NewIncomplete(0, "", false)

	if err := r.database.QueryRowContext(
		ctx,
		`SELECT wcfUserId,
        gameSerialHash,
        banned
        FROM users
        WHERE id = $1`,
		id.UUID(),
	).Scan(
		&inc.Data().WCFUserID,
		&inc.Data().GameSerialHash,
		&inc.Data().Banned,
	); err != nil {
		return nil, err
	}

	return newComplete(id, inc), nil
}

func (r *repository) GetByWCFUserID(ctx context.Context, wcfUserID user.WCFUserID) (user.Complete, error) {
	inc := user.NewIncomplete(wcfUserID, "", false)
	var id string

	if err := r.database.QueryRowContext(
		ctx,
		`SELECT id,
        banned,
        gameSerialHash
        FROM users
        WHERE wcfUserId = $1`,
		wcfUserID,
	).Scan(
		&id,
		&inc.Data().Banned,
		&inc.Data().GameSerialHash,
	); err != nil {
		return nil, err
	}

	return newComplete(
		newIdentifier(id),
		inc,
	), nil
}

func (r *repository) GetByGameSerialHash(ctx context.Context, hash string) (user.Complete, error) {
	inc := user.NewIncomplete(0, hash, false)
	var id string

	if err := r.database.QueryRowContext(
		ctx,
		`SELECT id,
        wcfUserId,
        banned
        FROM users
        WHERE gameSerialHash = $1`,
		hash,
	).Scan(
		&id,
		&inc.Data().WCFUserID,
		&inc.Data().Banned,
	); err != nil {
		return nil, err
	}

	return newComplete(
		newIdentifier(id),
		inc,
	), nil
}

func (r *repository) Create(ctx context.Context, inc user.Incomplete) (user.Complete, error) {
	rand, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	if _, err := r.database.ExecContext(
		ctx,
		`INSERT INTO users (
            id,
            wcfUserId,
            gameSerialHash,
            banned
        ) VALUES (
            $1,
            $2,
            $3,
            $4
        )`,
		rand.String(),
		inc.Data().WCFUserID,
		inc.Data().GameSerialHash,
		inc.Data().Banned,
	); err != nil {
		return nil, err
	}

	return newComplete(
		newIdentifier(rand.String()),
		inc,
	), nil
}

func (r *repository) Update(ctx context.Context, c user.Complete) error {
	_, err := r.database.ExecContext(
		ctx,
		`UPDATE users
        SET wcfUserId = $1,
        gameSerialHash = $2,
        banned = $3
        WHERE id = $4`,
		c.Data().WCFUserID,
		c.Data().GameSerialHash,
		c.Data().Banned,
		c.UUID(),
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id user.Identifier) error {
	_, err := r.database.ExecContext(
		ctx,
		`DELETE FROM users
        WHERE id = $1`,
		id.UUID(),
	)
	return err
}
