package cockroachdb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/51st-state/api/pkg/apis/inventory"
)

type db struct {
	database *sql.DB
}

// CreateSchema for the cockroachdb repository
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS inventories (
            id UUID PRIMARY KEY
        );
        CREATE UNIQUE INDEX IF NOT EXISTS inventories_idx_id ON inventories (id);
        CREATE TABLE IF NOT EXISTS inventory_items (
            inventoryId UUID references inventories (id),
            itemId TEXT NOT NULL,
            amount INTEGER NOT NULL,
            subset REAL NOT NULL,
            UNIQUE(inventoryId, itemId, subset)
        );
        CREATE UNIQUE INDEX IF NOT EXISTS inventory_items_idx_inventoryId_itemId_subset ON inventory_items (inventoryId, itemId, subset);
        CREATE INDEX IF NOT EXISTS inventory_items_idx_amount ON inventory_items (amount);`,
	)
	return
}

// NewRepository creates a new storage layer for the cockroachdb database
func NewRepository(d *sql.DB) inventory.Repository {
	return &db{d}
}

func (d *db) inventoryExists(ctx context.Context, id inventory.Identifier) error {
	var count int
	if err := d.database.QueryRowContext(
		ctx,
		`SELECT COUNT(*)
        FROM inventories
        WHERE id = $1`,
		id.GUID(),
	).Scan(
		&count,
	); err != nil {
		return err
	}

	if count != 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (d *db) Get(ctx context.Context, id inventory.Identifier) (inventory.Complete, error) {
	if err := d.inventoryExists(ctx, id); err != nil {
		return nil, err
	}

	inc := inventory.NewIncomplete(make([]*inventory.Item, 0))

	rows, err := d.database.QueryContext(
		ctx,
		`SELECT itemId,
        amount,
        subset
        FROM inventory_items
        WHERE inventoryId = $1
        AND amount > 0`,
		id.GUID(),
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item inventory.Item
		if err := rows.Scan(
			&item.ID,
			&item.Amount,
			&item.Subset,
		); err != nil {
			return nil, err
		}

		inc.Data().Items = append(inc.Data().Items, &item)
	}

	return &complete{
		id,
		inc,
	}, nil
}

func (d *db) Create(ctx context.Context) (inventory.Complete, error) {
	rand, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	if _, err := d.database.ExecContext(
		ctx,
		`INSERT INTO inventories (
            id
        ) VALUES (
            $1
        )`,
		rand.String(),
	); err != nil {
		return nil, err
	}

	return &complete{
		&identifier{rand.String()},
		inventory.NewIncomplete(make([]*inventory.Item, 0)),
	}, nil
}

func (d *db) AddItem(ctx context.Context, id inventory.Identifier, item *inventory.Item) error {
	_, err := d.database.ExecContext(
		ctx,
		`INSERT INTO inventory_items (
            inventoryId,
            itemId,
            amount,
            subset
        ) VALUES (
            $1,
            $2,
            $3,
            $4
        ) ON CONFLICT (
            inventoryId,
            itemId,
            subset
        ) DO UPDATE SET amount = inventory_items.amount + $5`,
		id.GUID(),
		item.ID,
		item.Amount,
		item.Subset,
		item.Amount,
	)
	return err
}

func (d *db) RemoveItem(ctx context.Context, id inventory.Identifier, item *inventory.Item) error {
	var guid string
	return d.database.QueryRowContext(
		ctx,
		`UPDATE inventory_items
        SET amount = amount - $1
        WHERE amount - $2 >= 0
        AND inventoryId = $3
        AND itemId = $4
        AND subset = $5
        RETURNING inventoryId`,
		item.Amount,
		item.Amount,
		id.GUID(),
		item.ID,
		item.Subset,
	).Scan(
		&guid,
	)
}

func (d *db) Delete(ctx context.Context, id inventory.Identifier) error {
	_, err := d.database.ExecContext(
		ctx,
		`DELETE FROM inventories
        WHERE id = $1`,
		id.GUID(),
	)
	return err
}
