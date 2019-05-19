package cockroachdb

import (
	"bytes"
	"context"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/google/uuid"

	"github.com/51st-state/api/pkg/apis/serviceaccount/key"
)

// CreateSchema creates a new cockroachdb database schema for service account keys
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS service_account_keys (
            guid UUID PRIMARY KEY,
            serviceAccountGUID UUID NOT NULL,
            name TEXT NOT NULL DEFAULT '',
            description TEXT NOT NULL DEFAULT '',
            pemPublicKey TEXT NOT NULL DEFAULT ''
        );
        CREATE UNIQUE INDEX IF NOT EXISTS service_account_keys_idx_guid ON service_account_keys (guid);`,
	)
	return
}

type db struct {
	database *sql.DB
}

// NewRepository creates a new repository using the cockroachdb database
func NewRepository(d *sql.DB) key.Repository {
	return &db{
		d,
	}
}

var (
	errInvalidPEMKey       = errors.New("invalid pem key")
	errInvalidPEMBlockType = errors.New("invalid pem block type")
)

func (d *db) Get(ctx context.Context, id key.Identifier) (key.Complete, error) {
	var pemPublicKey string
	inc := key.NewIncomplete("", "")

	if err := d.database.QueryRowContext(
		ctx,
		`SELECT serviceAccountGUID,
        name,
        description,
        pemPublicKey
        FROM service_account_keys
        WHERE guid = $1`,
		id.GUID(),
	).Scan(
		&inc.Data().ServiceAccountGUID,
		&inc.Data().Name,
		&inc.Data().Description,
		&pemPublicKey,
	); err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(pemPublicKey))
	if block == nil {
		return nil, errInvalidPEMKey
	}

	if block.Type != "PUBLIC KEY" {
		return nil, errInvalidPEMBlockType
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	inc.Data().PublicKey = publicKey

	return &complete{
		id,
		inc,
	}, nil
}

func (d *db) Update(ctx context.Context, c key.Complete) error {
	_, err := d.database.ExecContext(
		ctx,
		`UPDATE service_account_keys
        SET name = $1,
        description = $2
        WHERE guid = $3`,
		c.Data().Name,
		c.Data().Description,
		c.GUID(),
	)
	return err
}

func (d *db) Create(ctx context.Context, inc key.Incomplete) (key.Complete, error) {
	var buf = bytes.NewBuffer([]byte{})
	if err := pem.Encode(buf, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(inc.Data().PublicKey),
	}); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	rand, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	if _, err := d.database.ExecContext(
		ctx,
		`INSERT INTO service_account_keys (
            guid UUID PRIMARY KEY,
            serviceAccountGUID UUID NOT NULL,
            name TEXT NOT NULL DEFAULT '',
            description TEXT NOT NULL DEFAULT '',
            pemPublicKey TEXT NOT NULL DEFAULT ''
        ) SELECT $1,
        $2,
        $3,
        $4,
        $5`,
		rand.String(),
		inc.Data().ServiceAccountGUID,
		inc.Data().Name,
		inc.Data().Description,
		b,
	); err != nil {
		return nil, err
	}

	return &complete{
		key.NewIdentifier(rand.String()),
		inc,
	}, nil
}

func (d *db) Delete(ctx context.Context, id key.Identifier) error {
	_, err := d.database.ExecContext(
		ctx,
		`DELETE FROM service_account_keys
        WHERE guid = $1`,
		id.GUID(),
	)
	return err
}
