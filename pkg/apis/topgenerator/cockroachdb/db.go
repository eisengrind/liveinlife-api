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
            sex integer NOT NULL DEFAULT 0,
            undershirtId integer NOT NULL DEFAULT 0,
            undershirtTextureId integer NOT NULL DEFAULT 0,
            overshirtId integer NOT NULL DEFAULT 0,
            overshirtTextureId integer NOT NULL DEFAULT 0,
            torsoId integer NOT NULL DEFAULT 0,
            torsoTextureId integer NOT NULL DEFAULT 0,
            clothingType integer NOT NULL DEFAULT 0,
            valencyType integer NOT NULL DEFAULT 0,
            name text NOT NULL DEFAULT '',
            notice text NOT NULL DEFAULT '',
            status integer NOT NULL DEFAULT 0,
            polyesterPercentage integer NOT NULL DEFAULT 25,
            cottonPercentage integer NOT NULL DEFAULT 25,
            leatherPercentage integer NOT NULL DEFAULT 25,
            silkPercentage integer NOT NULL DEFAULT 25,
            relativeAmount integer NOT NULL DEFAULT 1,
            UNIQUE(sex, undershirtId, undershirtTextureId, overshirtId, overshirtTextureId)
        );
        CREATE INDEX IF NOT EXISTS tops_idx_sex ON tops (sex);
        CREATE INDEX IF NOT EXISTS tops_idx_undershirtId ON tops (undershirtId);
        CREATE INDEX IF NOT EXISTS tops_idx_undershirtTextureId ON tops (undershirtTextureId);
        CREATE INDEX IF NOT EXISTS tops_idx_overshirtId ON tops (overshirtId);
        CREATE INDEX IF NOT EXISTS tops_idx_overshirtTextureId ON tops (overshirtTextureId)`,
	)
	return
}

// New cockroachdb repository
func New(db *sql.DB) top.Repository {
	return &repository{db}
}

func (r *repository) Get(ctx context.Context, id top.Identifier) (top.Complete, error) {
	inc := top.NewIncomplete("", "", 0, 0, 0, 0, 0, 25, 25, 25, 25, 1)

	if err := r.db.QueryRowContext(
		ctx,
		`SELECT clothingType,
        valencyType,
        name,
        notice,
        status,
        torsoId,
        torsoTextureId,
        polyesterPercentage,
        cottonPercentage,
        leatherPercentage,
        silkPercentage,
        relativeAmount
        FROM tops
        WHERE sex = $1
        AND undershirtId = $2
        AND undershirtTextureId = $3
        AND overshirtId = $4
        AND overshirtTextureId = $5`,
		id.Sex(),
		id.UndershirtID(),
		id.UndershirtTextureID(),
		id.OvershirtID(),
		id.OvershirtTextureID(),
	).Scan(
		&inc.Data().ClothingType,
		&inc.Data().ValencyType,
		&inc.Data().Name,
		&inc.Data().Notice,
		&inc.Data().Status,
		&inc.Data().TorsoID,
		&inc.Data().TorsoTextureID,
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

func (r *repository) Upsert(ctx context.Context, c top.Complete) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tops (
            sex,
            undershirtId,
            undershirtTextureId,
            overshirtId,
            overshirtTextureId,
            clothingType,
            valencyType,
            name,
            notice,
            status,
            torsoId,
            torsoTextureId,
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
            $12,
            $13,
            $14,
            $15,
            $16,
            $17
        ) ON CONFLICT (
            sex,
            undershirtId,
            undershirtTextureId,
            overshirtId,
            overshirtTextureId
        ) DO UPDATE SET clothingType = $18,
        valencyType = $19,
        name = $20,
        notice = $21,
        status = $22,
        torsoId = $23,
        torsoTextureId = $24,
        polyesterPercentage = $25,
        cottonPercentage = $26,
        leatherPercentage = $27,
        silkPercentage = $28,
        relativeAmount = $29`,
		c.Sex(),
		c.UndershirtID(),
		c.UndershirtTextureID(),
		c.OvershirtID(),
		c.OvershirtTextureID(),
		c.Data().ClothingType,
		c.Data().ValencyType,
		c.Data().Name,
		c.Data().Notice,
		c.Data().Status,
		c.Data().TorsoID,
		c.Data().TorsoTextureID,
		c.Data().PolyesterPercentage,
		c.Data().CottonPercentage,
		c.Data().LeatherPercentage,
		c.Data().SilkPercentage,
		c.Data().RelativeAmount,
		c.Data().ClothingType,
		c.Data().ValencyType,
		c.Data().Name,
		c.Data().Notice,
		c.Data().Status,
		c.Data().TorsoID,
		c.Data().TorsoTextureID,
		c.Data().PolyesterPercentage,
		c.Data().CottonPercentage,
		c.Data().LeatherPercentage,
		c.Data().SilkPercentage,
		c.Data().RelativeAmount,
	)
	return err
}

type complete struct {
	top.Identifier
	top.Incomplete
}

func (c *complete) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Sex                 uint8  `json:"sex"`
		UndershirtID        uint   `json:"undershirt_id"`
		UndershirtTextureID uint   `json:"undershirt_texture_id"`
		OvershirtID         uint   `json:"overshirt_id"`
		OvershirtTextureID  uint   `json:"overshirt_texture_id"`
		Name                string `json:"name"`
		Notice              string `json:"notice"`
		ClothingType        uint8  `json:"clothing_type"`
		ValencyType         uint8  `json:"valency_type"`
		Status              uint8  `json:"status"`
		TorsoID             uint   `json:"torso_id"`
		TorsoTextureID      uint   `json:"torso_texture_id"`
		PolyesterPercentage uint   `json:"polyester_percentage"`
		CottonPercentage    uint   `json:"cotton_percentage"`
		LeatherPercentage   uint   `json:"leather_percentage"`
		SilkPercentage      uint   `json:"silk_percentage"`
		RelativeAmount      uint   `json:"relative_amount"`
	}{
		c.Sex(),
		c.UndershirtID(),
		c.UndershirtTextureID(),
		c.OvershirtID(),
		c.OvershirtTextureID(),
		c.Data().Name,
		c.Data().Notice,
		c.Data().ClothingType,
		c.Data().ValencyType,
		c.Data().Status,
		c.Data().TorsoID,
		c.Data().TorsoTextureID,
		c.Data().PolyesterPercentage,
		c.Data().CottonPercentage,
		c.Data().LeatherPercentage,
		c.Data().SilkPercentage,
		c.Data().RelativeAmount,
	})
}

func newComplete(id top.Identifier, inc top.Incomplete) top.Complete {
	return &complete{
		id,
		inc,
	}
}
