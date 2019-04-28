package event

import (
	"context"
	"encoding/json"

	"github.com/51st-state/api/pkg/pubsub"
)

// Producer of events
type Producer struct {
	producer pubsub.Producer
}

// Produce a event for API services
func (p *Producer) Produce(ctx context.Context, id ID, payload interface{}) error {
	e, err := new(id, payload)
	if err != nil {
		return err
	}

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return p.producer.Produce(ctx, b)
}

// NewProducer of API events
func NewProducer(p pubsub.Producer) *Producer {
	return &Producer{
		producer: p,
	}
}
