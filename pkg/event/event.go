package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const version = "1"

// ID type for a specific event.
// This datatype should be used for custom
// event type definitions.
type ID string

// Event object with payload, universal unique identifier, its event id
type Event struct {
	Meta    *Meta  `json:"meta"`
	Payload []byte `json:"payload"`
}

// Meta information of an event payload
type Meta struct {
	UUID      string    `json:"uuid"`
	ID        ID        `json:"id"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

// PayloadMeta information of a specific. optional to implement
type PayloadMeta struct {
	Version string `json:"version"`
}

// new event object with payload and id
func new(id ID, payload interface{}) (*Event, error) {
	rand, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	bPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		&Meta{
			UUID:      rand.String(),
			ID:        id,
			Version:   version,
			CreatedAt: time.Now(),
		},
		bPayload,
	}, nil
}

// Decode an event from a message queue for further usage
func Decode(b []byte) (*Event, error) {
	var e Event
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, err
	}

	return &e, nil
}
