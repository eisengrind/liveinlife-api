package preselect

import "encoding/json"

// Identifier of a pre-selection
type Identifier interface {
	Sex() uint64
	ComponentID() uint64
	DrawableID() uint64
	TextureID() uint64
}

type identifier struct {
	sex, componentID, drawableID, textureID uint64
}

func newIdentifier(s, c, d, t uint64) Identifier {
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

// Provider of data for an incomplete pre-selection object
type Provider interface {
	Data() *data
}

// Incomplete pre-selection object
type Incomplete interface {
	Provider
}

// Complete pre-selection object
type Complete interface {
	Identifier
	Incomplete
}

type complete struct {
	Identifier
	Incomplete
}

func newComplete(id Identifier, inc Incomplete) Complete {
	return &complete{id, inc}
}

func (c *complete) UnmarshalJSON(b []byte) error {
	var compl struct {
		Sex         uint64 `json:"sex"`
		ComponentID uint64 `json:"component_id"`
		DrawableID  uint64 `json:"drawable_id"`
		TextureID   uint64 `json:"texture_id"`
		Accepted    uint8  `json:"accepted"`
	}

	if err := json.Unmarshal(b, &compl); err != nil {
		return err
	}

	c.Identifier = newIdentifier(compl.Sex, compl.ComponentID, compl.DrawableID, compl.TextureID)
	c.Incomplete = NewIncomplete(compl.Accepted)

	return nil
}

type data struct {
	Accepted uint8 `json:"accepted"`
}

// NewIncomplete pre-selection object
func NewIncomplete(a uint8) Incomplete {
	return &data{
		a,
	}
}

func (d *data) Data() *data {
	return d
}

func (d *data) SetAccepted(to uint8) *data {
	d.Accepted = to
	return d
}
