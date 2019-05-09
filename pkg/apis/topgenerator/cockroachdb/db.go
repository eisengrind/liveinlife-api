package cockroachdb

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/51st-state/api/pkg/apis/topgenerator"
)

type repository struct {
	db *sql.DB
}

// CreateSchema for the cockroachdb repository
func CreateSchema(ctx context.Context, db *sql.DB) (err error) {
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS tops (
            sex bool NOT NULL DEFAULT 0,
            undershirtId integer NOT NULL DEFAULT 0,
            topId integer NOT NULL DEFAULT 0,
            torsoId integer NOT NULL DEFAULT 0,
            clothingType integer NOT NULL DEFAULT 0,
            valencyType integer NOT NULL DEFAULT 0,
            status integer NOT NULL DEFAULT 0,
            polyesterPercentage integer NOT NULL DEFAULT 25,
            cottonPercentage integer NOT NULL DEFAULT 25,
            leatherPercentage integer NOT NULL DEFAULT 25,
            silkPercentage integer NOT NULL DEFAULT 25,
            relativeAmount integer NOT NULL DEFAULT 1,
            UNIQUE(sex, undershirtId, topId)
        );
        CREATE INDEX IF NOT EXISTS tops_idx_sex ON tops (sex);
        CREATE INDEX IF NOT EXISTS tops_idx_undershirtId ON tops (undershirtId);
        CREATE INDEX IF NOT EXISTS tops_idx_topId ON tops (topId);`,
	)
	return
}

// New cockroachdb repository
func New(db *sql.DB) topgenerator.Repository {
	return &repository{db}
}

func (r *repository) Get(ctx context.Context, id topgenerator.Identifier) (topgenerator.Complete, error) {
	inc := topgenerator.NewIncomplete(0, 0, 0, 0, 25, 25, 25, 25, 1)

	if err := r.db.QueryRowContext(
		ctx,
		`SELECT clothingType,
        valencyType,
        status,
        torsoId,
        polyesterPercentage,
        cottonPercentage,
        leatherPercentage,
        silkPercentage,
        relativeAmount
        FROM tops
        WHERE sex = $1
        AND undershirtId = $2
        AND topId = $3`,
		id.Sex(),
		id.UndershirtID(),
		id.TopID(),
	).Scan(
		&inc.Data().ClothingType,
		&inc.Data().ValencyType,
		&inc.Data().Status,
		&inc.Data().TorsoID,
		&inc.Data().PolyesterPercentage,
		&inc.Data().CottonPercentage,
		&inc.Data().LeatherPercentage,
		&inc.Data().SilkPercentage,
		&inc.Data().RelativeAmount,
	); err != nil {
		return nil, err
	}

	return newComplete(id, inc), nil
}

func (r *repository) Upsert(ctx context.Context, c topgenerator.Complete) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tops (
            sex,
            undershirtId,
            topId,
            clothingType,
            valencyType,
            status,
            torsoId,
            polyesterPercentage,
            cottonPercentage,
            leatherPercentage,
            silkPercentage,
            relativeAmount
        ) VALUES (
            $1,
            $2,
            $3,
            $4,
            $5,
            $6,
            $7,
            $8,
            $9,
            $10,
            $11,
            $12
        ) ON CONFLICT (
            sex,
            undershirtId,
            topId
        ) DO UPDATE SET clothingType = $13,
        valencyType = $14,
        status = $15,
        torsoId = $16,
        polyesterPercentage = $17,
        cottonPercentage = $18,
        leatherPercentage = $19,
        silkPercentage = $20,
        relativeAmount = $21`,
		c.Sex(),
		c.UndershirtID(),
		c.TopID(),
		c.Data().ClothingType,
		c.Data().ValencyType,
		c.Data().Status,
		c.Data().TorsoID,
		c.Data().PolyesterPercentage,
		c.Data().CottonPercentage,
		c.Data().LeatherPercentage,
		c.Data().SilkPercentage,
		c.Data().RelativeAmount,
		c.Data().ClothingType,
		c.Data().ValencyType,
		c.Data().Status,
		c.Data().TorsoID,
		c.Data().PolyesterPercentage,
		c.Data().CottonPercentage,
		c.Data().LeatherPercentage,
		c.Data().SilkPercentage,
		c.Data().RelativeAmount,
	)
	return err
}

type complete struct {
	topgenerator.Identifier
	topgenerator.Incomplete
}

func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Sex                 bool   `json:"sex"`
		UndershirtID        uint64 `json:"undershirt_id"`
		TopID               uint64 `json:"top_id"`
		ClothingType        uint8  `json:"clothing_type"`
		ValencyType         uint8  `json:"valency_type"`
		Status              uint8  `json:"status"`
		TorsoID             uint   `json:"torso_id"`
		PolyesterPercentage uint   `json:"polyester_percentage"`
		CottonPercentage    uint   `json:"cotton_percentage"`
		LeatherPercentage   uint   `json:"leather_percentage"`
		SilkPercentage      uint   `json:"silk_percentage"`
		RelativeAmount      uint   `json:"relative_amount"`
	}{
		c.Sex(),
		c.UndershirtID(),
		c.TopID(),
		c.Data().ClothingType,
		c.Data().ValencyType,
		c.Data().Status,
		c.Data().TorsoID,
		c.Data().PolyesterPercentage,
		c.Data().CottonPercentage,
		c.Data().LeatherPercentage,
		c.Data().SilkPercentage,
		c.Data().RelativeAmount,
	})
}

func newComplete(id topgenerator.Identifier, inc topgenerator.Incomplete) topgenerator.Complete {
	return &complete{
		id,
		inc,
	}
}
