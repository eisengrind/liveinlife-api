package event

import (
	"context"
	"encoding/json"

	"github.com/51st-state/api/pkg/pubsub"
)

// HandlerFunc for API events
type HandlerFunc func(ctx context.Context, e *Event) error

// Consumer of API events
type Consumer struct {
	consumer pubsub.Consumer
}

// NewConsumer of API events
func NewConsumer(c pubsub.Consumer) *Consumer {
	return &Consumer{
		consumer: c,
	}
}

// Consume events from the API
func (c *Consumer) Consume(ctx context.Context, h HandlerFunc) error {
	return c.consumer.Consume(ctx, func(ctx context.Context, msg pubsub.Message) error {
		var event Event
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			return err
		}

		return h(ctx, &event)
	})
}
