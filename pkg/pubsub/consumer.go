package pubsub

import "context"

// Message of an event.
// The message has to be abstracted to be used in an business
// environment, so that a possible switch to a different
// message queue is possible in the future.
//go:generate counterfeiter -o ./mocks/message.go . Message
type Message interface {
	// Ack nowledges and accepts a message
	Ack()
	// Nack (not acklowdges) a message and requeue it
	Nack()
	// Data returns the byte data of a message
	Data() []byte
}

// HandlerFunc for the consumer
type HandlerFunc func(ctx context.Context, msg Message) error

// Consumer of messages from a message queue
//go:generate counterfeiter -o ./mocks/consumer.go . Consumer
type Consumer interface {
	Consume(context.Context, HandlerFunc) error
}
