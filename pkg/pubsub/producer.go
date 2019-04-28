package pubsub

import "context"

// Producer interface for message queue implementation
//go:generate counterfeiter -o ./mocks/producer.go . Producer
type Producer interface {
	Produce(ctx context.Context, data []byte) error
}
